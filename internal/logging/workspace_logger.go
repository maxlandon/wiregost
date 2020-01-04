package logging

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/evilsocket/islazy/fs"
	"github.com/maxlandon/wiregost/internal/dispatch"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/sirupsen/logrus"
)

// Channels. All channels convey log entries, with all information needed
// var moduleLogs = make(chan log.Entry)
// var workspaceLogs = make(chan log.Entry)
var ForwardLogs = make(chan *logrus.Entry, 100)

// Each workspace has its own dedicated WorkspaceLogger instance, because it can then
// modulate its log level independently of others, and save logs in the appropriate directory.
// This logger is embedded to other components (Server, ModuleStack, Workspace) and allows them
// to log their information.
type WorkspaceLogger struct {
	*logrus.Logger
	WorkspaceId   int
	WorkspaceName string
	LogFile       string
}

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
		// DisableColors:   true,
		// FullTimestamp:   true,
	})
	// Setup level
	logger.SetLevel(logrus.DebugLevel)
	// Add hook to forward each log to clients
	hook := new(ForwardToDispatch)
	logger.AddHook(hook)
	// Setup log file and log path
	logger.LogFile, _ = fs.Expand("~/.wiregost/workspaces/" + name + "/" + name + ".log")
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

func (wl *WorkspaceLogger) GetLogs(request messages.ClientRequest) {
	// Setup all necessary paths
	// framework, _ := fs.Expand("~/.wiregost/server/wiregost.log")
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
				hlength += 1
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
						count += 1
					}
				}
				hlength -= 1
			}
			// Send back logs to client
			logs := messages.LogResponse{
				Logs: list,
			}
			msg := messages.Message{
				ClientId: request.ClientId,
				Type:     "log",
				Content:  logs,
			}
			dispatch.Responses <- msg
			fmt.Println("Sent back logs")
		}
	}
}

// Hooks
type Hook interface {
	Levels() []logrus.Level
	Fire(*logrus.Entry) error
}

type ForwardToDispatch struct {
}

// Forward log items to a general log dispatcher, at the Endpoint level.
func (h *ForwardToDispatch) Fire(entry *logrus.Entry) error {
	ForwardLogs <- entry
	return nil
}

func (h *ForwardToDispatch) Levels() []logrus.Level {
	return logrus.AllLevels
}
