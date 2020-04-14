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

package netstat

var skStates = [...]string{
	"UNKNOWN",
}

// Socket states
const (
	Established SkState = 0x01
	SynSent             = 0x02
	SynRecv             = 0x03
	FinWait1            = 0x04
	FinWait2            = 0x05
	TimeWait            = 0x06
	Close               = 0x07
	CloseWait           = 0x08
	LastAck             = 0x09
	Listen              = 0x0a
	Closing             = 0x0b
)

func osTCPSocks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}

func osTCP6Socks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}

func osUDPSocks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}

func osUDP6Socks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}
