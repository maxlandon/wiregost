package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/agents"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/workspace"
)

// ModuleResponse is the message used to send back module/stack content and status to clients.
type ModuleResponse struct {
	ModuleName string
	Status     string
	Error      string
	Options    []Option
	ModuleList []string
	Modules    []*Module
}

// ModuleStack is a structure containing a list of previously loaded modules.
// It is used to keep their state and options alive during a session, and to have quick access to them.
type ModuleStack struct {
	WorkspaceID   int
	ID            int
	CurrentModule *Module
	Modules       []*Module
}

// Manager stores all modules stacks and perform operations either upon them or their modules.
// It also searches the list of available modules and manages them.
type Manager struct {
	Stacks map[int]*ModuleStack
}

// NewManager instantiates a new Module/Stack Manager, which handles all requests from clients or workspaces.
func NewManager() *Manager {
	man := &Manager{Stacks: make(map[int]*ModuleStack)}

	go man.handleWorkspaceRequests()
	go man.handleClientRequests()

	return man
}

func (msm *Manager) create(workspaceID int) {
	msm.Stacks[workspaceID] = &ModuleStack{ID: rand.Int(), WorkspaceID: workspaceID}
}

func (msm *Manager) handleWorkspaceRequests() {
	for {
		request := <-workspace.ModuleRequests
		switch request.Action {
		case "create":
			msm.create(request.WorkspaceID)
		case "load":
			msm.loadStack(request.Workspace, request.WorkspaceID)
		case "save":
			msm.saveStack(request.Workspace, request.WorkspaceID)
		}
	}
}

func (msm *Manager) handleClientRequests() {
	for {
		request := <-messages.ForwardModuleStack
		switch request.Command[0] {
		case "use":
			msm.useModule(request)
		// Cas "show" works for both info and options, as it sends the whole module
		case "show":
			msm.showModule(request)
		case "set":
			switch request.Command[1] {
			case "agent":
				msm.setAgent(request)
			default:
				msm.setOption(request)
			}
		case "run":
			msm.runModule(request)
		// List command for completers. This command "module" is not available in the shell
		case "module":
			getModuleList(request)
		// STACK COMMANDS
		case "stack":
			switch request.Command[1] {
			case "show":
				msm.getStackModuleList(request)
			case "pop":
				msm.popModule(request)
			case "list":
				msm.getStackModuleNames(request)
			}
		}
	}
}

func (msm *Manager) showModule(request messages.ClientRequest) {
	var module []*Module
	module = append(module, msm.Stacks[request.CurrentWorkspaceID].CurrentModule)
	response := ModuleResponse{
		Modules: module,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "module",
		Content:  response,
	}
	messages.Responses <- msg
}

func (msm *Manager) useModule(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceID]
	// Check if module already in stack
	name := request.Command[2]
	modPath := strings.Split(name, "/")
	modName := modPath[len(modPath)-1]
	if stack.Modules != nil {
		for _, mod := range stack.Modules {
			stackModNameSuf := mod.Path[len(mod.Path)-1]
			stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
			if strings.ToLower(stackModName) == strings.ToLower(modName) {
				stack.CurrentModule = mod
				// Dispatch response
				response := ModuleResponse{
					ModuleName: name,
				}
				msg := messages.Message{
					ClientID: request.ClientID,
					Type:     "module",
					Content:  response,
				}
				messages.Responses <- msg
				return
			}
		}
	}
	// If not, create it and add it to stack
	var mPath = path.Join("/home/para/go/src/github.com/Ne0nd0g/merlin",
		"data", "modules", name+".json")
	module, _ := Create(mPath)
	stack.Modules = append(stack.Modules, &module)

	stack.CurrentModule = &module
	if stack.CurrentModule != nil {
		// Dispatch response
		response := ModuleResponse{
			ModuleName: name,
		}
		msg := messages.Message{
			ClientID: request.ClientID,
			Type:     "module",
			Content:  response,
		}
		messages.Responses <- msg
		return
	}
}

