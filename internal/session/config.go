package session

// This file contains all the code for handling client-side configuration during use.
// This does not include configuration initialization (first time the client is used)

// It includes checking for presence of configuration directories and files, and
// loading them into the current shell session.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
)

type Config struct {
	// Directories
	UserDir      string
	LogDir       string
	ModulesDir   string
	PayloadsDir  string
	ResourceDir  string
	WorkspaceDir string
	ExportDir    string
	// Config file
	UserConfigFile string
	// History
	HistoryFile string
	// Global Variables
	GlobalVarFile string
	// Console Config
	ConsolePrompt string
}

func NewConfig() *Config {
	config := &Config{
		// Default settings are loaded, will be overwritten if
		// config file is found during LoadConfig()
		UserConfigFile: "~/.wiregost/client/ghost.conf",
		HistoryFile:    HistoryFile,
		UserDir:        "~/.wiregost/client/",
		LogDir:         "~/.wiregost/client/logs/",
		ModulesDir:     "~/.wiregost/client/modules/",
		PayloadsDir:    "~/.wiregost/client/payloads/",
		ResourceDir:    "~/.wiregost/client/resource/",
		WorkspaceDir:   "~/.wiregost/client/workspaces/",
		ExportDir:      "~/.wiregost/client/export/",
		ConsolePrompt:  "",
	}
	return config
}

// Parses and load all user-specific configuration files
func (conf *Config) LoadConfig() (err error) {
	// Check for personal directory. If no directory is found, exit the client.
	userDir, _ := fs.Expand(conf.UserDir)
	if fs.Exists(userDir) == false {
		fmt.Println(tui.Red(" ERROR: Personnal client directory does not exist."))
		fmt.Println(tui.Red("        Please run the ghost_setup.go script (in the scripts directory)," +
			" for initializing and configuring the client first"))
		os.Exit(1)
	} else {
		// Else load configuration
		fmt.Println(tui.Dim("Personal directory found."))
		path, _ := fs.Expand(conf.UserConfigFile)
		// If config file doesn't exist, exit the client
		if !fs.Exists(path) {
			fmt.Println(tui.Red("Configuration file not found: check for issues," +
				" or run the configuration script again"))
			os.Exit(1)
			// If config file is found, parse it.
		} else {
			configBlob, _ := ioutil.ReadFile(path)
			json.Unmarshal(configBlob, &conf)
		}
	}
	return err
}

func (conf *Config) ExportConfig() Config {

	config := Config{
		// Default settings are loaded, will be overwritten if
		// config file is found during LoadConfig()
		UserDir:      conf.UserDir,
		LogDir:       conf.LogDir,
		ModulesDir:   conf.ModulesDir,
		PayloadsDir:  conf.PayloadsDir,
		ResourceDir:  conf.ResourceDir,
		WorkspaceDir: conf.WorkspaceDir,
		ExportDir:    conf.ExportDir,

		UserConfigFile: conf.UserConfigFile,
		HistoryFile:    conf.HistoryFile,
		GlobalVarFile:  "",

		ConsolePrompt: conf.ConsolePrompt,
	}
	return config
}
