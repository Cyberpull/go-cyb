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
)

type ClientAuthSubscriber func(ref *ClientRef) (err error)

type Client struct {
	uuid                    string
	err                     error
	ref                     *ClientRef
	conn                    *tls.Conn
	mutex                   sync.Mutex
	timeout                 time.Duration
	opts                    ClientOptions
	authSubscribers         []ClientAuthSubscriber
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

func (c *Client) connect() (err error) {
	address := address(&c.opts)

	log.Magentafln("Connecting to %s...", address)

	c.isRunning = true

	defer func() {
		c.ref = nil

		c.isRunning = false

		c.uuid = ""

		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	var conn *tls.Conn

	dialer := &net.Dialer{Timeout: c.timeout}
	if conn, err = tls.DialWithDialer(dialer, "tcp", address, c.opts.TlsConfig); err != nil {
		return
	}

	c.ref = newClientRef(conn)

	defer c.ref.close()

	log.Successfln("Connected to %s", address)

	if err = c.sendClientInformation(); err != nil {
		return
	}

	if err = c.execAuth(); err != nil {
		return
	}

	if err = c.receiveIdentifier(); err != nil {
		return
	}

	err = c.runSession()

	return
}

func (c *Client) Start() {
	var err error

	defer func() {
		c.isStopped = false
	}()

	if err = sanitizeNameAndAlias(&c.opts); err != nil {
		log.Errorln(err)
		return
	}

	if err = sanitizeTlsConfig(&c.opts); err != nil {
		log.Errorln(err)
		return
	}

	for {
		if err = c.connect(); err != nil {
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
	// TODO: Process Data
}

func (c *Client) EnsureStarted() (err error) {
	for {
		if c.isRunning && !c.isRunningSession {
			continue
		}

		err = c.err

		break
	}

	return
}

/******************************************/

func NewClient(opts ClientOptions) *Client {
	return &Client{
		opts:                    opts,
		authSubscribers:         make([]ClientAuthSubscriber, 0),
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
