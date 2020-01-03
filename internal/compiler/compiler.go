package compiler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/workspace"
)

// The compiler package contains all the code necessary to compiler objects.
// These objects are in charge of preparing and compiling various kinds of agents,
// with given parameters.

// There is one compiler instance per workspace, like module stacks.
// It is accessible through the shell and used in a special menu context, with
// a restrained but dedicated set of commands.

type CompilerResponse struct {
	User    string
	Options []Option
	Status  string
	Error   string
}

type Manager struct {
	Compilers map[int]Compiler
}

type Compiler struct {
	// Compiler properties
	Name string
	// Array of options
	Options []Option
}

// Option is a structure containing the keys for the object
type Option struct {
	Name        string `json:"name"`        // Name of the option
	Value       string `json:"value"`       // Value of the option
	Required    bool   `json:"required"`    // Is this a required option?
	Description string `json:"description"` // A description of the option
}

func NewManager() *Manager {
	man := &Manager{
		Compilers: make(map[int]Compiler),
	}

	go man.handleWorkspaceRequests()
	go man.handleClientRequests()

	return man
}

func (m *Manager) Create(path string, id int) {
	comp := Compiler{
		// Generic Options
		Options: []Option{
			Option{Name: "Agent Name", Value: "wiregostAgent", Required: true, Description: "Name of the agent (without file extension)."},
			Option{Name: "Architecture", Value: "", Required: true, Description: "Processor architecture of the target. (x64, 386, arm)"},
			Option{Name: "Operating System", Value: "", Required: true, Description: "Operating system of the target. (windows, linux, darwin)"},
			Option{Name: "Binary directory (server)", Value: "", Required: true, Description: "Directory in which to save the agent executable binary."},
			Option{Name: "Server Address/URL", Value: "", Required: true, Description: "URL of the server the agent will check-in. (Example: https://192.168.1.15:443, https://domain.com:443)"},
			Option{Name: "URL Path ", Value: "/", Required: true, Description: "The URL path to be appended to the server's address, for HTTP requests of the agent."},
			Option{Name: "Protocol", Value: "h2", Required: true, Description: "Protocol for the agent to connect with [ https (HTTP/1.1), h2 (HTTP/2), h3 (QUIC or HTTP/3) ]"},
			Option{Name: "Proxy", Value: "", Required: false, Description: "Hardcoded proxy to use for HTTP/1.1 traffic only, that will override host configuration."},
			Option{Name: "Host", Value: "", Required: false, Description: "HTTP Host header"},
			Option{Name: "Binary local", Value: "", Required: false, Description: "If set to a value, the agent binary will be copied from server-side to this directory, client-side."},
			Option{Name: "Pre-Shared Key", Value: "wiregost", Required: true, Description: "Pre-Shared Key to encrypt initial communications (default 'wiregost')"},

			Option{Name: "Windows GUI", Value: "true", Required: true, Description: "Will launch the agent into a window that is not visible to the end user. (Set to empty '' to disable)"},
		},
	}

	// Set workspace name as compiler name
	wsPath := strings.Split(path, "/")
	wsName := wsPath[len(wsPath)-1]
	comp.Name = wsName

	// Save compiler options in directory
	compilerConf, _ := os.Create(path + "/" + "compiler.conf")
	defer compilerConf.Close()
	file, _ := fs.Expand(path + "/" + "compiler.conf")
	var jsonData []byte
	jsonData, err := json.MarshalIndent(comp, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to compiler configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(file, jsonData, 0755)
		fmt.Println("Populated compiler.conf for compiler " + comp.Name)
	}

	// Add compiler to list
	m.Compilers[id] = comp
}

func (m *Manager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardCompiler
		switch request.Command[0] {
		case "list":
			fmt.Println("Detected request for compiler options")
			for k, v := range m.Compilers {
				if request.CurrentWorkspaceId == k {
					fmt.Println("Detected appropriate workspace ID")
					response := CompilerResponse{
						User:    "para",
						Options: v.Options,
					}
					msg := messages.Message{
						ClientId: request.ClientId,
						Type:     "compiler",
						Content:  response,
					}
					dispatch.Responses <- msg
					fmt.Println("Sent back compiler options")
				}
			}
		case "set":
			for k, v := range m.Compilers {
				if request.CurrentWorkspaceId == k {
					opt, err := v.SetOption(request.Command[1], request.Command[2])
					if err != nil {
						response := CompilerResponse{
							User:  "para",
							Error: err.Error(),
						}
						msg := messages.Message{
							ClientId: request.ClientId,
							Type:     "compiler",
							Content:  response,
						}
						dispatch.Responses <- msg
					} else {
						response := CompilerResponse{
							User:   "para",
							Status: opt,
						}
						msg := messages.Message{
							ClientId: request.ClientId,
							Type:     "compiler",
							Content:  response,
						}
						dispatch.Responses <- msg
					}
				}
			}
		}
	}
}

func (m *Manager) handleWorkspaceRequests() {
	for {
		request := <-workspace.CompilerRequests
		switch request.Action {
		case "create":
			m.Create(request.WorkspacePath, request.WorkspaceId)
		case "spawn":
			m.LoadCompilers(request.WorkspacePath, request.WorkspaceId)
		case "delete":
			delete(m.Compilers, request.WorkspaceId)
		}
	}
}

func (m *Manager) LoadCompilers(path string, id int) {
	// Load compiler
	comp := Compiler{}
	confPath, _ := fs.Expand(path + "/" + "compiler.conf")
	configBlob, _ := ioutil.ReadFile(confPath)
	json.Unmarshal(configBlob, &comp)
	// Add to list
	m.Compilers[id] = comp
}

func (c *Compiler) ListOptions() {

}

func (c *Compiler) Compile() {

}

func (c *Compiler) SaveParams() {
	// Save workspace properties in directory
	compDir, _ := fs.Expand("~/.wiregost/workspaces" + "/" + c.Name)
	compConf, _ := os.Create(compDir + "/" + "compiler.conf")
	defer compConf.Close()
	file, _ := fs.Expand(compDir + "/" + "compiler.conf")
	var jsonData []byte
	jsonData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to workspace configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(file, jsonData, 0755)
		fmt.Println("Saved compiler.conf for " + c.Name + " compiler")
	}
}

// SetOption is used to change the passed in compiler option's value. Used when a user is configuring a agent
func (c *Compiler) SetOption(option string, value string) (string, error) {
	// Verify this option exists
	for k, v := range c.Options {
		if option == v.Name {
			c.Options[k].Value = value
			return fmt.Sprintf("[-] %s set to %s", v.Name, c.Options[k].Value), nil
		}
	}
	return "", fmt.Errorf("%s[!]%s invalid module option: %s", tui.RED, tui.RESET, option)
}
