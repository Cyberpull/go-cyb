package socket

import (
	"bytes"
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"time"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/objects"
)

type ClientAuthSubscriber func(ref *ClientRef) (err error)
type ClientUpdateSubscriber func(collection *ClientUpdateHandlerCollection) (err error)

type Client struct {
	uuid                    string
	err                     error
	ref                     *ClientRef
	conn                    *tls.Conn
	mutex                   sync.Mutex
	timeout                 time.Duration
	opts                    ClientOptions
	authSubscribers         []ClientAuthSubscriber
	updateSubscribers       []ClientUpdateSubscriber
	updateHandlerCollection *ClientUpdateHandlerCollection
	responseCollection      *ClientResponseCollection
	isRunningSession        bool
	isRunning               bool
	isStopped               bool
}

func (c *Client) Auth(subs ...ClientAuthSubscriber) {
	c.authSubscribers = append(c.authSubscribers, subs...)
}

func (c *Client) execAuth() (err error) {
	if len(c.authSubscribers) > 0 {
		log.Magentafln(`Authorizing %s...`, c.opts.Name)

		for _, subscriber := range c.authSubscribers {
			err = subscriber(c.ref)

			if err != nil {
				break
			}
		}

		if err == nil {
			log.Successln("Authorized")
		}
	}

	return
}

func (c *Client) Update(subs ...ClientUpdateSubscriber) {
	c.updateSubscribers = append(c.updateSubscribers, subs...)
}

func (c *Client) execUpdate() (err error) {
	if len(c.authSubscribers) > 0 {
		log.Magentaln(`Subscribing to updates...`)

		for _, subscriber := range c.updateSubscribers {
			err = subscriber(c.updateHandlerCollection)

			if err != nil {
				break
			}
		}

		if err == nil {
			log.Successln("Subscribed")
		}
	}

	return
}

func (c *Client) On(method, channel string, handler ClientUpdateHander) {
	c.updateHandlerCollection.On(method, channel, handler)
}

func (c *Client) receiveIdentifier() (err error) {
	log.Magentafln("Receiving identifier for %s...", c.opts.Name)

	var (
		data   []byte
		prefix string = "UUID "
	)

	if data, err = c.ref.ReadBytes('\n'); err != nil {
		return
	}

	if !bytes.HasPrefix(data, []byte(prefix)) {
		err = errors.New("Invalid UUID information")
		return
	}

	data = bytes.TrimPrefix(data, []byte(prefix))
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		err = errors.New("Invalid UUID")
		return
	}

	if _, err = c.ref.WriteStringln("RECEIVED"); err != nil {
		return
	}

	c.uuid = string(data)

	log.Successln("Received identifier")

	return
}

func (c *Client) sendClientInformation() (err error) {
	log.Magentafln("Registering %s...", c.opts.Name)

	if err = sendClientName(c); err != nil {
		return
	}

	if err = sendClientAlias(c); err != nil {
		return
	}

	log.Successln("Registered")

	return
}

func (c *Client) connect(errChan ...chan error) (err error) {
	address := address(&c.opts)

	log.Magentafln("Connecting to %s...", address)

	c.isRunning = true

	defer func() {
		c.ref = nil

		c.isRunning = false

		c.uuid = ""

		if r := recover(); r != nil {
			err = errors.From(r)
			writeOne(errChan, err)
		}
	}()

	var conn *tls.Conn

	dialer := &net.Dialer{Timeout: c.timeout}
	if conn, err = tls.DialWithDialer(dialer, "tcp", address, c.opts.TlsConfig); err != nil {
		writeOne(errChan, err)
		return
	}

	c.ref = newClientRef(conn)

	defer c.ref.close()

	log.Successfln("Connected to %s", address)

	if err = c.sendClientInformation(); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = c.execAuth(); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = c.execUpdate(); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = c.receiveIdentifier(); err != nil {
		writeOne(errChan, err)
		return
	}

	writeOne(errChan, nil)

	err = c.runSession()

	return
}

