package session

// This file contains all the descriptions for each command set, along with their respective parameters if they have some.

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
	{Name: "server", Description: "Commands and parameters for managing WireGost Server (connection, add, generate tokens, etc)"},
	{Name: "log", Description: "Commands for managing the various sets of logs used by WireGost"},
	{Name: "chat", Description: "Commands for using WireGost's messaging system"},
	{Name: "workspace", Description: "Manage WireGost workspaces"},
	{Name: "stack", Description: "Manage the module stack (all modules currently loaded in this session)"},
	{Name: "global", Description: "Manage all global variables in WireGost"},
	{Name: "db", Description: "Manage WireGost data services"},
	{Name: "listeners", Description: "Manage listeners instantiated in this session"},
	{Name: "exploit", Description: "Manage the currently active module, if the module is an exploit"},
	{Name: "payload", Description: "Manage the currently active module, if the module is a payload"},
	{Name: "hosts", Description: "Commands displaying hosts"},
	{Name: "services", Description: "Commands displaying services"},
	{Name: "creds", Description: "Commands displaying credentials"},
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
	{Name: "server.show", Description: "Show Server information (registered servers, client access token, etc...)"},
	{Name: "server.generate_token", Description: "Ask the server to generate a new client-side access token."},
	{Name: "server.connect", Params: []string{"name"}, Description: "Connect to one of the saved WireGost servers."},
	{Name: "server.add", Description: "Add server based on the current value of parameters"},
}

var serverParams = []ParamDescription{
	{Name: "server.address", Description: "IP or resolved address of the server"},
	{Name: "server.port", Description: "Listening port of the server"},
	{Name: "server.name", Description: "Name under which this server will be saved or displayed"},
	{Name: "server.certificate", Description: "Path to certificate needed for connection (not required)"},
	{Name: "server.default", Description: "Make this server the default server to connect to when client is started"},
}

var serverAdminCommands = []CommandDescription{
	{Name: "server.admin.show", Description: "Show registered users (and if they are active)"},
	{Name: "server.admin.add_user", Params: []string{"name"}, Description: "Register a new user (password will be sent first connection)"},
	{Name: "server.admin.delete_user", Params: []string{"name"}, Description: "Delete one or more registered users"},
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
	{Name: "workspace.created_at", Description: "Created at"},
	{Name: "workspace.uptated_at", Description: "Updated at "},
}

//-------------------------------------------------------------------------------------------------------------------------
// Stack commands
var stackCommands = []CommandDescription{
	{Name: "stack.show", Params: []string{"all | module_name"}, Description: "Show one or some modules currently loaded in the Module Stack"},
	{Name: "stack.pop", Params: []string{"all | module_name"}, Description: "Unload one or more modules from the Module Stack"},
}

//-------------------------------------------------------------------------------------------------------------------------
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
