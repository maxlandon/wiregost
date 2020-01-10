package agents

import (
	// Standard
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// 3rd party
	"github.com/evilsocket/islazy/tui"
	"github.com/mattn/go-shellwords"
	uuid "github.com/satori/go.uuid"

	// Wiregost
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/modules/shellcode"
)

// Manager is in charge of mapping agents to their servers, so that query is made
// with respect to a workspace context.
type Manager struct {
	// This map is just used for functions that require server filtering
	Agents map[uuid.UUID][]uuid.UUID
}

// NewManager instantiates an Agent Manager, which handles requests from Servers and Clients
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
		case "delete":
			am.deleteAgent(request)
		case "delete_all":
			for s, a := range am.Agents {
				if s == request.ServerID {
					for _, agent := range a {
						reqKill := messages.ClientRequest{Command: []string{"kill"}, CurrentAgentID: agent}
						fmt.Println(agent)
						am.killAgent(reqKill)
						time.Sleep(time.Second * 10)
					}
					delete(am.Agents, s)
				}
			}
			fmt.Println(am.Agents)
		}
	}
}

func (am *Manager) handleClientRequests() {
	for {
		request := <-messages.ForwardAgents
		switch request.Command[0] {
		case "agent":
			switch request.Command[1] {
			case "show":
				am.showAgents(request)
			case "list":
				am.listServerAgents(request)
			}
		case "info":
			am.showAgentInfo(request)
		case "kill":
			am.killAgent(request)
		case "ls":
			am.listDirectories(request)
		case "cd":
			am.changeAgentDirectory(request)
		case "pwd":
			am.printWorkingDirectory(request)
		case "cmd":
			am.agentCmd(request)
		case "download":
			am.agentDownload(request)
		case "upload":
			am.agentUpload(request)
		case "set":
			am.setAgentOption(request)
		case "execute-shellcode":
			am.executeShellCodeAgent(request)
		}
	}
}

