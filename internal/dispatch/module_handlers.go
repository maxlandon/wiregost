package dispatch

import (
	"github.com/maxlandon/wiregost/internal/messages"
)

// var moduleTest modules.Module

// Channels
var ForwardWorkspace = make(chan messages.ClientRequest)
var ForwardModuleStack = make(chan messages.ClientRequest)
var ForwardServerManager = make(chan messages.ClientRequest)

func handleModule(req messages.ClientRequest) {
	ForwardModuleStack <- req
	// if req.Command[0] == "show" {
	//         var mPath = path.Join("/home/para/go/src/github.com/Ne0nd0g/merlin", "data",
	//                 "modules", "windows/x64/powershell/credentials/LaZagneForensic"+".json")
	//         moduleTest, err := modules.Create(mPath)
	//         if err != nil {
	//                 fmt.Println(err.Error())
	//         }
	//         var options []modules.Option
	//         // Load options into list
	//         for _, o := range moduleTest.Options {
	//                 options = append(options, o)
	//         }
	//         response := modules.ModuleResponse{
	//                 User:    "para",
	//                 Options: options,
	//         }
	//         msg := messages.Message{
	//                 ClientId: req.ClientId,
	//                 Type:     "module",
	//                 Content:  response,
	//         }
	//         Responses <- msg
	// }
}

func handleWorkspace(req messages.ClientRequest) {
	ForwardWorkspace <- req
}

func handleServer(req messages.ClientRequest) {
	ForwardServerManager <- req
}
