package messages

// Dispatch
var (
	// Endpoint
	Requests      = make(chan ClientRequest, 20)
	Responses     = make(chan Message, 20)
	Notifications = make(chan Notification, 20)
	// Managers (buffered, so non-blocking)
	ForwardWorkspace     = make(chan ClientRequest, 20)
	ForwardModuleStack   = make(chan ClientRequest, 20)
	ForwardServerManager = make(chan ClientRequest, 20)
	ForwardCompiler      = make(chan ClientRequest, 20)
	ForwardLogger        = make(chan ClientRequest, 20)
)
