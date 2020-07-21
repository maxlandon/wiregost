package clients

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/maxlandon/wiregost/server/certs"
)

// LoadUserServerTLSConfig - Load TLS infrastructure for interacting with client consoles.
func LoadUserServerTLSConfig(host string) (config *tls.Config) {

	caCertPtr, _, err := certs.GetCertificateAuthority(certs.UserCA)
	if err != nil {
		// mtlsLog.Fatalf("Invalid ca type (%s): %v", certs.OperatorCA, host)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCertPtr)

	_, _, err = certs.UserServerGetCertificate(host)
	if err == certs.ErrCertDoesNotExist {
		certs.UserServerGenerateCertificate(host)
	}

	certPEM, keyPEM, err := certs.UserServerGetCertificate(host)
	if err != nil {
		// mtlsLog.Errorf("Failed to generate or fetch certificate %s", err)
		return nil
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		// mtlsLog.Fatalf("Error loading server certificate: %v", err)
	}

	config = &tls.Config{
		RootCAs:                  caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                caCertPool,
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS13,
	}
	config.BuildNameToCertificate()

	return
}
