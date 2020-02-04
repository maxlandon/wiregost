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

package certs

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/maxlandon/wiregost/server/db"
)

const (
	// UserCA - Directory containing user certificates
	UserCA = "user"

	clientNamespace = "client"
	serverNamespace = "server"
)

// UserClientGenerateCertificate - Generate a certificate signed with a given CA
func UserClientGenerateCertificate(user string) ([]byte, []byte, error) {
	cert, key := GenerateECCCertificate(UserCA, user, false, true)
	err := SaveCertificate(UserCA, ECCKey, fmt.Sprintf("%s.%s", clientNamespace, user), cert, key)
	return cert, key, err
}

// UserClientGetCertificate - Helper function to fetch a client cert
func UserClientGetCertificate(user string) ([]byte, []byte, error) {
	return GetCertificate(UserCA, ECCKey, fmt.Sprintf("%s.%s", clientNamespace, user))
}

// UserServerGetCertificate - Helper function to fetch a client cert
func UserServerGetCertificate(user string) ([]byte, []byte, error) {
	return GetCertificate(UserCA, ECCKey, fmt.Sprintf("%s.%s", serverNamespace, user))
}

// UserServerGenerateCertificate - Generate a certificate signed with a given CA
func UserServerGenerateCertificate(hostname string) ([]byte, []byte, error) {
	cert, key := GenerateECCCertificate(UserCA, hostname, false, false)
	err := SaveCertificate(UserCA, ECCKey, fmt.Sprintf("%s.%s", serverNamespace, hostname), cert, key)
	return cert, key, err
}

// UserClientListCertificates - Get all client certificates
func UserClientListCertificates() []*x509.Certificate {
	bucket, err := db.GetBucket(UserCA)
	if err != nil {
		return []*x509.Certificate{}
	}

	// The key structure is: <key type>_<namespace>.<user name>
	users, err := bucket.List(fmt.Sprintf("%s_%s", ECCKey, clientNamespace))
	if err != nil {
		return []*x509.Certificate{}
	}
	certsLog.Infof("Found %d user certs ...", len(users))

	certs := []*x509.Certificate{}
	for _, user := range users {

		certsLog.Infof("User = %v", user)
		keypairRaw, err := bucket.Get(user)
		if err != nil {
			certsLog.Warnf("Failed to fetch user keypair %v", err)
			continue
		}
		keypair := &CertificateKeyPair{}
		json.Unmarshal(keypairRaw, keypair)

		block, _ := pem.Decode(keypair.Certificate)
		if block == nil {
			certsLog.Warn("failed to parse certificate PEM")
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			certsLog.Warnf("failed to parse x.509 certificate %v", err)
			continue
		}
		certs = append(certs, cert)
	}
	return certs
}
