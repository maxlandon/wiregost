// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/evilsocket/islazy/tui"

	clientassets "github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	userLog = log.ServerLogger("users", "setup")
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

// NewUser - Add a User to DB
func NewUser(user, lhost string, lport uint32, isDefault bool) error {

	regex, _ := regexp.Compile("[^A-Za-z0-9]+") // Only allow alphanumeric chars
	user = regex.ReplaceAllString(user, "")

	publicKey, privateKey, err := certs.UserClientGenerateCertificate(user)
	if err != nil {
		return errors.New(tui.Red("Failed to generate default user certificate"))
	}

	caCertPEM, _, _ := certs.GetCertificateAuthorityPEM(certs.UserCA)
	config := ClientConfig{
		User:          user,
		LHost:         lhost,
		LPort:         int(lport),
		CACertificate: string(caCertPEM),
		PrivateKey:    string(privateKey),
		Certificate:   string(publicKey),
		IsDefault:     isDefault,
	}

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
		filename := fmt.Sprintf("%s_%s.cfg", filepath.Base(user), filepath.Base(lhost))
		saveTo = filepath.Join(saveTo, filename)
	}
	err = ioutil.WriteFile(saveTo, configJSON, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write config to: %s (%v) \n", saveTo, err)
	}

	return nil
}

// CreateDefaultUser - Check if users exist in db, if not create default one and save config
func CreateDefaultUser() error {

	userCerts := certs.UserClientListCertificates()

	if len(userCerts) == 0 {
		err := NewUser("wiregost", "localhost", uint32(1708), true)
		if err != nil {
			return err
		}
		userLog.Infoln("Created default Wiregost user config: name -> 'wiregost', lhost -> 'localhost', lport -> 1708")

		filename := fmt.Sprintf("%s_%s.cfg", filepath.Base("wiregost"), filepath.Base("localhost"))
		moveFrom := filepath.Join(assets.GetRootAppDir(), "users", filename)
		moveTo := filepath.Join(clientassets.GetConfigDir(), filename)

		err = os.Rename(moveFrom, moveTo)
		if err != nil {
			return errors.New((tui.Red("Failed to move default user config to ~/.wiregost-client/config directory")))
		}
	}

	return nil
}
