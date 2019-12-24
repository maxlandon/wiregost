package workspace

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
)

var Requests = make(chan map[string]int)
var Responses = make(chan messages.Message)

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
	// Channels
	Requests  chan messages.ClientRequest
	Responses chan messages.Message
}

func NewWorkspaceManager() *WorkspaceManager {
	ws := &WorkspaceManager{
		Requests:  make(chan messages.ClientRequest),
		Responses: make(chan messages.Message),
	}

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
			id := w.Create(req.Command[2], req.UserId)
			m := make(map[string]int)
			m["create"] = id
			fmt.Println(m["create"])
			Requests <- m
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
			}
		}
	}
}

func (w *WorkspaceManager) Create(name string, ownerId int) (Id int) {
	workspace := Workspace{
		Name:           name,
		Id:             rand.Int(),
		OwnerID:        ownerId,
		LimitToNetwork: false,
		CreatedAt:      time.Now(),
	}

	w.Workspaces = append(w.Workspaces, workspace)
	fmt.Println(workspace.OwnerID)
	return workspace.Id
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

func (w *WorkspaceManager) Delete(name string, ownerId int) {

}

func (w *WorkspaceManager) Save() {

}

func (w *WorkspaceManager) Load() {

}
