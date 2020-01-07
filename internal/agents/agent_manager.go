package agents

import (
	"github.com/maxlandon/wiregost/internal/messages"
	uuid "github.com/satori/go.uuid"
)

type Manager struct {
	Agents map[uuid.UUID][]uuid.UUID
}

func NewManager() *Manager {
	manager := &Manager{
		Agents: make(map[uuid.UUID][]uuid.UUID, 0),
	}

	go manager.handleServerRequests()
	go manager.handleClientRequests()

	return manager
}

func (am *Manager) handleServerRequests() {
	for {
		request := <-messages.AgentRequests
		switch request.Action {
		case "add":
			am.Agents[request.ServerID] = append(am.Agents[request.ServerID], request.AgentID)
		}
	}
}

func (am *Manager) handleClientRequests() {
	for {
		request := <-messages.ForwardAgents
		switch request.Command[1] {
		case "show":
			am.showInfo(request)
		case "list":
			am.listServerAgents(request)
		}
	}
}

func (am *Manager) listServerAgents(request messages.ClientRequest) {
	serverUUID, _ := uuid.FromString(request.Command[2])

	agentList := make(map[string]int)
	if _, ok := am.Agents[serverUUID]; ok {
		agentList[request.Command[2]] = len(am.Agents[serverUUID])
	}
	res := messages.AgentResponse{
		AgentNb: agentList,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) showInfo(request messages.ClientRequest) {
	agentList := make([]map[string]string, 0)
	for _, a := range am.Agents[request.CurrentServerID] {
		infos := make(map[string]string)
		infos["id"] = Agents[a].ID.String()
		infos["platform"] = Agents[a].Platform
		infos["arch"] = Agents[a].Architecture
		infos["username"] = Agents[a].UserName
		infos["hostname"] = Agents[a].HostName
		infos["protocol"] = Agents[a].Proto
		infos["statusCheckIn"] = Agents[a].StatusCheckIn.Format("2006-01-02T15:04:05")

		agentList = append(agentList, infos)
	}
	res := messages.AgentResponse{
		Infos: agentList,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}
