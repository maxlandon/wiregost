package messages

import uuid "github.com/satori/go.uuid"

// Client messages -----------------------------------------//

// ClientRequest is a message used by clients to perform all requests
// to the managers of Wiregost. They will use the fields they need,
// depending on the context and command.
type ClientRequest struct {
	ClientID           int
	UserID             int
	UserName           string
	UserPassword       string
	Context            string
	CurrentModule      string
	CurrentWorkspace   string
	CurrentWorkspaceID int
	CurrentServerID    uuid.UUID
	ServerRunning      bool
	CurrentAgentID     uuid.UUID
	Command            []string
	ServerParams       map[string]string
	WorkspaceParams    map[string]string
}

// ClientConnRequest is used for requesting connection to the Wiregost Endpoint.
type ClientConnRequest struct {
	UserName     string
	UserPassword string
}

// Endpoint messages -----------------------------------------//

// Message is acting as an envelope for other types of responses.
type Message struct {
	ClientID         int
	Type             string
	NotificationType string
	Content          interface{}
}

// Notification is used to push updates to clients. This message is used
// for all sorts of notifications, and fields will be used only if needed.
type Notification struct {
	Type                string
	Action              string
	NotConcerned        int
	WorkspaceID         int
	ServerID            uuid.UUID
	ServerRunning       bool
	FallbackWorkspaceID int
	Workspace           string
	PoppedModule        string
	FallbackModule      string
}

// WorkspaceResponse is used to send back status/content to a workspace command.
type WorkspaceResponse struct {
	User           string
	WorkspaceID    int
	Workspace      string
	WorkspaceInfos [][]string
	Result         string
}

// LogResponse is used to send back status/content to a log command.
type LogResponse struct {
	User string
	Log  string // Used to notify log is set
	Logs []string
}

// ServerResponse is used to send back status/content to a server command.
type ServerResponse struct {
	User          string
	Status        string
	Error         string
	ServerList    []map[string]string
	ServerID      uuid.UUID
	ServerRunning bool
}

// LogEvent is used to push log events to clients
type LogEvent struct {
	ClientID    int
	WorkspaceID int
	Level       string
	Message     string
}

// EndpointResponse is used to send connection confirmation to a client.
type EndpointResponse struct {
	User      string
	Connected bool
	Status    string
}

// StackRequest is used by WorkspaceManager to request an action from a ModuleStack.
type StackRequest struct {
	WorkspaceID int
	Workspace   string
	Action      string
}

type AgentRequest struct {
	ServerID uuid.UUID
	Action   string
	AgentID  uuid.UUID
}

// AgentResponse is used for sending back status/content about an agent.
type AgentResponse struct {
	Infos     []map[string]string
	AgentNb   map[string]int
	AgentInfo [][]string
	Status    string
}
