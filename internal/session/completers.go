package session

import (
	// Standard
	"strconv"
	"strings"

	// 3rd party
	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

func (s *Session) getCompleter(completer string) *readline.PrefixCompleter {

	// ------------------------------------------------------------
	// COMMANDS

	// Main menu.
	var main = readline.NewPrefixCompleter(
		// Core
		readline.PcItem("help",
			readline.PcItem("core"),
			readline.PcItem("log"),
			readline.PcItem("server"),
			readline.PcItem("endpoint"),
			readline.PcItem("workspace"),
			readline.PcItem("stack"),
			readline.PcItem("agent"),
			readline.PcItem("module"),
			readline.PcItem("compiler"),
		),
		readline.PcItem("mode",
			readline.PcItem("vim"),
			readline.PcItem("emacs"),
		),
		readline.PcItem("history",
			readline.PcItem("show"),
		),
		readline.PcItem("resource",
			readline.PcItem("make"),
			readline.PcItem("load"),
		),
		readline.PcItem("cd"),
		readline.PcItem("!"),
		readline.PcItem("exit"),
		readline.PcItem("get", readline.PcItemDynamic(s.listParams())),

		// Endpoint
		readline.PcItem("endpoint",
			readline.PcItem("connect", readline.PcItemDynamic(s.getEndpointList())),
			readline.PcItem("list"),
			readline.PcItem("delete", readline.PcItemDynamic(s.getEndpointList())),
			readline.PcItem("add"),
		),

		// Compiler
		readline.PcItem("compiler"),

		// Server
		readline.PcItem("server",
			readline.PcItem("start"), // Add getServerList here
			readline.PcItem("stop"),
			readline.PcItem("list"),
			readline.PcItem("reload"),
			readline.PcItem("generate_certificate"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("trace"),
				readline.PcItem("debug"),
				readline.PcItem("info"),
				readline.PcItem("warning"),
				readline.PcItem("error"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("agent", readline.PcItemDynamic(s.getAgentIds())),
				readline.PcItem("server"),
			),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("use",
				readline.PcItemDynamic(s.listStackModules())),
			readline.PcItem("show"),
			readline.PcItem("pop",
				readline.PcItemDynamic(s.listStackModules())),
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("list"),
			readline.PcItem("switch",
				readline.PcItemDynamic(s.listWorkspaces())),
			readline.PcItem("new"),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.listWorkspaces())),
		),
		// Agent
		readline.PcItem("agent",
			readline.PcItem("list"),
			readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),
			readline.PcItem("remove", readline.PcItemDynamic(s.getAgentList())),
		),
		readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),

		// Module
		readline.PcItem("use",
			readline.PcItem("module",
				readline.PcItemDynamic(s.listModules())),
		),
		// General parameters
		readline.PcItem("set",
			readline.PcItemDynamic(s.listParams()),
		),
	)

	var module = readline.NewPrefixCompleter(
		// Core
		readline.PcItem("help",
			readline.PcItem("core"),
			readline.PcItem("log"),
			readline.PcItem("server"),
			readline.PcItem("endpoint"),
			readline.PcItem("workspace"),
			readline.PcItem("stack"),
			readline.PcItem("agent"),
			readline.PcItem("module"),
		),
		readline.PcItem("mode",
			readline.PcItem("vim"),
			readline.PcItem("emacs"),
		),
		readline.PcItem("history",
			readline.PcItem("show"),
		),
		readline.PcItem("resource",
			readline.PcItem("make"),
			readline.PcItem("load"),
		),
		readline.PcItem("cd"),
		readline.PcItem("!"),
		readline.PcItem("exit"),
		readline.PcItem("get", readline.PcItemDynamic(s.listParams())),

		// Endpoint
		readline.PcItem("endpoint",
			readline.PcItem("connect", readline.PcItemDynamic(s.getEndpointList())),
			readline.PcItem("list"),
			readline.PcItem("delete", readline.PcItemDynamic(s.getEndpointList())),
			readline.PcItem("add"),
		),

		// Compiler
		readline.PcItem("compiler"),

		// Server
		readline.PcItem("server",
			readline.PcItem("start"), // Add getServerList here
			readline.PcItem("stop"),
			readline.PcItem("list"),
			readline.PcItem("reload"),
			readline.PcItem("generate_certificate"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("trace"),
				readline.PcItem("debug"),
				readline.PcItem("info"),
				readline.PcItem("warning"),
				readline.PcItem("error"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("server"),
				readline.PcItem("agent", readline.PcItemDynamic(s.getAgentIds())),
			),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("use",
				readline.PcItemDynamic(s.listStackModules())),
			readline.PcItem("show"),
			readline.PcItem("pop",
				readline.PcItemDynamic(s.listStackModules())),
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("list"),
			readline.PcItem("switch",
				readline.PcItemDynamic(s.listWorkspaces())),
			readline.PcItem("new"),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.listWorkspaces())),
		),

		// Agent
		readline.PcItem("agent",
			readline.PcItem("list"),
			readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),
			readline.PcItem("remove", readline.PcItemDynamic(s.getAgentList())),
		),
		readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),

		// Module
		readline.PcItem("use",
			readline.PcItem("module",
				readline.PcItemDynamic(s.listModules())),
		),
		readline.PcItem("info"),
		readline.PcItem("reload"),
		readline.PcItem("run"),
		readline.PcItem("back"),
		readline.PcItem("show",
			readline.PcItem("options"),
			readline.PcItem("info"),
		),
		readline.PcItem("set",
			readline.PcItem("agent",
				readline.PcItem("all"),
				readline.PcItemDynamic(s.getAgentList()),
			),
			readline.PcItemDynamic(s.getModuleOptions()),
		),
	)

	var agent = readline.NewPrefixCompleter(
		// Core
		readline.PcItem("help",
			readline.PcItem("core"),
			readline.PcItem("log"),
			readline.PcItem("server"),
			readline.PcItem("endpoint"),
			readline.PcItem("workspace"),
			readline.PcItem("stack"),
			readline.PcItem("hosts"),
			readline.PcItem("services"),
			readline.PcItem("creds"),
			readline.PcItem("agent"),
			readline.PcItem("module"),
			readline.PcItem("exploit"),
			readline.PcItem("payload"),
		),
		readline.PcItem("mode",
			readline.PcItem("vim"),
			readline.PcItem("emacs"),
		),
		readline.PcItem("history",
			readline.PcItem("show"),
		),
		readline.PcItem("resource",
			readline.PcItem("make"),
			readline.PcItem("load"),
		),
		readline.PcItem("cd"),
		readline.PcItem("!"),
		readline.PcItem("exit"),

		// Server
		readline.PcItem("server",
			readline.PcItem("start"), // Add getServerList here
			readline.PcItem("stop"),
			readline.PcItem("list"),
			readline.PcItem("reload"),
			readline.PcItem("generate_certificate"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("trace"),
				readline.PcItem("debug"),
				readline.PcItem("info"),
				readline.PcItem("warning"),
				readline.PcItem("error"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("server"),
				readline.PcItem("agent", readline.PcItemDynamic(s.getAgentIds())),
			),
		),

		// Module
		readline.PcItem("use",
			readline.PcItem("module",
				readline.PcItemDynamic(s.listModules())),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("use",
				readline.PcItemDynamic(s.listStackModules())),
			readline.PcItem("show"),
			readline.PcItem("pop",
				readline.PcItemDynamic(s.listStackModules())),
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("switch",
				readline.PcItemDynamic(s.listWorkspaces())),
			readline.PcItem("list"),
			readline.PcItem("new"),
			readline.PcItem("delete", readline.PcItemDynamic(s.listWorkspaces())),
		),

		// Agent
		readline.PcItem("back"),
		readline.PcItem("main"),
		readline.PcItem("info"),
		readline.PcItem("kill"),
		readline.PcItem("ls"),
		readline.PcItem("cd"),
		readline.PcItem("pwd"),
		readline.PcItem("cmd"),
		readline.PcItem("shell"),
		readline.PcItem("download"),
		readline.PcItem("upload"),
		readline.PcItem("execute-shellcode",
			readline.PcItem("self"),
			readline.PcItem("remote"),
			readline.PcItem("RtlCreateUserThread"),
			readline.PcItem("UserAPC"),
		),
		readline.PcItem("set",
			readline.PcItem("killdate"),
			readline.PcItem("maxretry"),
			readline.PcItem("padding"),
			readline.PcItem("skew"),
			readline.PcItem("sleep"),
		),
		readline.PcItem("agent",
			readline.PcItem("list"),
			readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),
		),
		readline.PcItem("interact", readline.PcItemDynamic(s.getAgentList())),
	)

	var compiler = readline.NewPrefixCompleter(
		readline.PcItem("compile"),
		readline.PcItem("back"),
		readline.PcItem("help"),
		readline.PcItem("list",
			readline.PcItem("servers"),
			readline.PcItem("parameters")),
		readline.PcItem("set", readline.PcItemDynamic(s.getCompilerOptions())),
		readline.PcItem("use"), // Add server completion function here
	)

	// ------------------------------------------------------------
	// PARAMETERS

	switch completer {
	case "main":
		return main
	case "module":
		return module
	case "agent":
		return agent
	case "compiler":
		return compiler
	}

	return main
}

