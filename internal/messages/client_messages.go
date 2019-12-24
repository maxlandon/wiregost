package messages

type ClientRequest struct {
	ClientId           int
	UserId             int
	UserName           string
	UserPassword       string
	Context            string
	CurrentWorkspace   string
	CurrentWorkspaceId int
	Command            []string
}

type ClientConnRequest struct {
	UserName     string
	UserPassword string
}

type Message struct {
	ClientId int
	Type     string
	Content  interface{}
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
	Log  string // Log message should normally be strings ?
}

type WorkspaceResponse struct {
	User string
	// WorkspaceList []Workspace // Return all workspaces with all their informations here.
	WorkspaceId    int // Return the current/chosen workspace here
	WorkspaceInfos [][]string
}

type StackResponse struct {
	User string
	// ModuleList []modules.Module // We will determine if we need to pass all modules or just their names/info
	// CurrentModule modules.Module // Maybe we will not need this line for changing shell state.
}

type ServerResponse struct {
	User      string
	Connected bool // Used upon connection, to notify shell it is correctly connected.
}
