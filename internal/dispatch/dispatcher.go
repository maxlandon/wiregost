package dispatch

import (
	"fmt"

	"github.com/maxlandon/wiregost/internal/messages"
)

var (
	// Channels
	requests  = make(chan messages.ClientRequest)
	Responses = make(chan messages.Message)
)

func DispatchRequest(req messages.ClientRequest) {
	// 1. Check commands: most of them can be run in either context
	// 2. For context-sensitive commands, check context
	fmt.Println(req.Command[0])
	switch req.Command[0] {
	// Server
	case "server":
		fmt.Println("launching handleServer")
		handleServer(req)
	// Log
	case "log":
		fmt.Println("Launching handleLog")
	// Stack
	case "stack":
		fmt.Println("Launching handleModule for stack")
		handleModule(req)
	// Workspace
	case "workspace":
		fmt.Println("Launching handleWorkspace")
		handleWorkspace(req)
	// Module
	case "run", "show", "reload", "module":
		fmt.Println("launching handleModule")
		handleModule(req)
	// Compiler:
	case "list", "compile", "compiler":
		fmt.Println("Dispatched request to handleCompiler")
		handleCompiler(req)
	// Agent
	case "agent", "interact", "cmd", "back", "download",
		"execute-shellcode", "kill", "main", "shell", "upload":
		fmt.Println("Launching handleAgent")
	// For both commands we need to check context
	case "use", "info", "set":
		switch req.Context {
		case "main":
			fmt.Println("Launching handleModule")
			handleModule(req)
		case "module":
			fmt.Println("Launching handleModule")
			handleModule(req)
		case "agent":
			fmt.Println("Launching handleAgent")
		case "compiler":
			handleCompiler(req)
		}
	}
}
