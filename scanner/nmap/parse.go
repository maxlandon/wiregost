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

package nmap

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/maxlandon/wiregost/data-service/models"
)

// Timestamp represents time as a UNIX timestamp in seconds.
type Timestamp time.Time

// ParseTime converts a UNIX timestamp string to a time.Time.
func (t *Timestamp) ParseTime(s string) error {
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(timestamp, 0))

	return nil
}

// Parse takes a byte array of nmap xml data and unmarshals it into a Run struct.
func Parse(content []byte) (*models.Run, error) {
	r := &models.Run{
		RawXML: content,
	}

	err := xml.Unmarshal(content, r)

	return r, err
}
