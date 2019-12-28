package cli

import (
	"fmt"
	"strconv"
	"strings"

	// Third Party
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

func generalHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Command Categories\n")))
	fmt.Println(tui.Dim("(Type 'help CATEGORY' to show category help)"))
	fmt.Println()

	maxLen := 0
	for _, c := range commandCategories {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"
	for _, c := range commandCategories {
		fmt.Printf("  "+tui.Yellow(pad)+" : %s\n", c.Name, c.Description)
	}
	fmt.Println()
}

func coreHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Core Commands\n")))
	var params string
	maxLen := 0
	for _, c := range coreCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"
	for _, c := range coreCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println()
}

func serverHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Server Commands\n")))
	var params string
	maxLen := 0
	for _, c := range serverCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range serverCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))
	fmt.Println(tui.Dim("Before instantiating or starting a server, make sure all paramaters have the correct/wished settings."))
	fmt.Println()

	// Parameters
	maxLen = 0
	for _, c := range serverParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range serverParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}
	fmt.Println()
}

func endpointHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Endpoint Commands\n")))
	var params string
	maxLen := 0
	for _, c := range endpointCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range endpointCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))

	// Parameters
	maxLen = 0
	for _, c := range endpointParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range endpointParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}
	fmt.Println()
}

func logHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Log Commands\n")))
	var params string
	maxLen := 0
	for _, c := range logCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range logCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	// Parameters
	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))
	maxLen = 0
	for _, c := range logParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range logParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}
	fmt.Println()
}

func chatHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Chat Commands\n")))
	var params string
	maxLen := 0
	for _, c := range chatCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range chatCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println()
}

func workspaceHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Workspace Commands\n")))
	var params string
	maxLen := 0
	for _, c := range workspaceCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range workspaceCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	// Parameters
	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))
	maxLen = 0
	for _, c := range workspaceParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range workspaceParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}
	fmt.Println()
}

func stackHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Module Stack Commands\n")))
	var params string
	maxLen := 0
	for _, c := range stackCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range stackCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println()
}

func moduleHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Module Commands\n")))
	var params string
	maxLen := 0
	for _, c := range moduleCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range moduleCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println()
}

func agentHelp() {
	// Commands
	fmt.Println(tui.Bold(tui.Blue("\n  Agent Commands\n")))
	fmt.Println(tui.Dim("These commands are only available when interacting with an agent. ('agent interact <agent_uuid>')"))
	fmt.Println()
	var params string
	maxLen := 0
	for _, c := range agentCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range agentCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}
	fmt.Println()
}

func helpHandler(args []string) error {
	filter := ""
	if len(args) == 2 {
		filter = str.Trim(args[1])
	}

	switch filter {
	case "":
		generalHelp()
	case "core":
		coreHelp()
	case "server":
		serverHelp()
	case "log":
		logHelp()
	case "chat":
		chatHelp()
	case "workspace":
		workspaceHelp()
	case "stack":
		stackHelp()
	case "endpoint":
		endpointHelp()
	case "module":
		moduleHelp()
	case "agent":
		agentHelp()
	}

	return nil
}

type CommandDescription struct {
	Name        string
	Params      []string
	Description string
}

type ParamDescription struct {
	Name        string
	Description string
	Default     string
}

