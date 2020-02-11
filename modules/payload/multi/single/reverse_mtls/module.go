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

// CHANGE THE NAME OF THE PACKAGE WITH THE NAME OF YOUR MODULE/DIRECTORY
package reverse_mtls

import (
	"fmt"
	"path/filepath"
	"strconv"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	c2 "github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// metadataFile - Full path to module metadata
var metadataFile = filepath.Join(assets.GetModulesDir(), "payload/multi/single/reverse_mtls/metadata.json")

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseMTLS - A single stage MTLS implant
type ReverseMTLS struct {
	Base *templates.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseMTLS {
	return &ReverseMTLS{Base: &templates.Module{}}
}

// Init - Module initialization, loads metadata. ** DO NOT ERASE **
func (s *ReverseMTLS) Init() error {
	return s.Base.Init(metadataFile)
}

// ToProtobuf - Returns protobuf version of module
func (s *ReverseMTLS) ToProtobuf() *pb.Module {
	return s.Base.ToProtobuf()
}

// SetOption - Sets a module option through its base object.
func (s *ReverseMTLS) SetOption(option, name string) {
	s.Base.SetOption(option, name)
}

// [ Module Methods ] ------------------------------------------------------------------------//

var rpcLog = log.ServerLogger("rpc", "server")

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseMTLS) Run(command string) (result string, err error) {

	switch command {

	case "to_listener":

		host := s.Base.Options["LHost"].Value
		portUint, _ := strconv.Atoi(s.Base.Options["LPort"].Value)
		port := uint16(portUint)

		ln, err := c2.StartMutualTLSListener(host, port)
		if err != nil {
			return "", err
		}

		job := &core.Job{
			ID:          core.GetJobID(),
			Name:        "mTLS",
			Description: "Mutual TLS listener",
			Protocol:    "tcp",
			Port:        port,
			JobCtrl:     make(chan bool),
		}

		go func() {
			<-job.JobCtrl
			rpcLog.Infof("Stopping mTLS listener (%d) ...", job.ID)
			ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

			core.Jobs.RemoveJob(job)

			core.EventBroker.Publish(core.Event{
				Job:       job,
				EventType: consts.StoppedEvent,
			})
		}()

		core.Jobs.AddJob(job)

		return fmt.Sprintf("Reverse Mutual TLS listener started at %s:%d", host, port), nil
	}

	return "ReverseMTLS listener started", nil
}
