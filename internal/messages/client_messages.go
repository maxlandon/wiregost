package messages

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

type Message struct {
	ClientId         int
	Type             string
	NotificationType string
	Content          interface{}
}

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

// type ModuleResponse struct {
//         User    string
//         Options []modules.Option
//         Modules []modules.Module
// }

type AgentResponse struct {
	User string
	// Agents agents.Agents // Change this
	Info [][]string
}

type LogResponse struct {
	User string
	Log  string // Used to notify log is set
	Logs []map[string]string
}

type LogEvent struct {
	ClientId    int
	WorkspaceId int
	Level       string
	Message     string
}

type WorkspaceResponse struct {
	User           string
	WorkspaceId    int // Return the current/chosen workspace here
	WorkspaceInfos [][]string
	Result         string
}

type StackResponse struct {
	User string
	// ModuleList []modules.Module // We will determine if we need to pass all modules or just their names/info
	// CurrentModule modules.Module // Maybe we will not need this line for changing shell state.
}

type EndpointResponse struct {
	User      string
	Connected bool // Used upon connection, to notify shell it is correctly connected.
	Status    string
}

type StackRequest struct {
	WorkspaceId int
	Workspace   string
	Action      string
}

type ServerResponse struct {
	User       string
	Status     string
	Error      string
	ServerList []map[string]string
}