func (msm *Manager) popModule(request messages.ClientRequest) {
	var poppedMod string
	// IDentify concerned stack and prepare it for changes.
	stack := msm.Stacks[request.CurrentWorkspaceID]
	newStack := stack.Modules[:0]
	// IDentify command
	switch len(request.Command) {
	// Pop current module
	case 2:
		for _, m := range stack.Modules {
			candidate := strings.ToLower(strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
			popped := strings.ToLower(strings.TrimSuffix(strings.Join(stack.CurrentModule.Path, "/"), ".json"))
			poppedMod = strings.TrimSuffix(strings.Join(stack.CurrentModule.Path, "/"), ".json")
			if candidate != popped {
				newStack = append(newStack, m)
			}
		}
		// If other modules in stack, use last one as current
		if len(newStack) != 0 {
			stack.CurrentModule = newStack[len(newStack)-1]
		} else {
			// If empty, just fill with an empty one
			stack.CurrentModule = &Module{}
		}
	case 3:
		// Pop all modules
		if request.Command[2] == "all" {
			stack = &ModuleStack{ID: rand.Int(), WorkspaceID: request.CurrentWorkspaceID}

		} else {
			// Pop selected one
			for _, m := range stack.Modules {
				candidate := strings.ToLower(strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
				if candidate != strings.ToLower(request.Command[2]) {
					newStack = append(newStack, m)
				}
			}
		}
	}
	// Set new stack
	stack.Modules = newStack

	// Send back new current module (empty response if no modules in stack left)
	currentMod := strings.TrimSuffix(strings.Join(stack.CurrentModule.Path, "/"), ".json")
	response := ModuleResponse{
		ModuleName: currentMod,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "module",
		Content:  response,
	}
	messages.Responses <- msg

	// Notify other clients to fallback.
	res := messages.Notification{
		Type:           "module",
		Action:         "pop",
		NotConcerned:   request.ClientID,
		WorkspaceID:    request.CurrentWorkspaceID,
		PoppedModule:   poppedMod,
		FallbackModule: currentMod,
	}
	messages.Notifications <- res
}

func (msm *Manager) setOption(request messages.ClientRequest) {
	// It is possible that the string formatting in this "set" case is overkill,
	// because we could juste compare names like in the "show" case juste above.
	// For now we keep it like that.
	for _, mod := range msm.Stacks[request.CurrentWorkspaceID].Modules {
		stackModNameSuf := strings.Join(mod.Path, "/")
		stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
		if strings.ToLower(stackModName) == strings.ToLower(request.CurrentModule) {
			opt, err := mod.SetOption(request.Command[1], request.Command[2])
			if err != nil {
				response := ModuleResponse{
					Error: err.Error(),
				}
				msg := messages.Message{
					ClientID: request.ClientID,
					Type:     "module",
					Content:  response,
				}
				messages.Responses <- msg
			} else {
				response := ModuleResponse{
					Status: opt,
				}
				msg := messages.Message{
					ClientID: request.ClientID,
					Type:     "module",
					Content:  response,
				}
				messages.Responses <- msg
			}
		}
	}

}

func (msm *Manager) setAgent(request messages.ClientRequest) {
	// It is possible that the string formatting in this "set" case is overkill,
	// because we could juste compare names like in the "show" case juste above.
	// For now we keep it like that.
	for _, mod := range msm.Stacks[request.CurrentWorkspaceID].Modules {
		stackModNameSuf := strings.Join(mod.Path, "/")
		stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
		if strings.ToLower(stackModName) == strings.ToLower(request.CurrentModule) {
			opt, err := mod.SetAgent(request.Command[2])
			if err != nil {
				response := ModuleResponse{
					Error: err.Error(),
				}
				msg := messages.Message{
					ClientID: request.ClientID,
					Type:     "module",
					Content:  response,
				}
				messages.Responses <- msg
			} else {
				response := ModuleResponse{
					Status: opt,
				}
				msg := messages.Message{
					ClientID: request.ClientID,
					Type:     "module",
					Content:  response,
				}
				messages.Responses <- msg
			}
		}
	}

}

func (msm *Manager) runModule(request messages.ClientRequest) {
	var status string
	var module *Module
	// IDentify concerned stack and prepare it for changes.
	stack := msm.Stacks[request.CurrentWorkspaceID]
	for _, mod := range stack.Modules {
		stackModNameSuf := strings.Join(mod.Path, "/")
		stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
		if strings.ToLower(stackModName) == strings.ToLower(request.CurrentModule) {
			module = mod
		}
	}
	// Run module, handle response.
	var m string
	r, err := module.Run()
	if err != nil {
		status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
		goto response
	}
	if len(r) <= 0 {
		status = fmt.Sprintf("%s[!]%s The %s module did not return a command to task an"+
			" agent with", tui.YELLOW, tui.RESET, module.Name)
		goto response
	}
	if strings.ToLower(module.Type) == "standard" {
		m, err = agents.AddJob(module.Agent, "cmd", r)
	} else {
		m, err = agents.AddJob(module.Agent, r[0], r[1:])
	}

	if err != nil {
		status = fmt.Sprintf("%s[!]%s There was an error adding the job to the specified agent:\n\t%s",
			tui.YELLOW, tui.RESET, err.Error())
		goto response
	} else {
		status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s",
			tui.GREEN, tui.RESET, m, module.Agent, time.Now().UTC().Format(time.RFC3339))
		goto response
	}

response:
	response := ModuleResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "module",
		Content:  response,
	}
	messages.Responses <- msg
}

