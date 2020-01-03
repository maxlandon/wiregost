package modules

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/workspace"
)

type ModuleResponse struct {
	ModuleName string
	Status     string
	Error      string
	Options    []Option
	ModuleList []string
	Modules    []*Module
}

type ModuleStack struct {
	WorkspaceId   int
	Id            int
	CurrentModule *Module
	Modules       []*Module
}

type ModuleStackManager struct {
	Stacks map[int]*ModuleStack
}

func NewModuleStackManager() *ModuleStackManager {
	man := &ModuleStackManager{Stacks: make(map[int]*ModuleStack)}

	go man.handleWorkspaceRequests()
	go man.handleClientRequests()

	return man
}

func (m *ModuleStackManager) Create(workspaceId int) {
	m.Stacks[workspaceId] = &ModuleStack{Id: rand.Int(), WorkspaceId: workspaceId}
}

func (ms *ModuleStackManager) handleWorkspaceRequests() {
	for {
		request := <-workspace.Requests
		fmt.Println(request)
		f := func(key string) bool { _, ok := request[key]; return ok }
		switch {
		case f("create"):
			fmt.Println("Identified create request")
			fmt.Println("Creating stack linked to workspace")
			ms.Create(request["create"])
		case f("load"):
			fmt.Println("Received load request")
		}
	}
}

func (msm *ModuleStackManager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardModuleStack
		switch request.Command[0] {
		case "use":
			msm.UseModule(request)
		// Cas "show" works for both info and options, as it sends the whole module
		case "show":
			msm.ShowModule(request)
		case "set":
			msm.SetOption(request)
		// List command for completers. This command "module" is not available in the shell
		case "module":
			GetModuleList(request)
		// STACK COMMANDS
		case "stack":
			switch request.Command[1] {
			case "show":
				msm.GetStackModuleList(request)
			case "pop":
				msm.PopModule(request)
			case "list":
				msm.GetStackModuleNames(request)
			}
		}
	}
}

func (msm *ModuleStackManager) ShowModule(request messages.ClientRequest) {
	var module []*Module
	module = append(module, msm.Stacks[request.CurrentWorkspaceId].CurrentModule)
	response := ModuleResponse{
		Modules: module,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "module",
		Content:  response,
	}
	dispatch.Responses <- msg
}

func (msm *ModuleStackManager) UseModule(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceId]
	// Check if module already in stack
	name := request.Command[2]
	modPath := strings.Split(name, "/")
	modName := modPath[len(modPath)-1]
	fmt.Printf("Module name after split: ")
	fmt.Println(modName)
	fmt.Printf("Module names in stack: ")
	for _, mod := range stack.Modules {
		fmt.Println(mod.Name)
	}
	for _, mod := range stack.Modules {
		stackModNameSuf := mod.Path[len(mod.Path)-1]
		stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
		fmt.Printf("Stack mod name based on path")
		fmt.Println(stackModName)
		if strings.ToLower(stackModName) == strings.ToLower(modName) {
			fmt.Println("Module already in stack, updating current module")
			stack.CurrentModule = mod
			// Dispatch response
			response := ModuleResponse{
				ModuleName: name,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "module",
				Content:  response,
			}
			dispatch.Responses <- msg
			fmt.Println("dispatch received response")
			return
		}
	}
	// If not, create it and add it to stack
	fmt.Println("Module not yet in stack, adding it")
	var mPath = path.Join("/home/para/go/src/github.com/Ne0nd0g/merlin",
		"data", "modules", name+".json")
	module, _ := Create(mPath)
	stack.Modules = append(stack.Modules, &module)
	fmt.Println("Stack modules after adding one")
	fmt.Println(stack.Modules)

	stack.CurrentModule = &module
	if stack.CurrentModule != nil {
		// Dispatch response
		response := ModuleResponse{
			ModuleName: name,
		}
		msg := messages.Message{
			ClientId: request.ClientId,
			Type:     "module",
			Content:  response,
		}
		dispatch.Responses <- msg
		fmt.Println("dispatch received response")
		return
	}
}

