package socket

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
)

type ClientAuthSubscriber func(conn *tls.Conn) (err error)

type Client struct {
	uuid                    string
	err                     error
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

func (c *Client) execAuth(conn *tls.Conn) (err error) {
	for _, subscriber := range c.authSubscribers {
		err = subscriber(conn)

		if err != nil {
			break
		}
	}

	return
}

func (c *Client) receiveIdentifier(reader *bufio.Reader) (err error) {
	// TODO: Receive Identifier
	return
}

func (c *Client) sendClientInformation() (err error) {
	if err = sendClientName(c); err != nil {
		return
	}

	err = sendClientAlias(c)

	return
}

func (c *Client) connect() (err error) {
	address := address(&c.opts)

	log.Magentafln("Connecting to %s...", address)

	c.isRunning = true

	defer func() {
		c.conn = nil

		c.isRunning = false

		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	dialer := &net.Dialer{Timeout: c.timeout}
	if c.conn, err = tls.DialWithDialer(dialer, "tcp", address, c.opts.TlsConfig); err != nil {
		return
	}

	defer c.conn.Close()

	log.Successfln("Connected to %s", address)

	log.Magentaln("Registering client information...")

	if err = c.sendClientInformation(); err != nil {
		return
	}

	log.Successln("Registered client information")

	if err = c.execAuth(c.conn); err != nil {
		return
	}

	log.Successln("Authorized")

	reader := bufio.NewReader(c.conn)

	if err = c.receiveIdentifier(reader); err != nil {
		return
	}

	err = c.runSession(reader)

	return
}

func (c *Client) Start() {
	var err error

	defer func() {
		c.isStopped = false

		if err != nil {
			log.Errorln(err)
		}
	}()

	if err = sanitizeNameAndAlias(&c.opts); err != nil {
		return
	}

	if err = sanitizeTlsConfig(&c.opts); err != nil {
		return
	}

	for {
		if err = c.connect(); err != nil {
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
	c.isStopped = true
	c.conn.Close()
	return
}

func (c *Client) runSession(reader *bufio.Reader) (err error) {
	c.isRunningSession = true

	defer func() {
		c.isRunningSession = false
	}()

	var data []byte

	for {
		if data, err = reader.ReadBytes('\n'); err != nil {
			break
		}

		if err2 := c.checkError(data); err2 != nil {
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

func (c *Client) checkError(data []byte) (err error) {
	c.mutex.Lock()

	defer c.mutex.Unlock()

	if !bytes.HasPrefix(data, []byte(ErrorPrefix)) {
		return
	}

	data = bytes.TrimPrefix(data, []byte(ErrorPrefix))
	err = errors.New(string(data))

	c.conn.Write(append([]byte(ErrorRcpt), '\n'))

	return
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
	return
}

func sendClientAlias(c *Client) (err error) {
	return
}