// DYNAMIC COMPLETER FUNCTIONS
func (s *Session) listParams() func(string) (names []string) {
	return func(string) []string {
		sessionParams := []string{
			// Server
			"server.address",
			"server.port",
			"server.protocol",
			"server.certificate",
			"server.key",
			"server.psk",
			"server.jwt",
			// Endpoint
			"endpoint.address",
			"endpoint.port",
			"endpoint.name",
			"endpoint.certificate",
			"endpoint.key",
			"endpoint.default",
			// Workspace
			"workspace.description",
			"workspace.boundary",
			"workspace.limit",
		}
		return sessionParams
	}
}

func (s *Session) listWorkspaces() func(string) (names []string) {
	return func(string) []string {
		s.send([]string{"workspace", "list"})
		workspace := <-s.workspaceReqs
		var list []string
		// Handle change of state here
		for _, ws := range workspace.WorkspaceInfos {
			list = append(list, ws[0])
		}
		return list
	}
}

func (s *Session) listModules() func(string) (names []string) {
	return func(string) []string {
		s.send([]string{"module", "list"})
		resp := <-s.moduleReqs
		list := resp.ModuleList
		// This is useless, but we should devise way to recursively update paths
		// so that we do not display all modules at once during completion.
		var testList []string
		for _, mod := range list {
			m := strings.TrimRight(mod, "/")
			testList = append(testList, m)
		}
		return testList
	}
}

