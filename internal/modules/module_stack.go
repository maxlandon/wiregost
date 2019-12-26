package modules

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/workspace"
)

type ModuleResponse struct {
	User       string
	Options    []Option
	ModuleName string
	Modules    []*Module
	ModuleList []string
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

func (ms *ModuleStackManager) handleClientRequests() {
	for {
		request := <-dispatch.ForwardModuleStack
		fmt.Println(request)
		switch request.Command[0] {
		case "use":
			for _, s := range ms.Stacks {
				if s.WorkspaceId == request.CurrentWorkspaceId {
					fmt.Println("Found corresponding workspace")
					modName, _ := s.AddModule(request.Command[2])
					response := ModuleResponse{
						User:       "para",
						ModuleName: modName,
					}
					msg := messages.Message{
						ClientId: request.ClientId,
						Type:     "module",
						Content:  response,
					}
					dispatch.Responses <- msg
					fmt.Println("dispatch received response")
				}
			}
		// Cas "show" works for both info and options, as it sends the whole module
		case "show":
			for _, s := range ms.Stacks {
				if s.WorkspaceId == request.CurrentWorkspaceId {
					var module []*Module
					module = append(module, s.CurrentModule)
					response := ModuleResponse{
						User:    "para",
						Modules: module,
					}
					msg := messages.Message{
						ClientId: request.ClientId,
						Type:     "module",
						Content:  response,
					}
					dispatch.Responses <- msg
				}
			}
			// It is possible that the string formatting in this "set" case is overkill,
			// because we could juste compare names like in the "show" case juste above.
			// For now we keep it like that.
		case "set":
			fmt.Println("detected set option command")
			for _, s := range ms.Stacks {
				if s.WorkspaceId == request.CurrentWorkspaceId {
					for _, mod := range s.Modules {
						stackModNameSuf := strings.Join(mod.Path, "/")
						stackModName := strings.TrimSuffix(stackModNameSuf, ".json")
						fmt.Printf("Stack mod name based on path")
						fmt.Println(strings.ToLower(stackModName))
						fmt.Println(strings.ToLower(request.CurrentModule))
						if strings.ToLower(stackModName) == strings.ToLower(request.CurrentModule) {
							opt, _ := mod.SetOption(request.Command[1], request.Command[2])
							fmt.Println(opt)
						}
					}
				}
			}
		// temporary list command for completers. This command "module" is not
		// available in the shell
		case "module":
			list := GetModuleList()
			response := ModuleResponse{
				User:       "para",
				ModuleList: list,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "module",
				Content:  response,
			}
			fmt.Println(list)
			dispatch.Responses <- msg
		// STACK COMMANDS
		case "stack":
			switch request.Command[1] {
			case "show":
				for _, s := range ms.Stacks {
					if s.WorkspaceId == request.CurrentWorkspaceId {
						modules := s.GetStackModuleList()
						response := ModuleResponse{
							User:    "para",
							Modules: modules,
						}
						msg := messages.Message{
							ClientId: request.ClientId,
							Type:     "module",
							Content:  response,
						}
						dispatch.Responses <- msg
					}
				}
			case "pop":
				switch len(request.Command) {
				case 2:
					for _, s := range ms.Stacks {
						if s.WorkspaceId == request.CurrentWorkspaceId {
							fmt.Println("Detected no module name, popping current one")
							s.PopModule(strings.TrimSuffix(strings.Join(s.CurrentModule.Path, "/"), ".json"))
						}
					}
				case 3:
					for _, s := range ms.Stacks {
						if s.WorkspaceId == request.CurrentWorkspaceId {
							fmt.Println("Detected module name, popping designated module")
							s.PopModule(request.Command[2])
						}
					}
				}
				// Case used for completing stack names
			case "list":
				for _, s := range ms.Stacks {
					if s.WorkspaceId == request.CurrentWorkspaceId {
						names := s.GetStackModuleNames()
						response := ModuleResponse{
							User:       "para",
							ModuleList: names,
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
	}
}

func (stack *ModuleStack) AddModule(name string) (current string, err error) {
	// Check if module already in stack
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
			stack.CurrentModule = mod // MODIFIED HERE, POPPED &
			return name, nil
		}
	}
	// If not, create it and add it to stack
	fmt.Println("Module not yet in stack, adding it")
	var mPath = path.Join("/home/para/go/src/github.com/Ne0nd0g/merlin",
		"data", "modules", name+".json")
	module, err := Create(mPath)
	if err != nil {
		return "", err
	}
	stack.Modules = append(stack.Modules, &module) // MODIFIED HERE , ADDED &
	fmt.Println("Stack modules after adding one")
	fmt.Println(stack.Modules)

	stack.CurrentModule = &module
	if stack.CurrentModule != nil {
		return name, nil
	}
	return "", nil
}

func (stack *ModuleStack) PopModule(name string) {
	newStack := stack.Modules[:0]
	for _, m := range stack.Modules {
		if strings.ToLower(strings.TrimSuffix(strings.Join(m.Path, "/"), ".json")) != strings.ToLower(name) {
			newStack = append(newStack, m)
		}
	}
	stack.Modules = newStack
}

func (stack *ModuleStack) SaveStack() {

}

func (stack *ModuleStack) GetStackModuleList() []*Module {
	var modules []*Module
	for _, m := range stack.Modules {
		modules = append(modules, m)
	}
	return modules
}

// Function used for completion
func (stack *ModuleStack) GetStackModuleNames() []string {
	modules := make([]string, 0)
	for _, m := range stack.Modules {
		modules = append(modules, strings.TrimSuffix(strings.Join(m.Path, "/"), ".json"))
	}
	return modules
}

func (stack *ModuleStack) LoadFromFile() {

}

func GetModuleList() []string {
	currentDir, _ := os.Getwd()
	ModuleDir := path.Join(filepath.ToSlash(currentDir), "data", "modules")
	o := make([]string, 0)

	err := filepath.Walk(ModuleDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", ModuleDir, err)
			return err
		}
		if strings.HasSuffix(f.Name(), ".json") {
			d := strings.SplitAfter(filepath.ToSlash(path), ModuleDir)
			if len(d) > 0 {
				m := d[1]
				m = strings.TrimLeft(m, "/")
				m = strings.TrimSuffix(m, ".json")
				if !strings.Contains(m, "templates") {
					o = append(o, m)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", ModuleDir, err)
	}
	return o
}
