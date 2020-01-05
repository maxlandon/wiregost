package messages

// Client messages -----------------------------------------//
type ClientRequest struct {
	// Client-specific
	ClientId int
	// User-specific
	UserId       int
	UserName     string
	UserPassword string
	// Context-specific
	Context            string
	CurrentModule      string
	CurrentWorkspace   string
	CurrentWorkspaceId int
	// Command-specific
	Command []string
	// Server-specific
	ServerParams map[string]string
	// Workspace-specific
	WorkspaceParams map[string]string
}

type ClientConnRequest struct {
	UserName     string
	UserPassword string
}

// Server messages -----------------------------------------//

// Message acting as an envelope for other types of responses
type Message struct {
	ClientId         int
	Type             string
	NotificationType string
	Content          interface{}
}

// Message used to push updates to clients
type Notification struct {
	Type         string
	Action       string
	NotConcerned int
	// Workspace
	WorkspaceId         int
	FallbackWorkspaceId int
	Workspace           string
	// Module
	PoppedModule   string
	FallbackModule string
}

// Response to a workspace command
type WorkspaceResponse struct {
	User           string
	WorkspaceId    int // Return the current/chosen workspace here
	WorkspaceInfos [][]string
	Result         string
}

// Response to a log command
type LogResponse struct {
	User string
	Log  string // Used to notify log is set
	Logs []map[string]string
}

// Response to a server command
type ServerResponse struct {
	User       string
	Status     string
	Error      string
	ServerList []map[string]string
}

// Message used to push log events to clients
type LogEvent struct {
	ClientId    int
	WorkspaceId int
	Level       string
	Message     string
}

// Message used to send connection confirmation to a client.
type EndpointResponse struct {
	User      string
	Connected bool
	Status    string
}

// Message used by WorkspaceManager to request an action from a ModuleStack.
type StackRequest struct {
	WorkspaceId int
	Workspace   string
	Action      string
}

type AgentResponse struct {
	User string
	// Agents agents.Agents // Change this
	Info [][]string
}
