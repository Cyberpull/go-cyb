package socket

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/objects"
	"cyberpull.com/go-cyb/uuid"
)

type ServerClientRef struct {
	uuid   string
	name   string
	alias  string
	opts   ServerOptions
	mutex  sync.Mutex
	conn   net.Conn
	reader *bufio.Reader
}

func (s *ServerClientRef) sendError(message string) (err error) {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	if _, err = s.WriteStringln(ErrorPrefix + message); err != nil {
		return
	}

	var rcpt string

	if rcpt, err = s.ReadString('\n'); err != nil {
		return
	}

	rcpt = strings.TrimSpace(rcpt)

	if rcpt != ErrorRcpt {
		err = errors.Newf(`Expected "%s", got "%s" instead.`, 500, ErrorRcpt, rcpt)
	}

	return
}

func (s *ServerClientRef) sendIdentifier() (err error) {
	data := fmt.Sprintf("UUID %s\n", s.uuid)

	if _, err = s.WriteString(data); err != nil {
		return
	}

	rcpt, err := s.ReadString('\n')

	if err != nil {
		return
	}

	rcpt = strings.TrimSpace(rcpt)

	if rcpt != "RECEIVED" {
		err = errors.Newf(`Expected "RECEIVED" from client, got "%s" instead`, 500, rcpt)
	}

	return
}

func (s *ServerClientRef) checkAndValidateInstance() {
	if s.conn == nil || s.reader == nil {
		err := errors.Newf("ServerClientRef not properly instanciated")
		panic(err)
	}
}

func (s *ServerClientRef) Write(b []byte) (int, error) {
	s.checkAndValidateInstance()
	return s.conn.Write(b)
}

func (s *ServerClientRef) WriteString(d string) (int, error) {
	return s.Write([]byte(d))
}

func (s *ServerClientRef) Writeln(b []byte) (int, error) {
	return s.Write(append(b, '\n'))
}

func (s *ServerClientRef) WriteStringln(d string) (int, error) {
	return s.Writeln([]byte(d))
}

func (s *ServerClientRef) ReadBytes(delim byte) ([]byte, error) {
	s.checkAndValidateInstance()
	return s.reader.ReadBytes(delim)
}

func (s *ServerClientRef) ReadString(delim byte) (string, error) {
	s.checkAndValidateInstance()
	return s.reader.ReadString(delim)
}

func (s *ServerClientRef) writeResponse(output *Output) (err error) {
	data := []byte(ResponsePrefix + output.uuid + "::")

	json, err := objects.ToJSON(output)

	if err != nil {
		return
	}

	data = append(data, json...)

	_, err = s.Writeln(data)

	return
}

func (s *ServerClientRef) close() error {
	return s.conn.Close()
}

/**********************************************/

func newServerClientRef(conn net.Conn, opts ServerOptions) (value *ServerClientRef, err error) {
	tmpValue := &ServerClientRef{
		conn:   conn,
		opts:   opts,
		reader: bufio.NewReader(conn),
	}

	if tmpValue.uuid, err = uuid.Generate(); err != nil {
		return
	}

	if err = getServerClientRefName(tmpValue); err != nil {
		return
	}

	if err = getServerClientRefAlias(tmpValue); err != nil {
		return
	}

	if err = sendServerRefName(tmpValue); err != nil {
		return
	}

	if err = sendServerRefAlias(tmpValue); err != nil {
		return
	}

	value = tmpValue

	return
}

func getServerClientRefName(ref *ServerClientRef) (err error) {
	var input string

	if _, err = ref.WriteStringln("CLIENT NAME:"); err != nil {
		return
	}

	if input, err = ref.ReadString('\n'); err != nil {
		return
	}

	input = strings.TrimSpace(input)

	if input == "" {
		err = errors.New("Invalid client name")
	}

	ref.name = input

	return
}

func getServerClientRefAlias(ref *ServerClientRef) (err error) {
	var input string

	if _, err = ref.WriteStringln("CLIENT ALIAS:"); err != nil {
		return
	}

	if input, err = ref.ReadString('\n'); err != nil {
		return
	}

	input = strings.TrimSpace(input)

	if input == "" {
		err = errors.New("Invalid client alias")
	}

	ref.alias = input

	return
}

func sendServerRefName(ref *ServerClientRef) (err error) {
	var input string

	if input, err = ref.ReadString('\n'); err != nil {
		return
	}

	input = strings.TrimSpace(input)

	if input != "SERVER NAME:" {
		err = errors.Newf(`Expected "SERVER NAME:", got "%s" instead`, input)
		return
	}

	_, err = ref.WriteStringln(ref.opts.Name)

	return
}

func sendServerRefAlias(ref *ServerClientRef) (err error) {
	var input string

	if input, err = ref.ReadString('\n'); err != nil {
		return
	}

	input = strings.TrimSpace(input)

	if input != "SERVER ALIAS:" {
		err = errors.Newf(`Expected "SERVER NAME:", got "%s" instead`, input)
		return
	}

	_, err = ref.WriteStringln(ref.opts.Alias)

	return
}
