package socket

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"sync"

	"cyberpull.com/go-cyb/errors"
)

type ClientRef struct {
	conn   *tls.Conn
	reader *bufio.Reader
	mutex  sync.Mutex
}

func (c *ClientRef) checkAndValidateInstance() {
	if c.conn == nil || c.reader == nil {
		err := errors.Newf("ClientRef not properly instanciated")
		panic(err)
	}
}

func (c *ClientRef) Write(b []byte) (int, error) {
	c.checkAndValidateInstance()
	return c.conn.Write(b)
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

func (c *ClientRef) ReadBytes(delim byte) ([]byte, error) {
	c.checkAndValidateInstance()
	return c.reader.ReadBytes(delim)
}

func (c *ClientRef) ReadString(delim byte) (string, error) {
	c.checkAndValidateInstance()
	return c.reader.ReadString(delim)
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

func newClientRef(conn *tls.Conn) *ClientRef {
	return &ClientRef{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}