func (msm *Manager) loadStack(name string, id int) {
	stack := ModuleStack{}
	confPath, _ := fs.Expand("~/.wiregost/server/workspaces/" + name + "/" + "stack.conf")
	if !fs.Exists(confPath) {
		msm.create(id)
		return
	}
	configBlob, _ := ioutil.ReadFile(confPath)
	json.Unmarshal(configBlob, &stack)
	msm.Stacks[stack.WorkspaceID] = &stack
}

func (msm *Manager) saveStack(name string, id int) {
	// Avoid saving for an empty workspace. WILL MOVE THAT ONCE DEFAULT IS USED AT SHELL SPAWN
	if name == "" {
		return
	}
	// Else save stack
	stack := msm.Stacks[id]
	stackDir, _ := fs.Expand("~/.wiregost/server/workspaces" + "/" + name)
	stackConf, _ := os.Create(stackDir + "/" + "stack.conf")
	defer stackConf.Close()
	// Save workspace properties in directory
	file, _ := fs.Expand(stackDir + "/" + "stack.conf")
	var jsonData []byte
	jsonData, err := json.MarshalIndent(stack, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to stack file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(file, jsonData, 0755)
	}
}

func (msm *Manager) getStackModuleList(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceID]
	var modules []*Module
	for _, m := range stack.Modules {
		modules = append(modules, m)
	}
	response := ModuleResponse{
		Modules: modules,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "module",
		Content:  response,
	}
	messages.Responses <- msg
}

// Function used for completion
func (msm *Manager) getStackModuleNames(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceID]
	modules := make([]string, 0)
	for _, m := range stack.Modules {
		modules = append(modules, strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
	}
	response := ModuleResponse{
		ModuleList: modules,
	}
	msg := messages.Message{
		Type:     "module",
		ClientID: request.ClientID,
		Content:  response,
	}
	messages.Responses <- msg
}

// Used for completion
func getModuleList(request messages.ClientRequest) {
	currentDir, _ := os.Getwd()
	ModuleDir := path.Join(filepath.ToSlash(currentDir), "data", "modules")
	list := make([]string, 0)

	err := filepath.Walk(ModuleDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", ModuleDir, err)
			return nil
		}
		if strings.HasSuffix(f.Name(), ".json") {
			d := strings.SplitAfter(filepath.ToSlash(path), ModuleDir)
			if len(d) > 0 {
				m := d[1]
				m = strings.TrimLeft(m, "/")
				m = strings.TrimSuffix(m, ".json")
				if !strings.Contains(m, "templates") {
					list = append(list, m)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", ModuleDir, err)
	}
	response := ModuleResponse{
		ModuleList: list,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "module",
		Content:  response,
	}
	messages.Responses <- msg
}
