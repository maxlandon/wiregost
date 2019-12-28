package modules

import (
	// Standard
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	// 3rd Party
	"github.com/evilsocket/islazy/tui"
	"github.com/fatih/color"
	uuid "github.com/satori/go.uuid"
	// Merlin
)

// Module is a structure containing the base information or template for modules
type Module struct {
	Agent        uuid.UUID   // The Agent that will later be associated with this module prior to execution
	Name         string      `json:"name"`                 // Name of the module
	Author       []string    `json:"author"`               // A list of module authors
	Credits      []string    `json:"credits"`              // A list of people to credit for underlying tool or techniques
	Path         []string    `json:"path"`                 // Path to the module (i.e. data/modules/powershell/powerview)
	Platform     string      `json:"platform"`             // Platform the module can run on (i.e. Windows, Linux, Darwin, or ALL)
	Arch         string      `json:"arch"`                 // The Architecture the module can run on (i.e. x86, x64, MIPS, ARM, or ALL)
	Lang         string      `json:"lang"`                 // What language does the module execute in (i.e. PowerShell, Python, or Perl)
	Priv         bool        `json:"privilege"`            // Does this module required a privileged level account like root or SYSTEM?
	Description  string      `json:"description"`          // A description of what the module does
	Notes        string      `json:"notes"`                // Additional information or notes about the module
	Commands     []string    `json:"commands"`             // A list of commands to be run on the agent
	SourceRemote string      `json:"remote"`               // Online or remote source code for a module (i.e. https://raw.githubusercontent.com/PowerShellMafia/PowerSploit/master/Exfiltration/Invoke-Mimikatz.ps1)
	SourceLocal  []string    `json:"local"`                // The local file path to the script or payload
	Options      []Option    `json:"options"`              // A list of configurable options/arguments for the module
	Powershell   interface{} `json:"powershell,omitempty"` // An option json object containing commands and configuration items specific to PowerShell
}

// Option is a structure containing the keys for the object
type Option struct {
	Name        string `json:"name"`        // Name of the option
	Value       string `json:"value"`       // Value of the option
	Required    bool   `json:"required"`    // Is this a required option?
	Flag        string `json:"flag"`        // The command line flag used for the option
	Description string `json:"description"` // A description of the option
}

// PowerShell structure is used to describe additional PowerShell features for modules that leverage PowerShell
type PowerShell struct {
	DisableAV   bool // Disable Windows Real Time "Set-MpPreference -DisableRealtimeMonitoring $true"
	Obfuscation bool // Unimplemented command to obfuscated powershell
	Base64      bool // Base64 encode the powershell command?
}

// Run function returns an array of commands to execute the module on an agent
func (m *Module) Run() ([]string, error) {
	if m.Agent == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000") {
		return nil, errors.New("agent not set for module")
	}

	// Check every 'required' option to make sure it isn't null
	for _, v := range m.Options {
		if v.Required {
			if v.Value == "" {
				return nil, errors.New(v.Name + " is required")
			}
		}
	}

	// Fill in or remove options values
	command := make([]string, len(m.Commands))
	copy(command, m.Commands)

	for _, o := range m.Options {
		for k := len(command) - 1; k >= 0; k-- {
			// Check if an option was set WITHOUT the Flag or Value qualifiers
			if strings.Contains(command[k], "{{"+o.Name+"}}") {
				if o.Value != "" {
					command[k] = strings.Replace(command[k], "{{"+o.Name+"}}", o.Flag+" "+o.Value, -1)
				} else {
					command = append(command[:k], command[k+1:]...)
				}
				// Check if an option was set WITH just the Flag qualifier
			} else if strings.Contains(command[k], "{{"+o.Name+".Flag}}") {
				if strings.ToLower(o.Value) == "true" {
					command[k] = strings.Replace(command[k], "{{"+o.Name+".Flag}}", o.Flag, -1)
				} else {
					command = append(command[:k], command[k+1:]...)
				}
				// Check if an option was set WITH just the Value qualifier
			} else if strings.Contains(command[k], "{{"+o.Name+".Value}}") {
				if o.Value != "" {
					command[k] = strings.Replace(command[k], "{{"+o.Name+".Value}}", o.Value, -1)
				} else {
					command = append(command[:k], command[k+1:]...)
				}
			}
		}
	}
	return command, nil
}

// SetOption is used to change the passed in module option's value. Used when a user is configuring a module
func (m *Module) SetOption(option string, value string) (string, error) {
	// Verify this option exists
	for k, v := range m.Options {
		if option == v.Name {
			m.Options[k].Value = value
			return fmt.Sprintf("[-] %s set to %s", v.Name, m.Options[k].Value), nil
		}
	}
	return "", fmt.Errorf("%s[!]%s invalid module option: %s", tui.RED, tui.RESET, option)
}

// SetAgent is used to set the agent associated with the module.
func (m *Module) SetAgent(agentUUID string) (string, error) {
	if strings.ToLower(agentUUID) == "all" {
		agentUUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	}
	i, err := uuid.FromString(agentUUID)
	if err != nil {
		return "", fmt.Errorf("invalid UUID")
	}
	m.Agent = i
	return fmt.Sprintf("agent set to %s", m.Agent.String()), nil
}

// Create is module function used to instantiate a module object using the provided file path to a module's json file
func Create(modulePath string) (Module, error) {
	var m Module

	// Read in the module's JSON configuration file
	f, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return m, err
	}

	// Unmarshal module's JSON message
	var moduleJSON map[string]*json.RawMessage
	errModule := json.Unmarshal(f, &moduleJSON)
	if errModule != nil {
		return m, errModule
	}

	// Determine all message types
	var keys []string
	for k := range moduleJSON {
		keys = append(keys, k)
	}

	// Validate that module's JSON contains at least the base message
	var containsBase bool
	for i := range keys {
		if keys[i] == "base" {
			containsBase = true
		}
	}

	// Marshal Base message type
	if !containsBase {
		return m, errors.New("the module's definition does not contain the 'BASE' message type")
	}
	errJSON := json.Unmarshal(*moduleJSON["base"], &m)
	if errJSON != nil {
		return m, errJSON
	}

	// Check for PowerShell configuration options
	for k := range keys {
		switch keys[k] {
		case "base":
		case "powershell":
			k := marshalMessage(*moduleJSON["powershell"])
			m.Powershell = (*json.RawMessage)(&k)
			var p PowerShell
			json.Unmarshal(k, &p)
		}
	}

	_, errValidate := validateModule(m)
	if errValidate != nil {
		return m, errValidate
	}
	return m, nil
}

// validate function is used to check a module's configuration for errors
func validateModule(m Module) (bool, error) {

	// Validate Platform
	switch strings.ToUpper(m.Platform) {
	case "WINDOWS":
	case "LINUX":
	case "DARWIN":
	default:
		return false, errors.New("invalid 'platform' value provided in module file")
	}

	// Validate Architecture
	switch strings.ToUpper(m.Arch) {
	case "X64":
	case "X32":
	default:
		return false, errors.New("invalid 'arch' value provided in module file")
	}
	return true, nil
}

// marshalMessage is a generic function used to marshal JSON messages
func marshalMessage(m interface{}) []byte {
	k, err := json.Marshal(m)
	if err != nil {
		color.Red("There was an error marshaling the JSON object")
		color.Red(err.Error())
	}
	return k
}