func (msm *ModuleStackManager) PopModule(request messages.ClientRequest) {
	var poppedMod string
	// Identify concerned stack and prepare it for changes.
	var stack *ModuleStack
	var newStack []*Module
	for _, s := range msm.Stacks {
		if s.WorkspaceId == request.CurrentWorkspaceId {
			stack = s
			newStack = s.Modules[:0]
		}
	}
	// Identify command
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
		fmt.Println("stack length: " + strconv.Itoa(len(newStack)))
		if len(newStack) != 0 {
			stack.CurrentModule = newStack[len(newStack)-1]
			fmt.Println(stack.CurrentModule.Name)
		} else {
			// If empty, just fill with an empty one
			stack.CurrentModule = &Module{}
		}
		fmt.Println("Detected no module name, popping current one")
	case 3:
		// Pop all modules
		if request.Command[2] == "all" {
			stack = &ModuleStack{Id: rand.Int(), WorkspaceId: request.CurrentWorkspaceId}

		} else {
			// Pop selected one
			for _, m := range stack.Modules {
				candidate := strings.ToLower(strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
				if candidate != strings.ToLower(request.Command[2]) {
					newStack = append(newStack, m)
				}
			}
			fmt.Println("Detected module name, popping designated module")
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
		ClientId: request.ClientId,
		Type:     "module",
		Content:  response,
	}
	dispatch.Responses <- msg

	// Notify other clients to fallback.
	res := messages.Notification{
		Type:           "module",
		Action:         "pop",
		NotConcerned:   request.ClientId,
		WorkspaceId:    request.CurrentWorkspaceId,
		PoppedModule:   poppedMod,
		FallbackModule: currentMod,
	}
	dispatch.Notifications <- res
}

func (msm *ModuleStackManager) SetOption(request messages.ClientRequest) {
	fmt.Println("detected set option command")
	// It is possible that the string formatting in this "set" case is overkill,
	// because we could juste compare names like in the "show" case juste above.
	// For now we keep it like that.
	for _, mod := range msm.Stacks[request.CurrentWorkspaceId].Modules {
		stackModNameSuf := strings.Join(mod.Path, "/")
		stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
		fmt.Printf("Stack mod name based on path")
		fmt.Println(strings.ToLower(stackModName))
		fmt.Println(strings.ToLower(request.CurrentModule))
		if strings.ToLower(stackModName) == strings.ToLower(request.CurrentModule) {
			opt, err := mod.SetOption(request.Command[1], request.Command[2])
			if err != nil {
				response := ModuleResponse{
					Error: err.Error(),
				}
				msg := messages.Message{
					ClientId: request.ClientId,
					Type:     "module",
					Content:  response,
				}
				dispatch.Responses <- msg
			} else {
				response := ModuleResponse{
					Status: opt,
				}
				msg := messages.Message{
					ClientId: request.ClientId,
					Type:     "module",
					Content:  response,
				}
				dispatch.Responses <- msg
			}
		}
	}

}

func (stack *ModuleStack) SaveStack() {

}

func (msm *ModuleStackManager) GetStackModuleList(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceId]
	var modules []*Module
	for _, m := range stack.Modules {
		modules = append(modules, m)
	}
	response := ModuleResponse{
		Modules: modules,
	}
	msg := messages.Message{
		ClientId: request.ClientId,
		Type:     "module",
		Content:  response,
	}
	dispatch.Responses <- msg
}

// Function used for completion
func (msm *ModuleStackManager) GetStackModuleNames(request messages.ClientRequest) {
	stack := msm.Stacks[request.CurrentWorkspaceId]
	modules := make([]string, 0)
	for _, m := range stack.Modules {
		modules = append(modules, strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
	}
	response := ModuleResponse{
		ModuleList: modules,
	}
	msg := messages.Message{
		Type:     "module",
		ClientId: request.ClientId,
		Content:  response,
	}
	dispatch.Responses <- msg
}

func (stack *ModuleStack) LoadFromFile() {

}

// Used for completion
func GetModuleList(request messages.ClientRequest) {
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
		ClientId: request.ClientId,
		Type:     "module",
		Content:  response,
	}
	fmt.Println(list)
	dispatch.Responses <- msg
}