//-------------------------------------------------------------------------------------------------------------------------
// List of all handler sets and their description
var commandCategories = []CommandDescription{
	{Name: "core", Description: "WireGost core commands, (resource loading and making, input & navigation mode, etc...)"},
	{Name: "endpoint", Description: "Commands and parameters for managing WireGost clients' Endpoint (connection, add, generate tokens, etc)"},
	{Name: "server", Description: "Commands and parameters for managing Agent Servers (state, certificates, tokens, etc.)"},
	{Name: "log", Description: "Commands for managing the various sets of logs used by WireGost"},
	// {Name: "chat", Description: "Commands for using WireGost's messaging system"},
	{Name: "workspace", Description: "Manage WireGost workspaces"},
	{Name: "stack", Description: "Manage the module stack (all modules currently loaded in this session)"},
	// {Name: "global", Description: "Manage all global variables in WireGost"},
	// {Name: "db", Description: "Manage WireGost data services"},
	// {Name: "listeners", Description: "Manage listeners instantiated in this session"},
	// {Name: "exploit", Description: "Manage the currently active module, if the module is an exploit"},
	{Name: "module", Description: "Manage the currently loaded module."},
	{Name: "agent", Description: "Manage the currently active agent."},
	// {Name: "payload", Description: "Manage the currently active module, if the module is a payload"},
	// {Name: "hosts", Description: "Commands displaying hosts"},
	// {Name: "services", Description: "Commands displaying services"},
	// {Name: "creds", Description: "Commands displaying credentials"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Core commands help list
var coreCommands = []CommandDescription{
	{Name: "quit", Description: "Exit WireGost (without confirmation prompt)"},
	{Name: "config", Description: "Show environment configuration (personal directories, config files, etc...)"},
	{Name: "help", Description: "Show help categories (core, server, log, chat, workspace, stack, db, listeners, exploit, payload, etc...)"},
	{Name: "mode", Description: "Manage Input & Navigation mode (Vim or Emacs). Type 'mode' for displaying current mode"},
	{Name: "cd", Params: []string{"dir"}, Description: "Change current directory within WireGost. (Supports basic variable subsitution: ~/.././ )"},
	{Name: "!", Params: []string{"cmd"}, Description: "Execute a shell command through WireGost"},
	{Name: "resource.make", Params: []string{"file.rc", "int"}, Description: "Make a resource file. (name: resource name, int: number of history commands to save in file."},
	{Name: "resource.load", Params: []string{"file.rc"}, Description: "Load a resource file. (Completion based on files stored in the resource directory)"},
	{Name: "history.show", Params: []string{"int"}, Description: "Show last (int) number of commands used."},
}

//-------------------------------------------------------------------------------------------------------------------------
// Server commands
var serverCommands = []CommandDescription{
	{Name: "server.start", Description: "Start a server, in the current workspace, listening for agents. (Default parameters below are used)"},
	{Name: "server.stop", Description: "Stop the server."},
	{Name: "server.reload", Description: "Restart the server with parameters below. (Parameters will be saved for subsequent starts)"},
	{Name: "server.generate_certificate", Params: []string{"name"}, Description: "Generate a certificate and key pair to use with the server in this workspace."},
	{Name: "server.generate_jwt", Description: "(Optional) Generate a new JSON Web Token that agents dedicated to this workspace/server will use for authenticating."},
	{Name: "server.list", Description: "List all servers and their state (running, agents, etc.), regardless of the current workspace."},
}

var serverParams = []ParamDescription{
	{Name: "server.address", Description: "IP address on which the server will listen"},
	{Name: "server.port", Description: "Listening port of the server"},
	{Name: "server.protocol", Description: "The protocol (i.e. HTTP/2 or HTTP/3) the server will use (type 'h2' or 'h3')"},
	{Name: "server.certificate", Description: "Path to x.509 certificate needed for connection (default is the certificate in this workspace's directory)"},
	{Name: "server.key", Description: "Path to x.509 private key used for decrypting communications with agents"},
	{Name: "server.psk", Description: "The pre-shared key password used prior to Password Authenticated Key Exchange (PAKE)"},
	{Name: "server.jwt", Description: "JSON Web Token used for authenticating agents"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Log commands
var logCommands = []CommandDescription{
	{Name: "log.global", Params: []string{"int"}, Description: "Show last (int) lines of global log (aggregation of all specfic logs)"},
	{Name: "log.exploit", Params: []string{"int"}, Description: "Show last (int) lines of exploit-specific logs."},
	{Name: "log.payload", Params: []string{"int"}, Description: "Show last (int) lines of payload-specific logs."},
	{Name: "log.listeners", Params: []string{"int"}, Description: "Show last (int) lines of listeners-specific logs."},
	{Name: "log.db", Params: []string{"int"}, Description: "Show last (int) lines of database-specific logs."},
}

// Log parameters
var logParams = []ParamDescription{
	{"log.global.path", "Path to global logs", "~/.wiregost/logs/global.log"},
	{"log.exploit.path", "Path to exploit logs", "~/.wiregost/logs/exploit.log"},
	{"log.payload.path", "Path to payload logs", "~/.wiregost/logs/payload.log"},
	{"log.listeners.path", "Path to listeners logs", "~/.wiregost/logs/listeners.log"},
	{"log.db.path", "Path to Data Service logs", "~/.wiregost/logs/db.log"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Server commands
// These commands are used for connecting clients to a given Wiregost instance, and for admin tasks

var endpointCommands = []CommandDescription{
	{Name: "endpoint.show", Description: "Show Server information (registered endpoints, client access token, etc...)"},
	{Name: "endpoint.generate_token", Description: "Ask the endpoint to generate a new client-side access token."},
	{Name: "endpoint.connect", Params: []string{"name"}, Description: "Connect to one of the saved WireGost endpoints."},
	{Name: "endpoint.add", Description: "Add endpoint based on the current value of parameters"},
}

var endpointParams = []ParamDescription{
	{Name: "endpoint.address", Description: "IP or resolved address of the endpoint"},
	{Name: "endpoint.port", Description: "Listening port of the endpoint"},
	{Name: "endpoint.name", Description: "Name under which this endpoint will be saved or displayed"},
	{Name: "endpoint.certificate", Description: "Path to certificate needed for connection (not required)"},
	{Name: "endpoint.default", Description: "Make this endpoint the default endpoint to connect to when client is started"},
}

var endpointAdminCommands = []CommandDescription{
	{Name: "endpoint.admin.show", Description: "Show registered users (and if they are active)"},
	{Name: "endpoint.admin.add_user", Params: []string{"name"}, Description: "Register a new user (password will be sent first connection)"},
	{Name: "endpoint.admin.delete_user", Params: []string{"name"}, Description: "Delete one or more registered users"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Chat commands
var chatCommands = []CommandDescription{
	{Name: "chat", Params: []string{"all | user"}, Description: "Send to all|user a string message."},
	{Name: "chat.show", Params: []string{"int"}, Description: "Show last (int) number of messages for this session"},
	{Name: "chat.connected", Description: "Show all currently connected users"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Workspace commands
var workspaceCommands = []CommandDescription{
	{Name: "workspace.show", Description: "Show all workspaces"},
	{Name: "workspace.switch", Description: "Switch to workspace. (Parameters, variables and modules in the current workspace will be automatically saved "},
	{Name: "workspace.new", Params: []string{"name"}, Description: "Create new workspace with specified name."},
	{Name: "workspace.delete", Params: []string{"name"}, Description: "Delete the specified workspace. (Completed)"},
}

var workspaceParams = []ParamDescription{
	{Name: "workspace.name", Description: "Workspace name"},
	{Name: "workspace.description", Description: "Description for the current workspace"},
	{Name: "workspace.boundary", Description: "Network address/range in which activity is allowed for this workspace."},
	{Name: "workspace.owner", Description: "Owner ID for this workspace"},
	{Name: "workspace.limit", Description: "false"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Stack commands
var stackCommands = []CommandDescription{
	{Name: "stack.show", Params: []string{"all | module_name"}, Description: "Show one or some modules currently loaded in the Module Stack"},
	{Name: "stack.pop", Params: []string{"all | module_name"}, Description: "Unload one or more modules from the Module Stack"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Module commands
var moduleCommands = []CommandDescription{
	{Name: "back", Description: "Exit from current module. (Doesn't unload it from stack)"},
	{Name: "show", Params: []string{"info | options"}, Description: "Show information about a module or its options."},
	{Name: "info", Description: "Show information about a module"},
	{Name: "reload", Description: "Reloads the module to a fresh clean state"},
	{Name: "set", Params: []string{"<option name>", "<option value>"}, Description: "Set the value for one of the module's options. (Auto-completed options)"},
	{Name: "run", Description: "Run or execute the module"},
}

//-------------------------------------------------------------------------------------------------------------------------
// Agent commands
var agentCommands = []CommandDescription{
	{Name: "info", Description: "Display all information about the agent"},
	{Name: "back | main", Description: "Return to the main menu"},
	{Name: "status", Description: "Print the current status of the agent."},
	{Name: "cd", Params: []string{"../.. | c:\\\\Users"}, Description: "Change directories in the agent's target system."},
	{Name: "ls", Params: []string{"/etc | c:\\\\Users"}, Description: "List directory contents"},
	{Name: "pwd", Description: "Print the current working directory in the target."},
	{Name: "cmd", Params: []string{"ping -c 3 8.8.8.8"}, Description: "Execute a command on the agent."},
	// We have renamed the shell as "cmd" because it is less ambiguous. Think of changing the corresponding handler if needed.
	// {"shell", "Execute a command on the agent", "shell ping -c 3 8.8.8.8"},
	{Name: "set", Params: []string{"<option name>", "<option value>"}, Description: "Set the value for one of the agent's options. (Auto-completed options)"},
	// Maybe useful to check if we want to autocomplete available options for agent
	//{"set", "Set the value for one of the agent's options", "killdate, maxretry, padding, skew, sleep"},
	{Name: "download", Params: []string{"remote_file"}, Description: "Download a file from the agent's target."},
	{Name: "upload", Params: []string{"local_file", "remote_file"}, Description: "Upload a file to the agent's target."},
	{Name: "execute-shellcode", Params: []string{"self, remote <pid> | RtlCreateUserThread <pid>"}, Description: "Execute shellcode on the target."},
	{Name: "kill", Description: "Instruct the agent to die or quit."},
}

//------------------------------------------------------------------------------------------------------------------------
// Database commands
// var dbCommands = []CommandDescription{
//     {Name: "db.source.connect", Params: {"source"}, Description: "Connect to data service. (Completion based on saved sources.)"},
//     {Name: "db.source.show", Description: "Show all data service sources."},
//     {Name: "db.source.add", Description: "Add data service source based on the current value of parameters below."},
//     {Name: "db.source.delete", Params: {"source"}, Description: "Delete one or more data service sources."},
//     {Name: "file.import", Description: "Import entities from a file."},
//     {Name: "file.export", Description: "Export the current data service content to a file. (Specify file in parameters below)"},
// }
//
// // Database parameters
// var dbParams = []ParamDescription{
//     {"db.source.name", "Data Service name"},
//     {"db.source.url", "Path to Data Service API."},
//     {"db.source.password", "Data Service user password"},
//     {"db.source.user", "Data Service user name"},
//     {"db.source.token", "Data Service user token"},
//     {"file.import.path", "Path to import file"},
//     {"file.export.path", "Path to export file", "~/.wiregost/db/exports/"},
//     {"file.export.format", "Format in which to export the Database content.", "json"},
// }
//
// //-------------------------------------------------------------------------------------------------------------------------
// // Entities commands
// var hostsCommands = []CommandDescription{
//     {"hosts", Description: "Show all hosts"},
//     {"hosts.ip", {"ip"}, Description: "Show hosts matching the provided IP address"},
//     {"hosts.mac", {"mac"}, Description: "Show hosts matching the provided MAC address"},
//     {"hosts.os", {"os"}, Description: "Show hosts matching the provided Operating System"},
//     {"hosts.name", {"name"}, Description: "Show hosts matching the provided Name"},
//     {"hosts.purpose", {"purpose"}, Description: "Show hosts matching the provided device purpose"},
//     {"hosts.filter", {"string"}, Description: "Show hosts matching the provided (string) filter"},
// }
//
// var servicesCommands = []CommandDescription{
//     {"services", Description: "Show all services"},
//     {"services.port", {"port"}, "Show all services listening on this port"},
//     {"services.host", {"host"}, "Show all services served on this IP Address"},
//     {"services.proto", {"proto"}, "Show all services communicating on this transport protocol"},
//     {"services.name", {"name"}, "Show all services communicating on this application-layer protocol"},
//     {"services.info", {"info"}, "Show all services matching this banner info"},
//     {"services.state", {"state"}, "Show all services matching this listening state"},
//     {"services.filter", {"string"}, "Show all services matching the provided (string) filter"},
// }
//
// var credsCommands = []CommandDescription{
//     {"creds", Description: "Show all credentials"},
//     {"creds.type", {"type"}, "Show credentials matching this type"},
//     {"creds.priv", {"priv"}, "Show all credentials matching this private type"},
//     {"creds.filter", {"string"}, "Show all credentials matching the provided (string) filter."},
// }
//
// //-------------------------------------------------------------------------------------------------------------------------
// // Editor commands
// var editorCommands = []CommandDescription{
//     {"edit", Description: "Edit the currently active module's code"},
//     {"loadpath", {"string"}, "Load the provided (string) path for modules"},
//     {"reload", {"all", "lib", "modules"}, "Reload all, libraries or modules code."},
// }
//
// //-------------------------------------------------------------------------------------------------------------------------
// // Listeners commands
// var listenersCommands = []CommandDescription{
//     {"listeners.show", {"all", "listener_name"}, "Show one or more listeners"},
//     {"listeners.kill", {"all", "listener_name"}, "Kill one or more listeners"},
//     {"listeners.rename", {"current", "new"}, "Rename (current) listeners with (new) name"},
//     {"listeners.duplicate", {"listener_name"}, "Duplicate (listener_name) and launch it"},
// }
