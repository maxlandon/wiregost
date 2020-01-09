package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/logging"
	"github.com/maxlandon/wiregost/internal/messages"
)

// ModuleRequests is a channel used for sending requests to Module Stack Manager
var ModuleRequests = make(chan messages.StackRequest, 20)

// ServerRequests is a channel used for sending requests to Server Manager
var ServerRequests = make(chan ServerRequest, 20)

// CompilerRequests is a channel used for sending requests to an associated Compiler
var CompilerRequests = make(chan CompilerRequest, 20)

// ServerRequest is the message used by Manager to request an action from a Server.
type ServerRequest struct {
	ClientID      int
	WorkspaceID   int
	Action        string
	WorkspacePath string
	Logger        *logging.WorkspaceLogger
}

// CompilerRequest is the message used by Manager to request an action from a Compiler.
type CompilerRequest struct {
	WorkspaceID   int
	Action        string
	WorkspacePath string
	Logger        *logging.WorkspaceLogger
}

// Workspace is an object used to structure work in Wiregost, with a dedicated agent server,
// compiler, module stack, and logger.
type Workspace struct {
	Name           string
	ID             int
	OwnerID        int
	Description    string
	Boundary       string
	LimitToNetwork bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Manager stores all workspaces and their associated loggers, and performs operations on them.
type Manager struct {
	Workspaces []Workspace
	Loggers    map[int]*logging.WorkspaceLogger
}

// NewManager instantiates a new Workspace Manager, which handles all requests from clients.
func NewManager() *Manager {
	ws := &Manager{
		Loggers: make(map[int]*logging.WorkspaceLogger),
	}
	// Load all workspaces
	ws.loadWorkspaces()

	go ws.handleRequests()
	go ws.handleLogRequests()
	return ws
}

func (wm *Manager) handleLogRequests() {
	for {
		req := <-messages.ForwardLogger
		wm.Loggers[req.CurrentWorkspaceID].GetLogs(req)
	}
}

func (wm *Manager) handleRequests() {
	for {
		req := <-messages.ForwardWorkspace
		switch req.Command[1] {
		// Create new workspace
		case "new":
			result := wm.create(req.Command[2], req.UserID, req.WorkspaceParams)
			// Respond to client
			res := messages.WorkspaceResponse{
				Result: result,
			}
			msg := messages.Message{
				ClientID: req.ClientID,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		// List workspaces
		case "list":
			res := messages.WorkspaceResponse{
				WorkspaceInfos: wm.getInfos(),
			}
			msg := messages.Message{
				ClientID: req.ClientID,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		case "switch":
			wm.switchWorkspace(req)
		case "delete":
			result := wm.deleteServer(req.Command[2], req.UserID)
			res := messages.WorkspaceResponse{
				Result: result,
			}
			msg := messages.Message{
				ClientID: req.ClientID,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		}
	}
}

func (wm *Manager) create(name string, ownerID int, params map[string]string) (result string) {
	// Create server object
	workspace := Workspace{
		Name:           name,
		ID:             rand.Int(),
		OwnerID:        ownerID,
		LimitToNetwork: false,
		CreatedAt:      time.Now(),
	}
	// Add optional parameters
	if v, ok := params["workspace.description"]; ok {
		workspace.Description = v
	}
	if v, ok := params["workspace.boundary"]; ok {
		workspace.Boundary = v
	}
	if v, ok := params["workspace.limit"]; ok {
		workspace.LimitToNetwork, _ = strconv.ParseBool(v)
	}

	// Add it to workspace list
	wm.Workspaces = append(wm.Workspaces, workspace)

	// Create workspace subdirectory
	workspaceDir, _ := fs.Expand("~/.wiregost/server/workspaces")
	if fs.Exists(workspaceDir) == false {
		os.MkdirAll(workspaceDir, 0755)
	}
	os.Mkdir(workspaceDir+"/"+workspace.Name, 0755)
	workspaceDir, _ = fs.Expand("~/.wiregost/server/workspaces" + "/" + workspace.Name)

	// Create subdirectory for agents
	os.Mkdir(workspaceDir+"/agents", 0755)

	// Save workspace properties in directory
	workspaceConf, _ := os.Create(workspaceDir + "/" + "workspace.conf")
	defer workspaceConf.Close()
	file, _ := fs.Expand(workspaceDir + "/" + "workspace.conf")
	var jsonData []byte
	jsonData, err := json.MarshalIndent(workspace, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to workspace configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(file, jsonData, 0755)
	}

	// Create its associated logger instance
	wm.Loggers[workspace.ID] = logging.NewWorkspaceLogger(workspace.Name, workspace.ID)
	// Create associated server
	ser := ServerRequest{
		WorkspaceID:   workspace.ID,
		WorkspacePath: workspaceDir,
		Action:        "create",
		Logger:        wm.Loggers[workspace.ID],
	}
	ServerRequests <- ser
	// Create new stack
	stackReq := messages.StackRequest{
		Action:      "create",
		WorkspaceID: workspace.ID,
		Workspace:   workspace.Name,
	}
	ModuleRequests <- stackReq
	// Create corresponding compiler
	compReq := CompilerRequest{
		WorkspaceID:   workspace.ID,
		WorkspacePath: workspaceDir,
		Action:        "create",
		Logger:        wm.Loggers[workspace.ID],
	}
	CompilerRequests <- compReq

	// Save directory config and return results
	workspace.saveServer()

	result = fmt.Sprintf("%s[*]%s Workspace '%s' created, with associated module stack, logger, and agent server.",
		tui.GREEN, tui.RESET, name)
	return result
}

func (wm *Manager) getWorkspaceList(ownerID int) []Workspace {
	var list []Workspace
	for _, ws := range wm.Workspaces {
		if ws.OwnerID == ownerID {
			list = append(list, ws)
		}
	}
	return list
}

func (wm *Manager) getInfos() [][]string {
	list := [][]string{}
	for _, w := range wm.Workspaces {
		var info []string
		info = append(info, w.Name)
		info = append(info, w.Description)
		info = append(info, w.Boundary)
		list = append(list, info)
	}
	return list
}

func (ws *Workspace) saveServer() {
	// Save workspace properties in directory
	workspaceDir, _ := fs.Expand("~/.wiregost/server/workspaces" + "/" + ws.Name)
	workspaceConf, _ := os.Create(workspaceDir + "/" + "workspace.conf")
	defer workspaceConf.Close()
	file, _ := fs.Expand(workspaceDir + "/" + "workspace.conf")
	var jsonData []byte
	jsonData, err := json.MarshalIndent(ws, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to workspace configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(file, jsonData, 0755)
	}
}

func (wm *Manager) deleteServer(name string, ownerID int) (result string) {
	var res string
	for _, ws := range wm.Workspaces {
		if ws.OwnerID == ownerID && name == ws.Name {
			path, _ := fs.Expand("~/.wiregost/server/workspaces/" + name)
			os.RemoveAll(path)
			res = fmt.Sprintf("%s[-]%s Deleted workspace %s and its directory content.",
				tui.GREEN, tui.RESET, ws.Name)
			// Delete tied HTTP/2 server
			servReq := ServerRequest{
				WorkspaceID:   ws.ID,
				WorkspacePath: path,
				Action:        "delete",
				Logger:        wm.Loggers[ws.ID],
			}
			ServerRequests <- servReq
			// Delete tied compiler
			compReq := CompilerRequest{
				WorkspaceID:   ws.ID,
				WorkspacePath: path,
				Action:        "delete",
				Logger:        wm.Loggers[ws.ID],
			}
			CompilerRequests <- compReq
			// Push notification to clients, fallback on default workspace
			defaultWorkspaceID := 0
			for _, w := range wm.Workspaces {
				if w.Name == "default" {
					defaultWorkspaceID = w.ID
				}
			}
			res := messages.Notification{
				Type:                "workspace",
				Action:              "delete",
				WorkspaceID:         ws.ID,
				FallbackWorkspaceID: defaultWorkspaceID,
			}
			messages.Notifications <- res
		}
	}
	// Update workspace list
	newList := wm.Workspaces[:0]
	for _, w := range wm.Workspaces {
		if name != w.Name {
			newList = append(newList, w)
		}
	}
	wm.Workspaces = newList

	return res
}

func (wm *Manager) switchWorkspace(request messages.ClientRequest) {
	for _, ws := range wm.Workspaces {
		if ws.Name == request.Command[2] {
			result := fmt.Sprintf("%s[*]%s => %s \n", tui.GREEN, tui.RESET, ws.Name)
			res := messages.WorkspaceResponse{
				WorkspaceID: ws.ID,
				Workspace:   ws.Name,
				Result:      result,
			}
			msg := messages.Message{
				ClientID: request.ClientID,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
			// Ask server manager to communicate status about associated server
			ser := ServerRequest{
				ClientID:    request.ClientID,
				WorkspaceID: ws.ID,
				Action:      "status",
			}
			ServerRequests <- ser
			// Ask StackManager to save stack for workspace
			stackReq := messages.StackRequest{
				Action:      "save",
				WorkspaceID: request.CurrentWorkspaceID,
				Workspace:   request.CurrentWorkspace,
			}
			ModuleRequests <- stackReq
		}
		// Save infos for current workspace
		if ws.ID == request.CurrentWorkspaceID {
			ws.saveServer()
		}
	}

}

func (wm *Manager) loadWorkspaces() {
	dir, _ := fs.Expand("~/.wiregost/server/workspaces/")
	dirs, _ := ioutil.ReadDir(dir)
	// If no workspaces are found, create a default one
	if dirs == nil {
		params := map[string]string{"Description": "Default Wiregost workspace."}
		wm.create("default", 1, params)
	}

	for _, d := range dirs {
		ws := Workspace{}
		confPath, _ := fs.Expand("~/.wiregost/server/workspaces/" + d.Name() + "/" + "workspace.conf")
		configBlob, _ := ioutil.ReadFile(confPath)
		json.Unmarshal(configBlob, &ws)
		wm.Workspaces = append(wm.Workspaces, ws)
		path, _ := fs.Expand("~/.wiregost/server/workspaces/" + d.Name())
		// Load associated loggers
		wm.Loggers[ws.ID] = logging.NewWorkspaceLogger(ws.Name, ws.ID)
		// Load associated compilers
		compReq := CompilerRequest{
			WorkspaceID:   ws.ID,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        wm.Loggers[ws.ID],
		}
		CompilerRequests <- compReq
		// Create new stack
		stackReq := messages.StackRequest{
			Action:      "load",
			WorkspaceID: ws.ID,
			Workspace:   ws.Name,
		}
		ModuleRequests <- stackReq
		// Load associated servers
		servReq := ServerRequest{
			WorkspaceID:   ws.ID,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        wm.Loggers[ws.ID],
		}
		ServerRequests <- servReq
	}
}
