package socket

import (
	"crypto/tls"
	"strings"

	"cyberpull.com/go-cyb/cert"
	"cyberpull.com/go-cyb/errors"
)

type Options interface {
	host() string
	port() string
	name() string
	alias() *string
	tlsConfig() *tls.Config
	setTlsConfig(config *tls.Config)
}

func address(opts Options) string {
	return opts.host() + ":" + opts.port()
}

func sanitizeNameAndAlias(opts Options) (err error) {
	if opts.name() == "" {
		err = errors.New(`"Name" is required`)
		return
	}

	if opts.alias() == nil || *opts.alias() == "" {
		*opts.alias() = opts.name()
	}

	*opts.alias() = strings.ToLower(*opts.alias())
	*opts.alias() = strings.TrimSpace(*opts.alias())
	*opts.alias() = strings.ReplaceAll(*opts.alias(), " ", "-")

	return
}

func sanitizeTlsConfig(opts Options, forServer ...bool) (err error) {
	if cert.IsEnabled() {
		var config *tls.Config

		if config, err = cert.GetTLSConfig(forServer...); err != nil {
			return
		}

		if opts.tlsConfig() == nil {
			opts.setTlsConfig(config)
		}

		if opts.tlsConfig().Certificates == nil || len(opts.tlsConfig().Certificates) == 0 {
			opts.tlsConfig().Certificates, err = cert.GetCertificates()
			cert.SanitizeTlsConfig(opts.tlsConfig(), forServer...)
		}
	}

	return
}
