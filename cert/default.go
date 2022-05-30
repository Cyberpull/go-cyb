package cert

import (
	"crypto/tls"
	"os"
	"strings"

	_ "cyberpull.com/go-cyb/env"

	"cyberpull.com/go-cyb/errors"
)

var (
	certEnabled string
	certEnv     string
	certFile    string
	keyFile     string
)

func init() {
	certEnabled = os.Getenv("CERT_ENABLED")
	certEnabled = strings.ToLower(certEnabled)

	if IsEnabled() {
		certEnv = os.Getenv("CERT_ENV")
		certEnv = strings.ToLower(certEnv)

		certFile = os.Getenv("CERT_CRT_FILE")
		keyFile = os.Getenv("CERT_KEY_FILE")
	}
}

func validate() (err error) {
	if certFile == "" || keyFile == "" {
		err = errors.New(`"CERT_CRT_FILE" and "CERT_KEY_FILE" environment variables are required`)
	}

	return
}

func IsEnabled() bool {
	return certEnabled == "yes"
}

func IsLocal() bool {
	return certEnv == "local"
}

func GetTLSConfig(forServer ...bool) (config *tls.Config, err error) {
	config = &tls.Config{}

	if IsEnabled() {
		config.Certificates, err = GetCertificates()
		err = SanitizeTlsConfig(config, forServer...)
	}

	return
}

func SanitizeTlsConfig(config *tls.Config, forServer ...bool) (err error) {
	if config == nil {
		err = errors.New("No config found.")
		return
	}

	if len(forServer) == 0 {
		forServer = append(forServer, false)
	}

	if !forServer[0] && IsLocal() {
		config.InsecureSkipVerify = true
		config.VerifyPeerCertificate = nil
		config.VerifyConnection = nil
	}

	return
}

func GetCertificates() (value []tls.Certificate, err error) {
	if err = validate(); err != nil {
		return
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return
	}

	value = []tls.Certificate{cert}

	return
}
