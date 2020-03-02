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
	"net/url"
	"os"
	"path"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/log"
)

var (
	buildLog = log.ServerLogger("generate", "build")
	// Fix #67: use an arch specific compiler
	defaultMingwPath = map[string]string{
		"386":   "/usr/bin/i686-w64-mingw32-gcc",
		"amd64": "/usr/bin/x86_64-w64-mingw32-gcc",
	}
)

const (
	// WINDOWS OS
	WINDOWS = "windows"

	// DARWIN / MacOS
	DARWIN = "darwin"

	// LINUX OS
	LINUX = "linux"

	clientsDirName = "clients"
	ghostsDirName  = "ghosts"

	encryptKeySize = 16

	// DefaultReconnectInterval - In seconds
	DefaultReconnectInterval = 60
	// DefaultMTLSLPort - Default listen port
	DefaultMTLSLPort = 8888
	// DefaultHTTPLPort - Default HTTP listen port
	DefaultHTTPLPort = 443 // Assume SSL, it'll fallback

	// GhostCC64EnvVar - Environment variable that can specify the 64 bit mingw path
	GhostCC64EnvVar = "SLIVER_CC_64"
	// GhostCC32EnvVar - Environment variable that can specify the 32 bit mingw path
	GhostCC32EnvVar = "SLIVER_CC_32"
)

// GhostConfig - Parameters when generating a implant
type GhostConfig struct {
	// Go
	GOOS   string `json:"go_os"`
	GOARCH string `json:"go_arch"`

	// Standard
	Name                string `json:"name"`
	CACert              string `json:"ca_cert"`
	Cert                string `json:"cert"`
	Key                 string `json:"key"`
	Debug               bool   `json:"debug"`
	ObfuscateSymbols    bool   `json:"obfuscate_symbols"`
	ReconnectInterval   int    `json:"reconnect_interval"`
	MaxConnectionErrors int    `json:"max_connection_errors"`

	C2            []GhostC2 `json:"c2s"`
	MTLSc2Enabled bool      `json:"c2_mtls_enabled"`
	HTTPc2Enabled bool      `json:"c2_http_enabled"`
	DNSc2Enabled  bool      `json:"c2_dns_enabled"`
	CanaryDomains []string  `json:"canary_domains"`

	// Limits
	LimitDomainJoined bool   `json:"limit_domainjoined"`
	LimitHostname     string `json:"limit_hostname"`
	LimitUsername     string `json:"limit_username"`
	LimitDatetime     string `json:"limit_datetime"`

	// Output Format
	Format clientpb.GhostConfig_OutputFormat `json:"format"`

	// For 	IsSharedLib bool `json:"is_shared_lib"`
	IsSharedLib bool `json:"is_shared_lib"`

	FileName string

	// Added for DB
	WorkspaceID uint
}

// ToProtobuf - Convert SliverConfig to protobuf equiv
func (c *GhostConfig) ToProtobuf() *clientpb.GhostConfig {
	config := &clientpb.GhostConfig{
		GOOS:             c.GOOS,
		GOARCH:           c.GOARCH,
		Name:             c.Name,
		CACert:           c.CACert,
		Cert:             c.Cert,
		Key:              c.Key,
		Debug:            c.Debug,
		ObfuscateSymbols: c.ObfuscateSymbols,
		CanaryDomains:    c.CanaryDomains,

		ReconnectInterval:   uint32(c.ReconnectInterval),
		MaxConnectionErrors: uint32(c.MaxConnectionErrors),

		LimitDatetime:     c.LimitDatetime,
		LimitDomainJoined: c.LimitDomainJoined,
		LimitHostname:     c.LimitHostname,
		LimitUsername:     c.LimitUsername,

		IsSharedLib: c.IsSharedLib,
		Format:      c.Format,

		FileName:    c.FileName,
		WorkspaceID: uint32(c.WorkspaceID),
	}
	config.C2 = []*clientpb.GhostC2{}
	for _, c2 := range c.C2 {
		config.C2 = append(config.C2, c2.ToProtobuf())
	}
	return config
}

