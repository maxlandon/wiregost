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
	"github.com/maxlandon/wiregost/internal/dispatch"
	testlog "github.com/maxlandon/wiregost/internal/logging"
	"github.com/maxlandon/wiregost/internal/messages"
)

// Used to request the ModuleStack Manager to create a Stack for each new workspace
// Or to load a saved stack for existing ones.
var Requests = make(chan map[string]int, 20)

// Used for communicating a workspace ID and its server.conf path.
var ServerRequests = make(chan ServerRequest, 20)

var CompilerRequests = make(chan CompilerRequest, 20)

// Responses are sent back to clients.
var Responses = make(chan messages.Message)

type ServerRequest struct {
	WorkspaceId   int
	Action        string
	WorkspacePath string
	Logger        *testlog.WorkspaceLogger
}

type CompilerRequest struct {
	WorkspaceId   int
	Action        string
	WorkspacePath string
	Logger        *testlog.WorkspaceLogger
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
	Loggers    map[int]*testlog.WorkspaceLogger
	// Channels
	Requests  chan messages.ClientRequest
	Responses chan messages.Message
}

func NewWorkspaceManager() *WorkspaceManager {
	ws := &WorkspaceManager{
		Requests:  make(chan messages.ClientRequest),
		Responses: make(chan messages.Message),
		Loggers:   make(map[int]*testlog.WorkspaceLogger),
	}
	// Load all workspaces
	ws.LoadWorkspaces()

	go ws.HandleRequests()
	return ws
}

func (w *WorkspaceManager) HandleRequests() {
	for {
		// req := <-Requests
		req := <-dispatch.ForwardWorkspace
		switch req.Command[1] {
		case "list":
			workspaces := w.GetWorkspaceList(req.UserId)
			res := messages.WorkspaceResponse{}
			for _, ws := range workspaces {
				infos := ws.GetInfos()
				res.WorkspaceInfos = append(res.WorkspaceInfos, infos)
			}
			msg := messages.Message{
				ClientId: req.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			Responses <- msg
		case "new":
			id, path := w.Create(req)
			ser := ServerRequest{
				WorkspaceId:   id,
				WorkspacePath: path,
				Action:        "create",
			}
			ServerRequests <- ser
			// Create new stack
			stackReq := make(map[string]int)
			stackReq["create"] = id
			Requests <- stackReq
			// w.Responses <- msg
		case "switch":
			for _, ws := range w.Workspaces {
				if ws.Name == req.Command[2] {
					res := messages.WorkspaceResponse{
						WorkspaceId: ws.Id,
					}
					msg := messages.Message{
						ClientId: req.ClientId,
						Type:     "workspace",
						Content:  res,
					}
					Responses <- msg
				}
				// Save infos for current workspace
				if ws.Id == req.CurrentWorkspaceId {
					ws.SaveConf()
				}
			}
		case "delete":
			result := w.Delete(req.Command[2], req.UserId)
			res := messages.WorkspaceResponse{
				Result: result,
			}
			msg := messages.Message{
				ClientId: req.ClientId,
				Type:     "workspace",
				Content:  res,
			}
			Responses <- msg
		}
	}
}

func (w *WorkspaceManager) Create(req messages.ClientRequest) (Id int, path string) {
	// Create server object
	workspace := Workspace{
		Name:           req.Command[2],
		Id:             rand.Int(),
		OwnerID:        req.UserId,
		LimitToNetwork: false,
		CreatedAt:      time.Now(),
	}
	// Add optional parameters
	params := req.WorkspaceParams
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

	// Respond to client
	result := fmt.Sprintf("%s[*]%s Workspace '%s' created, with associated module stack, logger, and agent server.", tui.GREEN, tui.RESET, workspace.Name)
	res := messages.WorkspaceResponse{
		Result: result,
	}
	msg := messages.Message{
		ClientId: req.ClientId,
		Type:     "workspace",
		Content:  res,
	}
	Responses <- msg

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
	w.Loggers[workspace.Id] = testlog.NewWorkspaceLogger(workspace.Name, workspace.Id)

	// Create corresponding compiler
	compReq := CompilerRequest{
		WorkspaceId:   workspace.Id,
		WorkspacePath: workspaceDir,
		Action:        "create",
		Logger:        w.Loggers[workspace.Id],
	}
	CompilerRequests <- compReq

	// Save directory config
	workspace.SaveConf()

	// Return results
	return workspace.Id, workspaceDir
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

func (ws *Workspace) GetInfos() []string {
	var info []string
	info = append(info, ws.Name)
	info = append(info, ws.Description)
	info = append(info, ws.Boundary)
	return info
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
			res = fmt.Sprintf("%s[-]%s Deleted workspace %s and its directory content.", tui.GREEN, tui.RESET, ws.Name)
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

func (w *WorkspaceManager) LoadWorkspaces() {
	dir, _ := fs.Expand("~/.wiregost/workspaces/")
	dirs, _ := ioutil.ReadDir(dir)

	for _, d := range dirs {
		ws := Workspace{}
		confPath, _ := fs.Expand("~/.wiregost/workspaces/" + d.Name() + "/" + "workspace.conf")
		configBlob, _ := ioutil.ReadFile(confPath)
		json.Unmarshal(configBlob, &ws)
		w.Workspaces = append(w.Workspaces, ws)
		fmt.Println(tui.Dim("Loaded workspace " + d.Name()))
		path, _ := fs.Expand("~/.wiregost/workspaces/" + d.Name())
		// Load associated loggers
		w.Loggers[ws.Id] = testlog.NewWorkspaceLogger(ws.Name, ws.Id)
		// Load associated servers
		servReq := ServerRequest{
			WorkspaceId:   ws.Id,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        w.Loggers[ws.Id],
		}
		ServerRequests <- servReq
		// Load associated compilers
		compReq := CompilerRequest{
			WorkspaceId:   ws.Id,
			WorkspacePath: path,
			Action:        "spawn",
			Logger:        w.Loggers[ws.Id],
		}
		CompilerRequests <- compReq
		// Create new stack
		stackReq := make(map[string]int)
		stackReq["create"] = ws.Id
		Requests <- stackReq
	}
}
