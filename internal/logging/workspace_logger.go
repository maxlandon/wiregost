package logging

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/evilsocket/islazy/fs"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/sirupsen/logrus"
)

// ForwardLogs is used to forward log events of a workspace to all ClientLoggers.
var ForwardLogs = make(chan *logrus.Entry, 100)

// WorkspaceLogger is in charge of logging all events happening in a single workspace, and of
// forwarding these logs to all ClientLoggers. It also saves logs to disk.
// This logger is embedded to other components (Server, ModuleStack, Workspace) and allows them
// to log their information.
type WorkspaceLogger struct {
	*logrus.Logger
	WorkspaceID   int
	WorkspaceName string
	LogFile       string
}

// NewWorkspaceLogger instantiates a new Logger attached to a workspace.
func NewWorkspaceLogger(name string, id int) *WorkspaceLogger {
	logger := &WorkspaceLogger{
		logrus.New(),
		id,
		name,
		"",
	}
	// Setup formatting
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "01-02 15:04:05",
	})
	// Setup level
	logger.SetLevel(logrus.DebugLevel)
	// Add hook to forward each log to clients
	hook := new(ForwardToDispatch)
	logger.AddHook(hook)
	// Setup log file and log path
	logger.LogFile, _ = fs.Expand("~/.wiregost/server/workspaces/" + name + "/" + name + ".log")
	if !fs.Exists(logger.LogFile) {
		os.Create(logger.LogFile)
	}
	file, err := os.OpenFile(logger.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logOutput := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(logOutput)

	return logger
}

// GetLogs is used to send back a list of last x logs to a client, for a given workspace.
func (wl *WorkspaceLogger) GetLogs(request messages.ClientRequest) {
	// Setup list of logs
	list := make([]map[string]string, 0)

	switch len(request.Command) {
	// If three elements in commmand, asked for a specific set of logs.
	case 3:
		switch request.Command[2] {
		case "server":
			// Determine file length for subsequent selection
			file, _ := os.Open(wl.LogFile)
			hlength := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				hlength++
			}
			file.Close()
			// Read file and add each JSON line to list
			file, _ = os.Open(wl.LogFile)
			defer file.Close()
			count := 1
			scan := bufio.NewScanner(file)
			for scan.Scan() {
				if hlength <= 20 {
					line := make(map[string]string)
					json.Unmarshal(scan.Bytes(), &line)
					if line["component"] == request.Command[2] {
						list = append(list, line)
						count++
					}
				}
				hlength--
			}
			// Send back logs to client
			logs := messages.LogResponse{
				Logs: list,
			}
			msg := messages.Message{
				ClientID: request.ClientID,
				Type:     "log",
				Content:  logs,
			}
			messages.Responses <- msg
		}
	}
}

// Hook is an interface needed by logrus logger for triggering actions upon log receival.
type Hook interface {
	Levels() []logrus.Level
	Fire(*logrus.Entry) error
}

// ForwardToDispatch is necessary for hooks
type ForwardToDispatch struct {
}

// Fire forwards all log entries to ClientLoggers, at the Endpoint level.
func (h *ForwardToDispatch) Fire(entry *logrus.Entry) error {
	ForwardLogs <- entry
	return nil
}

// Levels is needed to satisfy the Hook interface
func (h *ForwardToDispatch) Levels() []logrus.Level {
	return logrus.AllLevels
}
