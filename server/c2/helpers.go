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

package c2

import (
	secureRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"net"
	"strings"
)

func fingerprintSHA256(block *pem.Block) string {
	hash := sha256.Sum256(block.Bytes)
	b64hash := base64.RawStdEncoding.EncodeToString(hash[:])
	return strings.TrimRight(b64hash, "=")
}

// TODO: Add a geofilter to make it look like we're in various regions of the world.
// We don't really need to use srand but it's just easier to operate on bytes here.
func randomIP() net.IP {
	randBuf := make([]byte, 4)
	secureRand.Read(randBuf)

	// Ensure non-zeros with various bitmasks
	return net.IPv4(randBuf[0]|0x10, randBuf[1]|0x10, randBuf[2]|0x1, randBuf[3]|0x10)
}