func (s *Session) listStackModules() func(string) (names []string) {
	return func(string) []string {
		s.send([]string{"stack", "list"})
		resp := <-s.moduleReqs
		list := resp.ModuleList
		// This is useless, but we should devise way to recursively update paths
		// so that we do not display all modules at once during completion.
		var testList []string
		for _, mod := range list {
			m := strings.TrimRight(mod, "/")
			testList = append(testList, m)
		}
		return testList
	}
}

func (s *Session) getModuleOptions() func(string) (options []string) {
	return func(string) []string {
		s.send([]string{"show", "options"})
		mod := <-s.moduleReqs
		opts := mod.Modules[0]
		list := make([]string, 0)
		for _, opt := range opts.Options {
			list = append(list, opt.Name)
		}
		return list
	}
}

func (s *Session) getCompilerOptions() func(string) (options []string) {
	return func(string) []string {
		s.send([]string{"list", "parameters"})
		comp := <-s.compilerReqs
		opts := comp.Options
		list := make([]string, 0)
		for _, opt := range opts {
			list = append(list, opt.Name)
		}
		return list
	}
}

func (s *Session) getEndpointList() func(string) (endpoints []string) {
	return func(string) []string {
		endpoints := []string{}
		for _, l := range s.SavedEndpoints {
			server := l.FQDN + " " + tui.Dim("at ") + l.IPAddress + ":" + strconv.Itoa(l.Port)
			endpoints = append(endpoints, server)
		}
		return endpoints

	}
}

func (s *Session) getAgentList() func(string) (agents []string) {
	return func(string) []string {
		s.send([]string{"agent", "show"})
		agents := <-s.agentReqs
		list := make([]string, 0)
		for _, a := range agents.Infos {
			agent := a["id"] + " " + tui.Dim("as "+a["username"]) + tui.Bold("@") + tui.Dim(a["hostname"])
			list = append(list, agent)
		}
		return list
	}
}

func (s *Session) getAgentIds() func(string) (agents []string) {
	return func(string) []string {
		s.send([]string{"agent", "show"})
		agents := <-s.agentReqs
		list := make([]string, 0)
		for _, a := range agents.Infos {
			list = append(list, a["id"])
		}
		return list
	}
}