func (c *Client) Start(errChan ...chan error) {
	var err error

	defer func() {
		c.isStopped = false
	}()

	if err = sanitizeNameAndAlias(&c.opts); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = sanitizeTlsConfig(&c.opts); err != nil {
		writeOne(errChan, err)
		return
	}

	for {
		if err = c.connect(errChan...); err != nil {
			if c.isStopped {
				break
			}

			duration := time.Second * ClientRetry

			log.Errorln(err)
			log.Printfln("Retrying in %d seconds...", ClientRetry)

			time.Sleep(duration)

			continue
		}

		return
	}
}

func (c *Client) Stop() (err error) {
	if c.ref != nil {
		c.isStopped = true
		c.ref.close()
	}

	return
}

func (c *Client) runSession() (err error) {
	c.isRunningSession = true

	defer func() {
		c.isRunningSession = false
	}()

	log.Successln("Session started")

	var data []byte

	for {
		if data, err = c.ref.ReadBytes('\n'); err != nil {
			break
		}

		if err2 := c.ref.checkError(data); err2 != nil {
			log.Errorln(err2)
			continue
		}

		go c.processData(data)
	}

	return
}

func (c *Client) processData(data []byte) {
	var err error

	defer func() {
		if err != nil {
			log.Errorln(err)
		}
	}()

	if len(data) == 0 {
		return
	}

	delimIndex := bytes.Index(data, []byte("::"))
	prefix := string(data[:delimIndex])
	data = data[delimIndex+2:]

	switch prefix {
	case ResponseTxt:
		delimIndex = bytes.Index(data, []byte("::"))

		if delimIndex < 0 {
			err = errors.New("Invalid response")
			return
		}

		requestUUID := string(data[:delimIndex])
		requestUUID = strings.TrimSpace(requestUUID)

		if requestUUID == "" {
			err = errors.New("Invalid response uuid")
			return
		}

		data = data[delimIndex+2:]

		out := &Output{}

		if err = objects.ParseJSON(data, out); err != nil {
			return
		}

		err = c.responseCollection.Set(requestUUID, out)

	case UpdateTxt:
		out := &Output{}

		if err = objects.ParseJSON(data, out); err != nil {
			return
		}

		go c.updateHandlerCollection.updateAll(out)

	default:
		err = errors.New("Unable to process data")
	}
}

func (c *Client) request(method, channel string, data any, timeout ...time.Duration) (out *Output, err error) {
	method = strings.ToUpper(method)

	var request *Request

	if request, err = newRequest(c); err != nil {
		return
	}

	if err = request.SetData(data); err != nil {
		return
	}

	request.Method = method
	request.Channel = channel

	var requestBytes []byte

	if requestBytes, err = objects.ToJSON(request); err != nil {
		return
	}

	if _, err = c.ref.Writeln(requestBytes); err != nil {
		return
	}

	var tmpOut *Output

	if tmpOut, err = c.responseCollection.Get(request, timeout...); err != nil {
		return
	}

	if err = tmpOut.GetError(); err != nil {
		return
	}

	out = tmpOut

	return
}

/******************************************/

func NewClient(opts ClientOptions) *Client {
	return &Client{
		opts:                    opts,
		authSubscribers:         make([]ClientAuthSubscriber, 0),
		updateSubscribers:       make([]ClientUpdateSubscriber, 0),
		updateHandlerCollection: newClientUpdateHandlerCollection(),
		responseCollection:      newClientResponseCollection(),
	}
}

func sendClientName(c *Client) (err error) {
	var data string

	if data, err = c.ref.ReadString('\n'); err != nil {
		return
	}

	data = strings.TrimSpace(data)

	if data != "CLIENT NAME:" {
		err = errors.Newf(`Expected "CLIENT NAME:", got "%s" instead.`, 500, data)
		return
	}

	_, err = c.ref.WriteStringln(c.opts.Name)

	return
}

func sendClientAlias(c *Client) (err error) {
	var data string

	if data, err = c.ref.ReadString('\n'); err != nil {
		return
	}

	data = strings.TrimSpace(data)

	if data != "CLIENT ALIAS:" {
		err = errors.Newf(`Expected "CLIENT ALIAS:", got "%s" instead.`, 500, data)
		return
	}

	_, err = c.ref.WriteStringln(c.opts.Alias)

	return
}
