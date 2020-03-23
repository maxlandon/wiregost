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

package nmap

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/maxlandon/wiregost/data-service/models"
)

// ScanRunner represents something that can run a scan.
type ScanRunner interface {
	Run() (result *models.Run, warnings []string, err error)
}

// Scanner represents an Nmap scanner.
type Scanner struct {
	cmd *exec.Cmd

	args       []string
	binaryPath string
	ctx        context.Context

	portFilter func(models.Port) bool
	hostFilter func(models.Host) bool

	stderr, stdout bufio.Scanner
}

// NewScanner creates a new Scanner, and can take options to apply to the scanner.
func NewScanner(options ...func(*Scanner)) (*Scanner, error) {
	scanner := &Scanner{}

	for _, option := range options {
		option(scanner)
	}

	if scanner.binaryPath == "" {
		var err error
		scanner.binaryPath, err = exec.LookPath("nmap")
		if err != nil {
			return nil, ErrNmapNotInstalled
		}
	}

	if scanner.ctx == nil {
		scanner.ctx = context.Background()
	}

	return scanner, nil
}

func NewCommandLineScanner(options []string) (*Scanner, error) {
	scanner := &Scanner{}
	if scanner.binaryPath == "" {
		var err error
		scanner.binaryPath, err = exec.LookPath("nmap")
		if err != nil {
			return nil, ErrNmapNotInstalled
		}
	}

	// Append options
	scanner.args = append(scanner.args, options...)

	if scanner.ctx == nil {
		scanner.ctx = context.Background()
	}

	return scanner, nil
}

// Run runs nmap synchronously and returns the result of the scan.
func (s *Scanner) Run() (result *models.Run, warnings []string, err error) {
	var stdout, stderr bytes.Buffer

	// Enable XML output
	s.args = append(s.args, "-oX")

	// Get XML output in stdout instead of writing it in a file
	s.args = append(s.args, "-")

	// Prepare nmap process
	cmd := exec.Command(s.binaryPath, s.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run nmap process
	err = cmd.Start()
	if err != nil {
		return nil, warnings, err
	}

	// Make a goroutine to notify the select when the scan is done.
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for nmap process or timeout
	select {
	case <-s.ctx.Done():

		// Context was done before the scan was finished.
		// The process is killed and a timeout error is returned.
		_ = cmd.Process.Kill()

		return nil, warnings, ErrScanTimeout
	case <-done:

		// Process nmap stderr output containing none-critical errors and warnings
		// Everyone needs to check whether one or some of these warnings is a hard issue in their use case
		if stderr.Len() > 0 {
			warnings = strings.Split(strings.Trim(stderr.String(), "\n"), "\n")
		}

		// Check for warnings that will inevitable lead to parsing errors, hence, have priority
		for _, warning := range warnings {
			switch {
			case strings.Contains(warning, "Malloc Failed!"):
				return nil, warnings, ErrMallocFailed
			// TODO: Add cases for other known errors we might want to guard.
			default:
			}
		}

		// Parse nmap xml output. Usually nmap always returns valid XML, even if there is a scan error.
		// Potentially available warnings are returned too, but probably not the reason for a broken XML.
		result, err := Parse(stdout.Bytes())
		if err != nil {
			warnings = append(warnings, err.Error()) // Append parsing error to warnings for those who are interested.
			return nil, warnings, ErrParseOutput
		}

		// Critical scan errors are reflected in the XML.
		if result != nil && len(result.Stats.Finished.ErrorMsg) > 0 {
			switch {
			case strings.Contains(result.Stats.Finished.ErrorMsg, "Error resolving name"):
				return result, warnings, ErrResolveName
			// TODO: Add cases for other known errors we might want to guard.
			default:
				return result, warnings, fmt.Errorf(result.Stats.Finished.ErrorMsg)
			}
		}

		// Call filters if they are set.
		if s.portFilter != nil {
			result = choosePorts(result, s.portFilter)
		}
		if s.hostFilter != nil {
			result = chooseHosts(result, s.hostFilter)
		}

		// Return result, optional warnings but no error
		return result, warnings, nil
	}
}

// RunAsync runs nmap asynchronously and returns error.
// TODO: RunAsync should return warnings as well.
func (s *Scanner) RunAsync() error {
	// Enable XML output.
	s.args = append(s.args, "-oX")

	// Get XML output in stdout instead of writing it in a file.
	s.args = append(s.args, "-")
	s.cmd = exec.Command(s.binaryPath, s.args...)

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to get error output from asynchronous nmap run: %v", err)
	}

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to get standard output from asynchronous nmap run: %v", err)
	}

	s.stdout = *bufio.NewScanner(stdout)
	s.stderr = *bufio.NewScanner(stderr)

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("unable to execute asynchronous nmap run: %v", err)
	}

	go func() {
		<-s.ctx.Done()
		_ = s.cmd.Process.Kill()
	}()

	return nil
}

// Wait waits for the cmd to finish and returns error.
func (s *Scanner) Wait() error {
	return s.cmd.Wait()
}

// GetStdout returns stdout variable for scanner.
func (s *Scanner) GetStdout() bufio.Scanner {
	return s.stdout
}

//  GetStdout returns stderr variable for scanner.
func (s *Scanner) GetStderr() bufio.Scanner {
	return s.stderr
}

func chooseHosts(result *models.Run, filter func(models.Host) bool) *models.Run {
	var filteredHosts []models.Host

	for _, host := range result.Hosts {
		if filter(host) {
			filteredHosts = append(filteredHosts, host)
		}
	}

	result.Hosts = filteredHosts

	return result
}

func choosePorts(result *models.Run, filter func(models.Port) bool) *models.Run {
	for idx := range result.Hosts {
		var filteredPorts []models.Port

		for _, port := range result.Hosts[idx].Ports {
			if filter(port) {
				filteredPorts = append(filteredPorts, port)
			}
		}

		result.Hosts[idx].Ports = filteredPorts
	}

	return result
}
