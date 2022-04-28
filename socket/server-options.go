package socket

import (
	"crypto/tls"
)

type ServerOptions struct {
	Host      string
	Port      string
	Name      string
	Alias     string
	TlsConfig *tls.Config
}

func (s ServerOptions) host() string {
	return s.Host
}

func (s ServerOptions) port() string {
	return s.Port
}

func (s *ServerOptions) name() string {
	return s.Name
}

func (s *ServerOptions) alias() *string {
	return &s.Alias
}

func (s *ServerOptions) tlsConfig() *tls.Config {
	return s.TlsConfig
}

func (s *ServerOptions) setTlsConfig(config *tls.Config) {
	s.TlsConfig = config
}
