package dispatch

import (
	"github.com/maxlandon/wiregost/internal/messages"
)

// var moduleTest modules.Module

// Channels
var ForwardWorkspace = make(chan messages.ClientRequest)
var ForwardModuleStack = make(chan messages.ClientRequest)
var ForwardServerManager = make(chan messages.ClientRequest)
var ForwardCompiler = make(chan messages.ClientRequest)

func handleModule(req messages.ClientRequest) {
	ForwardModuleStack <- req
}

func handleWorkspace(req messages.ClientRequest) {
	ForwardWorkspace <- req
}

func handleServer(req messages.ClientRequest) {
	ForwardServerManager <- req
}

func handleCompiler(req messages.ClientRequest) {
	ForwardCompiler <- req
}
