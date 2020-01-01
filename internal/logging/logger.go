package logging

import (
	"io"
	"log"
	"os"

	"github.com/evilsocket/islazy/fs"
	"github.com/sirupsen/logrus"
)

// Channels. All channels convey log entries, with all information needed
// var moduleLogs = make(chan log.Entry)
// var workspaceLogs = make(chan log.Entry)
var ForwardLogs = make(chan *logrus.Entry, 100)

type LogEvent struct {
	Time        string
	Level       string
	Message     string
	Workspace   string
	WorkspaceId int
}

// Each workspace has its own dedicated WorkspaceLogger instance, because it can then
// modulate its log level independently of others, and save logs in the appropriate directory.
// This logger is embedded to other components (Server, ModuleStack, Workspace) and allows them
// to log their information.
type WorkspaceLogger struct {
	*logrus.Logger
	WorkspaceId   int
	WorkspaceName string
}

func NewWorkspaceLogger(name string, id int) *WorkspaceLogger {
	logger := &WorkspaceLogger{
		logrus.New(),
		id,
		name,
	}
	// Setup formatting
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
	})
	// Setup level
	logger.SetLevel(logrus.DebugLevel)
	// Add hook to forward each log to clients
	hook := new(ForwardToClients)
	logger.AddHook(hook)
	// Setup log file
	logfile, _ := fs.Expand("~/.wiregost/workspaces/" + name + "/" + name + ".log")
	// if fs.Exists("~/.wiregost/workspaces/"+name+"/"+name+".log") == false {
	if !fs.Exists(logfile) {
		os.Create(logfile)
		// os.Create("~/.wiregost/workspaces/" + name + "/" + name + ".log")
	}
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// file, err := os.OpenFile("~/.wiregost/workspaces/"+name+"/"+name+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logOutput := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(logOutput)

	return logger
}

// Hooks
type Hook interface {
	Levels() []logrus.Level
	Fire(*logrus.Entry) error
}
type ForwardToClients struct {
}

// Forward log items to a general log dispatcher, at the Endpoint level.
func (h *ForwardToClients) Fire(entry *logrus.Entry) error {
	ForwardLogs <- entry
	return nil
}

func (h *ForwardToClients) Levels() []logrus.Level {
	return logrus.AllLevels
}
