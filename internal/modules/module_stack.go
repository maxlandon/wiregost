package modules

import (
	"fmt"
	"math/rand"
	"path"

	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/workspace"
)

var WorkspaceReqs = make(chan int)

type ModuleResponse struct {
	User       string
	Options    []Option
	ModuleName string
	Modules    []Module
}

type ModuleStack struct {
	WorkspaceId   int
	Id            int
	CurrentModule *Module
	Modules       []Module
}

type ModuleStackManager struct {
	Stacks map[int]ModuleStack
}

func NewModuleStackManager() *ModuleStackManager {
	man := &ModuleStackManager{Stacks: make(map[int]ModuleStack)}

	go man.handleWorkspaceRequests()
	go man.handleClientRequests()

	return man
}

func (m *ModuleStackManager) Create(workspaceId int) {
	m.Stacks[workspaceId] = ModuleStack{Id: rand.Int(), WorkspaceId: workspaceId}
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
			fmt.Println("use module")
			for _, s := range ms.Stacks {
				if s.WorkspaceId == request.CurrentWorkspaceId {
					fmt.Println("Found corresponding workspace")
					modName, _ := s.AddModule(request.Command[2])
					fmt.Println(modName)
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
		}
	}
}

func (stack *ModuleStack) AddModule(name string) (current string, err error) {
	var mPath = path.Join("/home/para/go/src/github.com/Ne0nd0g/merlin",
		"data", "modules", name+".json")
	module, err := Create(mPath)
	if err != nil {
		return "", err
	}
	stack.Modules = append(stack.Modules, module)

	stack.CurrentModule = &module
	if stack.CurrentModule != nil {
		return name, nil
	}
	return "", nil
}

func (stack *ModuleStack) PopModule(name string) {
	newStack := stack.Modules[:0]
	for _, m := range stack.Modules {
		if m.Name != name {
			newStack = append(newStack, m)
		}
	}
	stack.Modules = newStack
}

func (stack *ModuleStack) SaveStack() {

}

func (stack *ModuleStack) GetModuleList() map[string]string {
	list := make(map[string]string)
	for _, m := range stack.Modules {
		list[m.Name] = m.Platform
	}
	return list
}

func (stack *ModuleStack) LoadFromFile() {

}