func (am *Manager) deleteAgent(request messages.AgentRequest) {
	serverAgents := am.Agents[request.ServerID]
	newList := make([]uuid.UUID, 0)

	for _, a := range serverAgents {
		if a != request.AgentID {
			newList = append(newList, a)
		}
	}
	am.Agents[request.ServerID] = newList
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

func (am *Manager) showAgents(request messages.ClientRequest) {
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
		infos["status"] = GetAgentStatus(Agents[a].ID)

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

// ShowInfo lists all of the agent's structure value in a table
func (am *Manager) showAgentInfo(request messages.ClientRequest) {

	agentID := request.CurrentAgentID
	data := [][]string{
		{"Status", GetAgentStatus(agentID)},
		{"ID", Agents[agentID].ID.String()},
		{"Platform", Agents[agentID].Platform},
		{"Architecture", Agents[agentID].Architecture},
		{"UserName", Agents[agentID].UserName},
		{"User GUID", Agents[agentID].UserGUID},
		{"Hostname", Agents[agentID].HostName},
		{"Process ID", strconv.Itoa(Agents[agentID].Pid)},
		{"IP", fmt.Sprintf("%v", Agents[agentID].Ips)},
		{"Initial Check In", Agents[agentID].InitialCheckIn.Format(time.RFC3339)},
		{"Last Check In", Agents[agentID].StatusCheckIn.Format(time.RFC3339)},
		{"Agent Version", Agents[agentID].Version},
		{"Agent Build", Agents[agentID].Build},
		{"Agent Wait Time", Agents[agentID].WaitTime},
		{"Agent Wait Time Skew", strconv.FormatInt(Agents[agentID].Skew, 10)},
		{"Agent Message Padding Max", strconv.Itoa(Agents[agentID].PaddingMax)},
		{"Agent Max Retries", strconv.Itoa(Agents[agentID].MaxRetry)},
		{"Agent Failed Check In", strconv.Itoa(Agents[agentID].FailedCheckin)},
		{"Agent Kill Date", time.Unix(Agents[agentID].KillDate, 0).UTC().Format(time.RFC3339)},
		{"Agent Communication Protocol", Agents[agentID].Proto},
	}

	res := messages.AgentResponse{
		AgentInfo: data,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) listDirectories(request messages.ClientRequest) {
	cmd := request.Command
	var m string
	var status string
	var err error
	if len(cmd) > 1 {
		arg := strings.Join(cmd[0:], " ")
		argS, errS := shellwords.Parse(arg)
		if errS != nil {
			status = fmt.Sprintf("%s[!]%s There was an error parsing command line "+
				"arguments: %s\r\n%s", tui.YELLOW, tui.RESET, strings.Join(cmd, " "), errS.Error())
			goto response
		}
		m, err = AddJob(request.CurrentAgentID, "ls", argS)
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
			goto response
		}
		status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
			m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
	} else {
		m, err = AddJob(request.CurrentAgentID, cmd[0], cmd)
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
			goto response
		}
		status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
			m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
	}

	status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
		m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
	goto response

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) killAgent(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) > 0 {
		m, err := AddJob(request.CurrentAgentID, "kill", cmd[0:])
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
		} else {
			status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
				m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
		}
	}
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) changeAgentDirectory(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	var err error
	var m string
	if len(cmd) > 1 {
		arg := strings.Join(cmd[0:], " ")
		argS, errS := shellwords.Parse(arg)
		if errS != nil {
			status = fmt.Sprintf("%s[!]%s There was an error parsing command line argments: %s\r\n%s",
				tui.YELLOW, tui.RESET, strings.Join(cmd, " "), errS.Error())
			goto response
		}
		m, err = AddJob(request.CurrentAgentID, "cd", argS)
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
			goto response
		}
	} else {
		m, err = AddJob(request.CurrentAgentID, "cd", cmd)
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
			goto response
		}
	}
	status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
		m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
	goto response

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) printWorkingDirectory(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	var err error
	var m string
	m, err = AddJob(request.CurrentAgentID, "pwd", cmd)
	if err != nil {
		status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
		goto response
	}
	status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
		m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
	goto response

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) agentDownload(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) >= 2 {
		arg := strings.Join(cmd[1:], " ")
		argS, errS := shellwords.Parse(arg)
		if errS != nil {
			status = fmt.Sprintf("%s[!]%s There was an error parsing command line "+
				"arguments: %s\r\n%s", tui.YELLOW, tui.RESET, strings.Join(cmd, " "), errS.Error())
			goto response
		}
		if len(argS) >= 1 {
			m, err := AddJob(request.CurrentAgentID, "download", argS[0:1])
			if err != nil {
				status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				goto response
			} else {
				status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
					m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				goto response
			}
		}
	} else {
		status = fmt.Sprintf("%s[!]%s Invalid command. \n\tdownload <remote_file_path>", tui.YELLOW, tui.RESET)
		goto response
	}

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) agentCmd(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) > 1 {
		m, err := AddJob(request.CurrentAgentID, "cmd", cmd[1:])
		if err != nil {
			status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
		} else {
			status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
				m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
		}
	}
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) agentUpload(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) >= 3 {
		arg := strings.Join(cmd[1:], " ")
		argS, errS := shellwords.Parse(arg)
		if errS != nil {
			status = fmt.Sprintf("%s[!]%s There was an error parsing command line "+
				""+
				"arguments: %s\r\n%s", tui.YELLOW, tui.RESET, strings.Join(cmd, " "), errS.Error())
			goto response
		}
		if len(argS) >= 2 {
			_, errF := os.Stat(argS[0])
			if errF != nil {
				status = fmt.Sprintf("%s[!]%s There was an error accessing the source "+
					"upload file:\r\n%s", tui.YELLOW, tui.RESET, errF.Error())
				goto response
			}
			m, err := AddJob(request.CurrentAgentID, "upload", argS[0:2])
			if err != nil {
				status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				goto response
			} else {
				status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
					m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
			}
		}
	} else {
		status = fmt.Sprintf("%s[!]%s Invalid command. \n\tupload <local_file_path> <remote_file_path>", tui.YELLOW, tui.RESET)
		goto response
	}

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}

