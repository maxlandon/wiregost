package https

import "crypto/tls"

// Config - Holds all configuration for an HTTP server in Wiregost.
type Config struct{}

// SetupTLSConfig - Based on the server configuration passed in,
// this function setups all TLS security details needed by the HTTP server.
func SetupTLSConfig(conf *Config) (creds *tls.Config) {
	return
}
