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
	// GhostCA - Directory containing sliver certificates
	GhostCA = "ghost"
)

// GhostGenerateECCCertificate - Generate a certificate signed with a given CA
func GhostGenerateECCCertificate(ghostName string) ([]byte, []byte, error) {
	cert, key := GenerateECCCertificate(GhostCA, ghostName, false, true)
	err := SaveCertificate(GhostCA, ECCKey, ghostName, cert, key)
	return cert, key, err
}

// GhostGenerateRSACertificate - Generate a certificate signed with a given CA
func GhostGenerateRSACertificate(ghostName string) ([]byte, []byte, error) {
	cert, key := GenerateRSACertificate(GhostCA, ghostName, false, true)
	err := SaveCertificate(GhostCA, RSAKey, ghostName, cert, key)
	return cert, key, err
}
