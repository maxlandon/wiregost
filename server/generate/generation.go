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
	"errors"
	"fmt"
	"path"

	consts "github.com/maxlandon/wiregost/client/constants"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/assets"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/gogo"
)

// GhostSharedLibrary - Generates a ghost shared library (DLL/dylib/so) binary
func GhostSharedLibrary(config *GhostConfig) (string, error) {
	// Compile go code
	appDir := assets.GetRootAppDir()
	crossCompiler := getCCompiler(config.GOARCH)
	if crossCompiler == "" {
		return "", errors.New("No cross-compiler (mingw) found")
	}
	goConfig := &gogo.GoConfig{
		CGO:    "1",
		CC:     crossCompiler,
		GOOS:   config.GOOS,
		GOARCH: config.GOARCH,
		GOROOT: gogo.GetGoRootDir(appDir),
	}
	pkgPath, err := renderGhostGoCode(config, goConfig)
	if err != nil {
		return "", err
	}

	dest := path.Join(goConfig.GOPATH, "bin", config.Name)
	if goConfig.GOOS == WINDOWS {
		dest += ".dll"
	}
	if goConfig.GOOS == DARWIN {
		dest += ".dylib"
	}
	if goConfig.GOOS == LINUX {
		dest += ".so"
	}

	tags := []string{"netgo"}
	ldflags := []string{"-s -w -buildid="}
	if !config.Debug && goConfig.GOOS == WINDOWS {
		ldflags[0] += " -H=windowsgui"
	}
	// Keep those for potential later use
	gcflags := fmt.Sprintf("")
	asmflags := fmt.Sprintf("")
	// trimpath is now a separate flag since Go 1.13
	trimpath := "-trimpath"
	_, err = gogo.GoBuild(*goConfig, pkgPath, dest, "c-shared", tags, ldflags, gcflags, asmflags, trimpath)
	config.FileName = path.Base(dest)
	saveFileErr := GhostFileSave(config.Name, dest)
	saveCfgErr := GhostConfigSave(config)
	if saveFileErr != nil || saveCfgErr != nil {
		buildLog.Errorf("Failed to save file to db %s %s", saveFileErr, saveCfgErr)
	}
	return dest, err
}

// GhostExecutable - Generates a ghost executable binary
func GhostExecutable(config *GhostConfig) (string, error) {

	// Compile go code
	appDir := assets.GetRootAppDir()
	cgo := "0"
	if config.IsSharedLib {
		cgo = "1"
	}
	goConfig := &gogo.GoConfig{
		CGO:    cgo,
		GOOS:   config.GOOS,
		GOARCH: config.GOARCH,
		GOROOT: gogo.GetGoRootDir(appDir),
	}
	pkgPath, err := renderGhostGoCode(config, goConfig)
	if err != nil {
		return "", err
	}

	dest := path.Join(goConfig.GOPATH, "bin", config.Name)
	if goConfig.GOOS == WINDOWS {
		dest += ".exe"
	}
	tags := []string{"netgo"}
	ldflags := []string{"-s -w -buildid="}
	if !config.Debug && goConfig.GOOS == WINDOWS {
		ldflags[0] += " -H=windowsgui"
	}
	gcflags := fmt.Sprintf("")
	asmflags := fmt.Sprintf("")
	// trimpath is now a separate flag since Go 1.13
	trimpath := "-trimpath"
	_, err = gogo.GoBuild(*goConfig, pkgPath, dest, "", tags, ldflags, gcflags, asmflags, trimpath)
	config.FileName = path.Base(dest)
	saveFileErr := GhostFileSave(config.Name, dest)
	saveCfgErr := GhostConfigSave(config)
	if saveFileErr != nil || saveCfgErr != nil {
		buildLog.Errorf("Failed to save file to db %s %s", saveFileErr, saveCfgErr)
	}
	return dest, err
}

// CompileGhost concurrently compiles a ghost implant with the provided config
func CompileGhost(config GhostConfig) {
	c2s := []string{}
	for _, c := range config.C2 {
		c2s = append(c2s, c.String())
	}
	description := fmt.Sprintf("Platform: %s/%s - Type: %s => %v", config.GOOS, config.GOARCH, config.Format, c2s)

	// Send job start
	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "Compiler",
		Description: description,
		JobCtrl:     make(chan bool),
	}

	go func() {
		<-job.JobCtrl
		buildLog.Infof("Done compiling ghost implant %s", config.Name)
		core.Jobs.RemoveJob(job)

		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
		})
	}()

	core.Jobs.AddJob(job)

	// Compile according to Format
	switch config.Format {
	case clientpb.GhostConfig_EXECUTABLE:
		go func() {
			path, err := GhostExecutable(&config)
			if err != nil {
				job.Err = err.Error()
				job.JobCtrl <- true
			}
			job.Err = path
			job.JobCtrl <- true
		}()

	case clientpb.GhostConfig_SHARED_LIB:
		go func() {
			path, err := GhostSharedLibrary(&config)
			if err != nil {
				job.Err = err.Error()
				job.JobCtrl <- true
			}
			job.Err = path
			job.JobCtrl <- true
		}()

	case clientpb.GhostConfig_SHELLCODE:
		go func() {
			path, err := GhostSharedLibrary(&config)
			if err != nil {
				job.Err = err.Error()
				job.JobCtrl <- true
			}
			path, err = ShellcodeRDIToFile(path, "")
			if err != nil {
				job.Err = err.Error()
				job.JobCtrl <- true
			}

			job.Err = path
			job.JobCtrl <- true
		}()
	}
}
