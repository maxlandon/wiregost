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

import "errors"

var (
	// ErrNmapNotInstalled means that upon trying to manually locate nmap in the user's path,
	// it was not found. Either use the WithBinaryPath method to set it manually, or make sure that
	// the nmap binary is present in the user's $PATH.
	ErrNmapNotInstalled = errors.New("nmap binary was not found")

	// ErrScanTimeout means that the provided context was done before the scanner finished its scan.
	ErrScanTimeout = errors.New("nmap scan timed out")

	// ErrMallocFailed means that nmap crashed due to insufficient memory, which may happen on large target networks.
	ErrMallocFailed = errors.New("malloc failed, probably out of space")

	// ErrParseOutput means that nmap's output was not parsed successfully.
	ErrParseOutput = errors.New("unable to parse nmap output, see warnings for details")

	// ErrResolveName means that Nmap could not resolve a name.
	ErrResolveName = errors.New("nmap could not resolve a name")
)
