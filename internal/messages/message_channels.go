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
	ForwardAgents        = make(chan ClientRequest, 20)
	AgentRequests        = make(chan AgentRequest, 20)
	ForwardEnpoint       = make(chan WorkspaceResponse, 20)
	ForwardServer        = make(chan ServerResponse, 20)
	FromEndpoint         = make(chan ClientRequest, 20)
	EndpointToServer     = make(chan ClientRequest, 20)
)
