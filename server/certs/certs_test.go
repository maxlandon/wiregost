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
	"bytes"
	"testing"
)

func TestGhostGenerateECCCertificate(t *testing.T) {
	GenerateCertificateAuthority(GhostCA)
	eccCert1, eccKey1, err := GhostGenerateECCCertificate("test1")
	if err != nil {
		t.Errorf("Failed to generate ecc certificate %v", err)
		return
	}
	eccCert2, eccKey2, err := GetCertificate(GhostCA, ECCKey, "test1")
	if err != nil {
		t.Errorf("Failed to get certificate %v", err)
		return
	}
	if !bytes.Equal(eccCert1, eccCert2) || !bytes.Equal(eccKey1, eccKey2) {
		t.Errorf("Stored ecc cert/key does match generated cert/key")
		return
	}

}

func TestGhostGenerateRSACertificate(t *testing.T) {
	rsaCert1, rsaKey1, err := GhostGenerateRSACertificate("test2")
	if err != nil {
		t.Errorf("Failed to generate rsa certificate %v", err)
		return
	}
	rsaCert2, rsaKey2, err := GetCertificate(GhostCA, RSAKey, "test2")
	if err != nil {
		t.Errorf("Failed to get certificate %v", err)
		return
	}
	if !bytes.Equal(rsaCert1, rsaCert2) || !bytes.Equal(rsaKey1, rsaKey2) {
		t.Errorf("Stored rsa cert/key does match generated cert/key")
		return
	}
}

func TestOperatorGenerateCertificate(t *testing.T) {
	GenerateCertificateAuthority(UserCA)
	cert1, key1, err := OperatorClientGenerateCertificate("test3")
	if err != nil {
		t.Errorf("Failed to store ecc certificate %v", err)
		return
	}

	cert2, key2, err := OperatorClientGetCertificate("test3")
	if err != nil {
		t.Errorf("Failed to get ecc certificate %v", err)
		return
	}

	if !bytes.Equal(cert1, cert2) || !bytes.Equal(key1, key2) {
		t.Errorf("Stored ecc cert/key does match generated cert/key: %v != %v", cert1, cert2)
		return
	}
}

func TestGenerateServerCertificate(t *testing.T) {
	GenerateCertificateAuthority(ServerCA)
	ServerGenerateECCCertificate("test3.com")
	_, _, err := ServerGenerateRSACertificate("test3.com")
	if err != nil {
		t.Errorf("Failed to generate server rsa certificate")
		return
	}
}
