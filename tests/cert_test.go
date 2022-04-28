package tests

import (
	"testing"

	_ "cyberpull.com/go-cyb/env"

	"cyberpull.com/go-cyb/cert"
)

func TestCert_GetCertificates(t *testing.T) {
	_, err := cert.GetCertificates()

	if err != nil {
		t.Error(err)
	}
}

func TestCert_GetTlsConfig(t *testing.T) {
	_, err := cert.GetTLSConfig()

	if err != nil {
		t.Error(err)
	}
}
