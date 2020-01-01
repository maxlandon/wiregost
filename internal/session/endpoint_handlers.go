package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
)

func (s *Session) EndpointConnect(cmd []string) {
	// In case where shell is already connected to a server, disconnect it
	if s.connection != nil {
		s.Disconnect()
	}

	// Set Current endpoint
	for _, v := range s.SavedEndpoints {
		if v.FQDN == cmd[2] {
			s.CurrentEndpoint = v
		}
	}
	s.Connect()
	// Send(cmd)
	// endpoint := <-endpointReqs
	// fmt.Println(endpoint)
}

func (s *Session) AddEndpoint() error {
	// Load params from Env
	params := make(map[string]string)
	for k, v := range s.Env {
		if strings.HasPrefix(k, "endpoint") {
			params[k] = v
		}
	}

	// Check for parameters
	if _, ok := params["endpoint.name"]; !ok {
		fmt.Println()
		fmt.Printf("%s[!]%s %s is not set.", tui.RED, tui.RESET, "'endpoint.name'")
		fmt.Println()
		return nil
	}
	if _, ok := params["endpoint.address"]; !ok {
		fmt.Println()
		fmt.Printf("%s[!]%s %s is not set.", tui.RED, tui.RESET, "'endpoint.address'")
		fmt.Println()
		return nil
	}
	if _, ok := params["endpoint.port"]; !ok {
		fmt.Println()
		fmt.Printf("%s[!]%s %s is not set.", tui.RED, tui.RESET, "'endpoint.port'")
		fmt.Println()
		return nil
	}
	if _, ok := params["endpoint.certificate"]; !ok {
		fmt.Println()
		fmt.Printf("%s[!]%s %s is not set.", tui.RED, tui.RESET, "'endpoint.certificate'")
		fmt.Println()
		return nil
	}
	if _, ok := params["endpoint.key"]; !ok {
		fmt.Println()
		fmt.Printf("%s[!]%s %s is not set.", tui.RED, tui.RESET, "'endpoint.key'")
		fmt.Println()
		return nil
	}
	// Check for default
	isDefault := false
	if val, ok := params["endpoint.default"]; !ok {
	} else {
		if val == "true" {
			isDefault = true
		}
	}

	// Load into template, add to ServersList and save file
	port, _ := strconv.Atoi(params["endpoint.port"])
	template := Endpoint{
		IPAddress:   params["endpoint.address"],
		Port:        port,
		Certificate: params["endpoint.certificate"],
		Key:         params["endpoint.key"],
		FQDN:        params["endpoint.name"],
		IsDefault:   isDefault,
	}
	s.SavedEndpoints = append(s.SavedEndpoints, template)
	fmt.Println()
	fmt.Printf("%s[*]%s Added endpoint %s at %s:%d \n", tui.GREEN, tui.RESET, template.FQDN, template.IPAddress, template.Port)
	s.WriteEndpointList()

	return nil
}

func (s *Session) DeleteEndpoint(cmd []string) error {
	newList := s.SavedEndpoints[:0]
	for _, v := range s.SavedEndpoints {
		if v.FQDN != cmd[2] {
			newList = append(newList, v)
		}
	}
	s.SavedEndpoints = newList
	s.WriteEndpointList()

	return nil
}

func (s *Session) SetDefaultEndpoint(cmd []string) error {
	for _, v := range s.SavedEndpoints {
		if v.IsDefault == true {
			v.IsDefault = false
		}
		if v.FQDN == cmd[2] {
			v.IsDefault = true
		}
	}

	s.WriteEndpointList()

	return nil
}

// List Servers
func (s *Session) ListEndpoints() error {
	columns := []string{
		tui.Yellow("FQDN (Common Name)"),
		tui.Yellow("Address"),
		tui.Yellow("Certificate"),
		tui.Yellow("Connected"),
		tui.Yellow("Default"),
	}

	rows := [][]string{}

	for _, l := range s.SavedEndpoints {
		row := []string{}
		// Name
		row = append(row, l.FQDN)
		// IP:Port
		address := l.IPAddress + ":" + strconv.Itoa(l.Port)
		row = append(row, address)
		// Certificate name (removing path)
		row = append(row, l.Certificate)
		// Connected
		if s.CurrentEndpoint == l {
			row = append(row, tui.Green("Connected"))
		}
		if s.CurrentEndpoint != l {
			row = append(row, " ")
		}
		// Default
		if l.IsDefault == true {
			row = append(row, "default")
		}
		if l.IsDefault == false {
			row = append(row, " ")
		}
		//Append to servers list
		rows = append(rows, row)
	}
	// Print table
	tui.Table(os.Stdout, columns, rows)
	return nil
}

func (s *Session) WriteEndpointList() error {
	endpointFile, _ := fs.Expand("~/.wiregost/client/server.conf")
	var servConf *os.File
	if !fs.Exists(endpointFile) {
		servConf, _ = os.Create(endpointFile)
	} else {
		servConf, _ = os.Open(endpointFile)
	}
	defer servConf.Close()

	// Marshal to JSON
	var jsonData []byte
	jsonData, err := json.MarshalIndent(s.SavedEndpoints, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to server configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(endpointFile, jsonData, 0755)
	}

	return nil
}