// GhostConfigFromProtobuf - Create a native config struct from Protobuf
func GhostConfigFromProtobuf(pbConfig *clientpb.GhostConfig) *GhostConfig {
	cfg := &GhostConfig{}

	cfg.GOOS = pbConfig.GOOS
	cfg.GOARCH = pbConfig.GOARCH
	cfg.Name = pbConfig.Name
	cfg.CACert = pbConfig.CACert
	cfg.Cert = pbConfig.Cert
	cfg.Key = pbConfig.Key
	cfg.Debug = pbConfig.Debug
	cfg.ObfuscateSymbols = pbConfig.ObfuscateSymbols
	cfg.CanaryDomains = pbConfig.CanaryDomains

	cfg.ReconnectInterval = int(pbConfig.ReconnectInterval)
	cfg.MaxConnectionErrors = int(pbConfig.MaxConnectionErrors)

	cfg.LimitDomainJoined = pbConfig.LimitDomainJoined
	cfg.LimitDatetime = pbConfig.LimitDatetime
	cfg.LimitUsername = pbConfig.LimitUsername
	cfg.LimitHostname = pbConfig.LimitHostname

	cfg.Format = pbConfig.Format
	cfg.IsSharedLib = pbConfig.IsSharedLib

	cfg.C2 = copyC2List(pbConfig.C2)
	cfg.MTLSc2Enabled = isC2Enabled([]string{"mtls"}, cfg.C2)
	cfg.HTTPc2Enabled = isC2Enabled([]string{"http", "https"}, cfg.C2)
	cfg.DNSc2Enabled = isC2Enabled([]string{"dns"}, cfg.C2)

	cfg.FileName = pbConfig.FileName
	cfg.WorkspaceID = uint(pbConfig.WorkspaceID)
	return cfg
}

func copyC2List(src []*clientpb.GhostC2) []GhostC2 {
	c2s := []GhostC2{}
	for _, srcC2 := range src {
		c2URL, err := url.Parse(srcC2.URL)
		if err != nil {
			buildLog.Warnf("Failed to parse c2 url %v", err)
			continue
		}
		c2s = append(c2s, GhostC2{
			Priority: srcC2.Priority,
			URL:      c2URL.String(),
			Options:  srcC2.Options,
		})
	}
	return c2s
}

func isC2Enabled(schemes []string, c2s []GhostC2) bool {
	for _, c2 := range c2s {
		c2URL, err := url.Parse(c2.URL)
		if err != nil {
			buildLog.Warnf("Failed to parse c2 url %v", err)
			continue
		}
		for _, scheme := range schemes {
			if scheme == c2URL.Scheme {
				return true
			}
		}
	}
	buildLog.Debugf("No %v URLs found in %v", schemes, c2s)
	return false
}

// GhostC2 - C2 struct
type GhostC2 struct {
	Priority uint32 `json:"priority"`
	URL      string `json:"url"`
	Options  string `json:"options"`
}

// ToProtobuf - Convert to protobuf version
func (s GhostC2) ToProtobuf() *clientpb.GhostC2 {
	return &clientpb.GhostC2{
		Priority: s.Priority,
		URL:      s.URL,
		Options:  s.Options,
	}
}

func (s GhostC2) String() string {
	return s.URL
}

// GetGhostsDir - Get the binary directory
func GetGhostsDir() string {
	appDir := assets.GetRootAppDir()
	ghostsDir := path.Join(appDir, ghostsDirName)
	if _, err := os.Stat(ghostsDir); os.IsNotExist(err) {
		buildLog.Infof("Creating bin directory: %s", ghostsDir)
		err = os.MkdirAll(ghostsDir, 0700)
		if err != nil {
			buildLog.Fatal(err)
		}
	}
	return ghostsDir
}

// GhostEgg - Generates a sliver egg (stager) binary
func GhostEgg(config GhostConfig) (string, error) {

	return "", nil
}
