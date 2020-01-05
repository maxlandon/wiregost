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

// Channels used for requesting other managers, to perform tasks such as loading, refreshing, etc.
var ModuleRequests = make(chan messages.StackRequest, 20) // Request ModuleStack to take action.
var ServerRequests = make(chan ServerRequest, 20)         // Request Server to take action.
var CompilerRequests = make(chan CompilerRequest, 20)     // Request Compiler to take action.

// Message used by WorkspaceManager to request an action from a Server.
type ServerRequest struct {
	ClientId      int
	WorkspaceId   int
	Action        string
	WorkspacePath string
	Logger        *logging.WorkspaceLogger
}

// Message used by WorkspaceManager to request an action from a Compiler.
type CompilerRequest struct {
	WorkspaceId   int
	Action        string
	WorkspacePath string
	Logger        *logging.WorkspaceLogger
}

type Workspace struct {
	Name           string
	Id             int
	OwnerID        int
	Description    string
	Boundary       string
	LimitToNetwork bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WorkspaceManager struct {
	Workspaces []Workspace
	Loggers    map[int]*logging.WorkspaceLogger
	// Channels
	Requests  chan messages.ClientRequest
	Responses chan messages.Message
}

func NewWorkspaceManager() *WorkspaceManager {
	ws := &WorkspaceManager{
		Requests:  make(chan messages.ClientRequest),
		Responses: make(chan messages.Message),
		Loggers:   make(map[int]*logging.WorkspaceLogger),
	}
	// Load all workspaces
	ws.LoadWorkspaces()

	go ws.HandleRequests()
	go ws.handleLogRequests()
	return ws
}

func (wm *WorkspaceManager) handleLogRequests() {
	for {
		req := <-messages.ForwardLogger
		wm.Loggers[req.CurrentWorkspaceId].GetLogs(req)
	}
}

func (wm *WorkspaceManager) HandleRequests() {
	for {
		req := <-messages.ForwardWorkspace
		switch req.Command[1] {
		// Create new workspace
		case "new":
			result := wm.Create(req.Command[2], req.UserId, req.WorkspaceParams)
			// Respond to client
			res := messages.WorkspaceResponse{
				Result: result,
			}
			msg := messages.Message{
				ClientId: req.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		// List workspaces
		case "list":
			res := messages.WorkspaceResponse{
				WorkspaceInfos: wm.GetInfos(),
			}
			msg := messages.Message{
				ClientId: req.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		case "switch":
			wm.SwitchWorkspace(req)
		case "delete":
			result := wm.Delete(req.Command[2], req.UserId)
			res := messages.WorkspaceResponse{
				Result: result,
			}
			msg := messages.Message{
				ClientId: req.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
		}
	}
}

func (w *WorkspaceManager) Create(name string, ownerId int, params map[string]string) (result string) {
	// Create server object
	workspace := Workspace{
		Name:           name,
		Id:             rand.Int(),
		OwnerID:        ownerId,
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
	w.Workspaces = append(w.Workspaces, workspace)

	// Create workspace subdirectory
	workspaceDir, _ := fs.Expand("~/.wiregost/workspaces")
	if fs.Exists(workspaceDir) == false {
		os.MkdirAll(workspaceDir, 0755)
		fmt.Println(" General workspace directory created")
	} else {
		fmt.Println(" General workspace directory found")
	}
	os.Mkdir(workspaceDir+"/"+workspace.Name, 0755)
	workspaceDir, _ = fs.Expand("~/.wiregost/workspaces" + "/" + workspace.Name)
	if fs.Exists(workspaceDir) == false {
		fmt.Println(" There was an error creating workspace directory")
	} else {
		fmt.Println(" Workspace directory created for " + workspace.Name)
	}

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
		fmt.Println("Populated workspace.conf for " + workspace.Name)
	}

	// Create its associated logger instance
	w.Loggers[workspace.Id] = logging.NewWorkspaceLogger(workspace.Name, workspace.Id)
	// Create associated server
	ser := ServerRequest{
		WorkspaceId:   workspace.Id,
		WorkspacePath: workspaceDir,
		Action:        "create",
		Logger:        w.Loggers[workspace.Id],
	}
	ServerRequests <- ser
	// Create new stack
	stackReq := messages.StackRequest{
		Action:      "create",
		WorkspaceId: workspace.Id,
		Workspace:   workspace.Name,
	}
	ModuleRequests <- stackReq
	// Create corresponding compiler
	compReq := CompilerRequest{
		WorkspaceId:   workspace.Id,
		WorkspacePath: workspaceDir,
		Action:        "create",
		Logger:        w.Loggers[workspace.Id],
	}
	CompilerRequests <- compReq

	// Save directory config and return results
	workspace.SaveConf()

	result = fmt.Sprintf("%s[*]%s Workspace '%s' created, with associated module stack, logger, and agent server.",
		tui.GREEN, tui.RESET, name)
	return result
}

func (w *WorkspaceManager) GetWorkspaceList(ownerId int) []Workspace {
	var list []Workspace
	for _, ws := range w.Workspaces {
		if ws.OwnerID == ownerId {
			list = append(list, ws)
		}
	}
	return list
}

func (wm *WorkspaceManager) GetInfos() [][]string {
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

func (ws *Workspace) SaveConf() {
	// Save workspace properties in directory
	workspaceDir, _ := fs.Expand("~/.wiregost/workspaces" + "/" + ws.Name)
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
		fmt.Println("Saved workspace.conf for " + ws.Name)
	}
}

func (w *WorkspaceManager) Delete(name string, ownerId int) (result string) {
	var res string
	for _, ws := range w.Workspaces {
		if ws.OwnerID == ownerId && name == ws.Name {
			path, _ := fs.Expand("~/.wiregost/workspaces/" + name)
			os.RemoveAll(path)
			res = fmt.Sprintf("%s[-]%s Deleted workspace %s and its directory content.",
				tui.GREEN, tui.RESET, ws.Name)
			// Delete tied HTTP/2 server
			servReq := ServerRequest{
				WorkspaceId:   ws.Id,
				WorkspacePath: path,
				Action:        "delete",
				Logger:        w.Loggers[ws.Id],
			}
			ServerRequests <- servReq
			// Delete tied compiler
			compReq := CompilerRequest{
				WorkspaceId:   ws.Id,
				WorkspacePath: path,
				Action:        "delete",
				Logger:        w.Loggers[ws.Id],
			}
			CompilerRequests <- compReq
			// Push notification to clients, fallback on default workspace
			defaultWorkspaceId := 0
			for _, w := range w.Workspaces {
				if w.Name == "default" {
					defaultWorkspaceId = w.Id
				}
			}
			res := messages.Notification{
				Type:                "workspace",
				Action:              "delete",
				WorkspaceId:         ws.Id,
				FallbackWorkspaceId: defaultWorkspaceId,
			}
			messages.Notifications <- res
		}
	}
	// Update workspace list
	newList := w.Workspaces[:0]
	for _, w := range w.Workspaces {
		if name != w.Name {
			newList = append(newList, w)
		}
	}
	w.Workspaces = newList

	return res
}

func (wm *WorkspaceManager) SwitchWorkspace(request messages.ClientRequest) {
	for _, ws := range wm.Workspaces {
		if ws.Name == request.Command[2] {
			result := fmt.Sprintf("%s[*]%s => %s \n", tui.GREEN, tui.RESET, ws.Name)
			res := messages.WorkspaceResponse{
				WorkspaceId: ws.Id,
				Result:      result,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			messages.Responses <- msg
			// Ask server manager to communicate status about associated server
			ser := ServerRequest{
				ClientId:    request.ClientId,
				WorkspaceId: ws.Id,
				Action:      "status",
			}
			ServerRequests <- ser
			// Ask StackManager to save stack for workspace
			stackReq := messages.StackRequest{
				Action:      "save",
				WorkspaceId: request.CurrentWorkspaceId,
				Workspace:   request.CurrentWorkspace,
			}
			fmt.Println("Request save with workspace: " + request.CurrentWorkspace)
			ModuleRequests <- stackReq
		}
		// Save infos for current workspace
		if ws.Id == request.CurrentWorkspaceId {
			ws.SaveConf()
		}
	}

}

func (wm *WorkspaceManager) LoadWorkspaces() {
	dir, _ := fs.Expand("~/.wiregost/workspaces/")
	dirs, _ := ioutil.ReadDir(dir)
	// If no workspaces are found, create a default one
	if dirs == nil {
		params := map[string]string{"Description": "Default Wiregost workspace."}
		wm.Create("default", 1, params)
	}

	for _, d := range dirs {
		ws := Workspace{}
		confPath, _ := fs.Expand("~/.wiregost/workspaces/" + d.Name() + "/" + "workspace.conf")
		configBlob, _ := ioutil.ReadFile(confPath)
		json.Unmarshal(configBlob, &ws)
		wm.Workspaces = append(wm.Workspaces, ws)
		fmt.Println(tui.Dim("Loaded workspace " + d.Name()))
		path, _ := fs.Expand("~/.wiregost/workspaces/" + d.Name())
		// Load associated loggers
		wm.Loggers[ws.Id] = logging.NewWorkspaceLogger(ws.Name, ws.Id)
		// Load associated compilers
		compReq := CompilerRequest{
			WorkspaceId:   ws.Id,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        wm.Loggers[ws.Id],
		}
		CompilerRequests <- compReq
		// Create new stack
		stackReq := messages.StackRequest{
			Action:      "load",
			WorkspaceId: ws.Id,
			Workspace:   ws.Name,
		}
		// Load associated servers
		servReq := ServerRequest{
			WorkspaceId:   ws.Id,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        wm.Loggers[ws.Id],
		}
		ServerRequests <- servReq
		fmt.Println("Sending stack request for workspace " + strconv.Itoa(ws.Id))
		ModuleRequests <- stackReq
	}
}
