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

package util

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func randomData() []byte {
	buf := make([]byte, 128)
	rand.Read(buf)
	return buf
}

func TestGzip(t *testing.T) {
	sample := randomData()
	gzipData := bytes.NewBuffer([]byte{})
	gz := new(Gzip)
	gz.Encode(gzipData, sample)
	data, err := gz.Decode(gzipData.Bytes())
	if err != nil {
		t.Errorf("gzip decode returned an error %v", err)
		return
	}
	if !bytes.Equal(sample, data) {
		t.Errorf("sample does not match returned\n%#v != %#v", sample, data)
	}
}
