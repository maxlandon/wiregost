package certs

import (
	"os"
	"path"

	"golang.org/x/crypto/acme/autocert"
	// "github.com/maxlandon/wiregost/server/log"
)

const (
	// ACMEDirName - Name of dir to store ACME certs
	ACMEDirName = "acme"
)

var (
// acmeLog = log.ServerLogger("certs", "acme")
)

// GetACMEDir - Dir to store ACME certs
func GetACMEDir() string {
	acmePath := path.Join(getCertDir(), ACMEDirName)
	if _, err := os.Stat(acmePath); os.IsNotExist(err) {
		// acmeLog.Infof("[mkdir] %s", acmePath)
		os.MkdirAll(acmePath, os.ModePerm)
	}
	return acmePath
}

// GetACMEManager - Get an ACME cert/tls config with the certs
func GetACMEManager(domain string) *autocert.Manager {
	acmeDir := GetACMEDir()
	return &autocert.Manager{
		Cache:  autocert.DirCache(acmeDir),
		Prompt: autocert.AcceptTOS,
	}
}
