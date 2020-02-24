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

package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"gopkg.in/yaml.v2"

	"github.com/maxlandon/wiregost/data_service/assets"
	"github.com/maxlandon/wiregost/data_service/models"
)

// Env contains all configuration options for DB access and data web service. It passes its DB connection
// pool to HTTP handlers that need it, and its HTTP service parameters to remote/ functions.
type Env struct {
	DB *models.DB
	// Database
	Database *Database
	// Web service
	Service *Service
}

// Database
type Database struct {
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
}

// Web service
type Service struct {
	Address     string `yaml:"address"`
	Port        int    `yaml:"port"`
	URL         string `yaml:"url"`
	Certificate string `yaml:"certificate"`
	Key         string `yaml:"key"`
}

// LoadEnv - Load Data Service configuration
func LoadEnv() *Env {
	// Create a default Data Service config, eventually parse one if found
	env := &Env{
		Database: &Database{
			DbName:     "wiregost_db",
			DbUser:     "wiregost",
			DbPassword: "wiregost",
		},
		Service: &Service{
			Address:     "localhost",
			Port:        8001,
			URL:         "/",
			Certificate: "~/.wiregost/data-service/certs/wiregost_pub.pem",
			Key:         "~/.wiregost/data-service/certs/wiregost_priv.pem",
		},
	}

	// Load config
	file := filepath.Join(assets.GetDataServiceDir(), "config.yaml")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = saveDataServiceEnv(env)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		err = saveDataServiceEnv(env)
	}

	err = yaml.Unmarshal(data, &env)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to unmarshal config.yaml file."))
	}

	// Check certificate and key exist
	cert := filepath.Join(assets.GetDataServiceDir(), "certs", "wiregost_pub.pem")

	// If not, generate them
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		err = createDataServiceCertificates()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}

	// Load them
	env.Service.Certificate, err = fs.Expand(env.Service.Certificate)
	env.Service.Key, err = fs.Expand(env.Service.Key)

	// Connect to postgreSQL
	env.DB, err = models.New(env.Database.DbName, env.Database.DbUser, env.Database.DbPassword)
	if err != nil {
		fmt.Println(err.Error())
	}

	return env
}

// saveDataServiceEnv either saves the current Data Service Env config, or writes a default one.
func saveDataServiceEnv(env *Env) error {
	saveTo := assets.GetDataServiceDir()
	envYAML, _ := yaml.Marshal(env)

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write to wiregost-client root directory %s", err))
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "config.yaml"
		saveTo = filepath.Join(saveTo, filename)
	}

	err = ioutil.WriteFile(saveTo, envYAML, 0600)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write config to: %s (%v) \n", saveTo, err))
	}

	f, err := os.OpenFile(saveTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	f.Close()

	return nil
}

func createDataServiceCertificates() error {
	certsDir := filepath.Join(assets.GetDataServiceDir(), "certs")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := key.PublicKey

	// RSA PRIVATE KEY
	outFile, err := os.Create(filepath.Join(certsDir, "wiregost_priv.pem"))
	if err != nil {
		return err
	}

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		return err
	}
	outFile.Close()

	// CERTIFICATE
	//Generate cryptographically strong pseudo-random between 0 - max
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)

	template := &x509.Certificate{
		BasicConstraintsValid: true,
		SerialNumber:          n,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &publicKey, key)
	if err != nil {
		fmt.Println(err)
	}

	var certPem = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	pemfile, err := os.Create(filepath.Join(certsDir, "wiregost_pub.pem"))
	if err != nil {
		return err
	}

	err = pem.Encode(pemfile, certPem)
	if err != nil {
		return err
	}
	pemfile.Close()

	return nil
}
