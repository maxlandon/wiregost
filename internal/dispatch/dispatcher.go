package dispatch

import (
	"fmt"

	"github.com/maxlandon/wiregost/internal/messages"
)

var (
	// Endpoint
	requests      = make(chan messages.ClientRequest, 20)
	Responses     = make(chan messages.Message, 20)
	Notifications = make(chan messages.Notification, 20)
	// Managers (buffered, so non-blocking)
	ForwardWorkspace     = make(chan messages.ClientRequest, 20)
	ForwardModuleStack   = make(chan messages.ClientRequest, 20)
	ForwardServerManager = make(chan messages.ClientRequest, 20)
	ForwardCompiler      = make(chan messages.ClientRequest, 20)
	ForwardLogger        = make(chan messages.ClientRequest, 20)
)

func DispatchRequest(req messages.ClientRequest) {
	// 1. Check commands: most of them can be run in either context
	// 2. For context-sensitive commands, check context
	fmt.Println(req.Command[0])
	switch req.Command[0] {
	// Server
	case "server":
		fmt.Println("launching handleServer")
		ForwardServerManager <- req
	// Log
	case "log":
		ForwardLogger <- req
		fmt.Println("Launching handleLog")
	// Stack
	case "stack":
		fmt.Println("Launching handleModule for stack")
		ForwardModuleStack <- req
	// Workspace
	case "workspace":
		fmt.Println("Launching handleWorkspace")
		ForwardWorkspace <- req
	// Module
	case "run", "show", "reload", "module":
		fmt.Println("launching handleModule")
		ForwardModuleStack <- req
	// Compiler:
	case "list", "compile", "compiler":
		fmt.Println("Dispatched request to handleCompiler")
		ForwardCompiler <- req
	// Agent
	case "agent", "interact", "cmd", "back", "download",
		"execute-shellcode", "kill", "main", "shell", "upload":
		fmt.Println("Launching handleAgent")
	// For both commands we need to check context
	case "use", "info", "set":
		switch req.Context {
		case "main":
			fmt.Println("Launching handleModule")
			ForwardModuleStack <- req
		case "module":
			fmt.Println("Launching handleModule")
			ForwardModuleStack <- req
		case "agent":
			fmt.Println("Launching handleAgent")
		case "compiler":
			ForwardCompiler <- req
		}
	}
}
