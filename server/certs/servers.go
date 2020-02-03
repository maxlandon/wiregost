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

const (
	// ServerCA - Directory containing server certificates
	ServerCA = "server"
)

// ServerGenerateECCCertificate - Generate a server certificate signed with a given CA
func ServerGenerateECCCertificate(host string) ([]byte, []byte, error) {
	cert, key := GenerateECCCertificate(ServerCA, host, false, false)
	err := SaveCertificate(ServerCA, ECCKey, host, cert, key)
	return cert, key, err
}

// ServerGenerateRSACertificate - Generate a server certificate signed with a given CA
func ServerGenerateRSACertificate(host string) ([]byte, []byte, error) {
	cert, key := GenerateRSACertificate(ServerCA, host, false, false)
	err := SaveCertificate(ServerCA, RSAKey, host, cert, key)
	return cert, key, err
}