func (am *Manager) setAgentOption(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) > 1 {
		switch cmd[1] {
		case "killdate":
			if len(cmd) > 2 {
				_, errU := strconv.ParseInt(cmd[2], 10, 64)
				if errU != nil {
					status = fmt.Sprintf("%s[!]%s There was an error converting %s to an"+
						" int64 \n\tKill date takes in a UNIX epoch timestamp such as"+
						" 811123200 for September 15, 1995", tui.YELLOW, tui.RESET, cmd[2])
					break
				}
				m, err := AddJob(request.CurrentAgentID, "killdate", cmd[1:])
				if err != nil {
					status = fmt.Sprintf("There was an error adding a killdate "+
						"agent control message:\r\n%s", err.Error())
				} else {
					status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
						m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				}
			}
		case "maxretry":
			if len(cmd) > 2 {
				m, err := AddJob(request.CurrentAgentID, "maxretry", cmd[1:])
				if err != nil {
					status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				} else {
					status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
						m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				}
			}
		case "padding":
			if len(cmd) > 2 {
				m, err := AddJob(request.CurrentAgentID, "padding", cmd[1:])
				if err != nil {
					status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				} else {
					status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
						m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				}
			}
		case "sleep":
			if len(cmd) > 2 {
				m, err := AddJob(request.CurrentAgentID, "sleep", cmd[1:])
				if err != nil {
					status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				} else {
					status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
						m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				}
			}
		case "skew":
			if len(cmd) > 2 {
				m, err := AddJob(request.CurrentAgentID, "skew", cmd[1:])
				if err != nil {
					status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				} else {
					status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
						m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				}
			}
		}
		res := messages.AgentResponse{
			Status: status,
		}
		msg := messages.Message{
			ClientID: request.ClientID,
			Type:     "agent",
			Content:  res,
		}
		messages.Responses <- msg
	}
}

func (am *Manager) executeShellCodeAgent(request messages.ClientRequest) {
	cmd := request.Command
	var status string
	if len(cmd) > 2 {
		options := make(map[string]string)
		switch strings.ToLower(cmd[1]) {
		case "self":
			options["method"] = "self"
			options["pid"] = ""
			options["shellcode"] = strings.Join(cmd[2:], " ")
		case "remote":
			if len(cmd) > 3 {
				options["method"] = "remote"
				options["pid"] = cmd[2]
				options["shellcode"] = strings.Join(cmd[3:], " ")
			} else {
				status = fmt.Sprintf("%s[!]%s Not enough arguments. Try using the help command."+
					"\n\texecute-shellcode remote <pid> <shellcode>", tui.YELLOW, tui.RESET)
				goto response
			}
		case "rtlcreateuserthread":
			if len(cmd) > 3 {
				options["method"] = "rtlcreateuserthread"
				options["pid"] = cmd[2]
				options["shellcode"] = strings.Join(cmd[3:], " ")
			} else {
				status = fmt.Sprintf("%s[!]%s Not enough arguments. Try using the help command."+
					"\n\texecute-shellcode RtlCreateUserThread <pid> <shellcode>", tui.YELLOW, tui.RESET)
				goto response
			}
		case "userapc":
			if len(cmd) > 3 {
				options["method"] = "userapc"
				options["pid"] = cmd[2]
				options["shellcode"] = strings.Join(cmd[3:], " ")
			} else {
				status = fmt.Sprintf("%s[!]%s Not enough arguments. Try using the help command."+
					"\n\texecute-shellcode UserAPC <pid> <shellcode>", tui.YELLOW, tui.RESET)
				goto response
			}
		default:
			status = fmt.Sprintf("%s[!]%s invalid method provided", tui.YELLOW, tui.RESET)
			goto response
		}
		if len(options) > 0 {
			sh, errSh := shellcode.Parse(options)
			if errSh != nil {
				status = fmt.Sprintf("%s[!]%s there was an error parsing the shellcode:\r\n%s",
					tui.YELLOW, tui.RESET, errSh.Error())
				goto response
			}
			m, err := AddJob(request.CurrentAgentID, sh[0], sh[1:])
			if err != nil {
				status = fmt.Sprintf("%s[!]%s %s", tui.YELLOW, tui.RESET, err.Error())
				goto response
			} else {
				status = fmt.Sprintf("%s[*]%s Created job %s for agent %s at %s", tui.GREEN, tui.RESET,
					m, request.CurrentAgentID, time.Now().UTC().Format(time.RFC3339))
				goto response
			}
		}
	} else {
		status = fmt.Sprintf("%s[!]%s Not enough arguments were provided"+
			"\n\texecute-shellcode self <shellcode>"+
			"\n\texecute-shellcode remote <pid> <shellcode>"+
			"\n\texecute-shellcode RtlCreateUserThread <pid> <shellcode>", tui.YELLOW, tui.RESET)
		goto response
	}

response:
	res := messages.AgentResponse{
		Status: status,
	}
	msg := messages.Message{
		ClientID: request.ClientID,
		Type:     "agent",
		Content:  res,
	}
	messages.Responses <- msg
}
