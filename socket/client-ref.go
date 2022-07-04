package socket

import (
	"bufio"
	"bytes"
	"net"
	"sync"

	"cyberpull.com/go-cyb/errors"
)

type ClientRef struct {
	conn   net.Conn
	reader *bufio.Reader
	mutex  sync.Mutex
}

func (c *ClientRef) validate() (err error) {
	if c.conn == nil || c.reader == nil {
		err = errors.New("ClientRef not properly instantiated")
	}

	return
}

func (c *ClientRef) Write(b []byte) (i int, err error) {
	if err = c.validate(); err != nil {
		return
	}

	i, err = c.conn.Write(b)

	return
}

func (c *ClientRef) WriteString(d string) (int, error) {
	return c.Write([]byte(d))
}

func (c *ClientRef) Writeln(b []byte) (int, error) {
	return c.Write(append(b, '\n'))
}

func (c *ClientRef) WriteStringln(d string) (int, error) {
	return c.Writeln([]byte(d))
}

func (c *ClientRef) ReadBytes(delim byte) (value []byte, err error) {
	if err = c.validate(); err != nil {
		return
	}

	value, err = c.reader.ReadBytes(delim)

	return
}

func (c *ClientRef) ReadString(delim byte) (value string, err error) {
	if err = c.validate(); err != nil {
		return
	}

	value, err = c.reader.ReadString(delim)

	return
}

func (c *ClientRef) checkError(data []byte) (err error) {
	c.mutex.Lock()

	defer c.mutex.Unlock()

	if !bytes.HasPrefix(data, []byte(ErrorPrefix)) {
		return
	}

	data = bytes.TrimPrefix(data, []byte(ErrorPrefix))
	err = errors.New(string(data))

	c.WriteStringln(ErrorRcpt)

	return
}

func (c *ClientRef) close() error {
	return c.conn.Close()
}

/**********************************************/

func newClientRef(conn net.Conn) *ClientRef {
	return &ClientRef{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}
