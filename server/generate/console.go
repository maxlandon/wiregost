package generate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	db "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/certs"
)

// ClientConfig - Struct containing user information and certificate
type ClientConfig struct {
	User          string `json:"user"`
	LHost         string `json:"lhost"`
	LPort         int    `json:"lport"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	Certificate   string `json:"certificate"`
	IsDefault     bool   `json:"is_default"`
}

// CompileObfuscatedConsole - Each user can compile obfuscated binaries with built-in certificates and tokens.
func CompileObfuscatedConsole(user string, name string, send bool) (err error) {
	return
}

// UserConsoleConfig - If user needs a file with credentials for connecting to this server
func UserConsoleConfig(user *db.User, pub []byte, priv []byte, isDefault bool) (err error) {

	// Make config
	caCertPEM, _, _ := certs.GetCertificateAuthorityPEM(certs.UserCA)
	config := ClientConfig{
		User:          user.Name,
		LHost:         assets.ServerConfiguration.ServerHost,
		LPort:         int(assets.ServerConfiguration.ServerPort),
		CACertificate: string(caCertPEM),
		PrivateKey:    string(priv),
		Certificate:   string(pub),
		IsDefault:     isDefault,
	}

	// Save to file
	configJSON, _ := json.Marshal(config)

	saveTo, _ := filepath.Abs(path.Join(assets.GetRootAppDir(), "users"))

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Cannot write to wiregost root directory %s", err)
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		def := ""
		if isDefault {
			def = "default"
		} else {
			def = "normal"
		}
		filename := fmt.Sprintf("%s_%s_%s.cfg", filepath.Base(user.Name), filepath.Base(assets.ServerConfiguration.ServerHost), def)
		saveTo = filepath.Join(saveTo, filename)
	}
	err = ioutil.WriteFile(saveTo, configJSON, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write config to: %s (%v) \n", saveTo, err)
	}

	return
}
