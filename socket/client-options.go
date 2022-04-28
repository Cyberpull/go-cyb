package socket

import "crypto/tls"

type ClientOptions struct {
	ServerHost string
	ServerPort string
	Name       string
	Alias      string
	TlsConfig  *tls.Config
}

func (s ClientOptions) host() string {
	return s.ServerHost
}

func (s ClientOptions) port() string {
	return s.ServerPort
}

func (s *ClientOptions) name() string {
	return s.Name
}

func (s *ClientOptions) alias() *string {
	return &s.Alias
}

func (s *ClientOptions) tlsConfig() *tls.Config {
	return s.TlsConfig
}

func (s *ClientOptions) setTlsConfig(config *tls.Config) {
	s.TlsConfig = config
}
