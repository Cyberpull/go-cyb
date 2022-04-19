package cert

import (
	"crypto/tls"
	"os"
	"strings"

	_ "cyberpull.com/go-cyb/env"
)

var (
	certEnv  string
	certFile string
	keyFile  string
)

func init() {
	certEnv = os.Getenv("CERT_ENV")
	certEnv = strings.ToLower(certEnv)

	certFile = os.Getenv("CERT_CRT_FILE")
	keyFile = os.Getenv("CERT_KEY_FILE")

	if certFile == "" {
		panic(`""CERT_CRT_FILE" environment variable is required`)
	}

	if keyFile == "" {
		panic(`""CERT_KEY_FILE" environment variable is required`)
	}
}

func GetTLSConfig() (config *tls.Config, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	config = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if certEnv == "local" {
		config.InsecureSkipVerify = true
		config.VerifyPeerCertificate = nil
		config.VerifyConnection = nil
	}

	return
}
