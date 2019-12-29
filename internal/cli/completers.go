package cli

import (
	"strings"

	"github.com/chzyer/readline"
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
			// readline.PcItem("hosts"),
			// readline.PcItem("services"),
			// readline.PcItem("creds"),
			readline.PcItem("agent"),
			readline.PcItem("module"),
			readline.PcItem("compiler"),
			// readline.PcItem("exploit"),
			// readline.PcItem("payload"),
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

		// Compiler
		readline.PcItem("compiler"),

		// Server
		readline.PcItem("server",
			readline.PcItem("connect"), // Add getServerList here
			readline.PcItem("list"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("debug"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("exploit"),
				readline.PcItem("agent"),
			),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("show"), // Add getStackList here
			readline.PcItem("pop",
				readline.PcItemDynamic(s.ListStackModules())), // Same
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("list"),
			readline.PcItem("switch",
				readline.PcItemDynamic(s.ListWorkspaces())),
			readline.PcItem("new"),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.ListWorkspaces())),
		),

		// Agent
		readline.PcItem("agent",
			readline.PcItem("list"),     // Add getAgentsList here
			readline.PcItem("interact"), // same
			readline.PcItem("remove"),   // same
		),
		readline.PcItem("interact"), // Same

		// Module
		readline.PcItem("use",
			readline.PcItem("module",
				readline.PcItemDynamic(s.ListModules())),
		),
		readline.PcItem("set",
			readline.PcItemDynamic(s.ListParams()),
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
			// readline.PcItem("hosts"),
			// readline.PcItem("services"),
			// readline.PcItem("creds"),
			readline.PcItem("agent"),
			readline.PcItem("module"),
			// readline.PcItem("exploit"),
			// readline.PcItem("payload"),
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

		// Compiler
		readline.PcItem("compiler"),

		// Server
		readline.PcItem("server",
			readline.PcItem("connect"), // Add getServerList here
			readline.PcItem("list"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("debug"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("exploit"),
				readline.PcItem("agent"),
			),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("show"), // Add getStackList here
			readline.PcItem("pop",
				readline.PcItemDynamic(s.ListStackModules())), // Same
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("list"),
			readline.PcItem("switch",
				readline.PcItemDynamic(s.ListWorkspaces())),
			readline.PcItem("new"),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.ListWorkspaces())),
		),

		// Agent
		readline.PcItem("agent",
			readline.PcItem("list"),     // Add getAgentsList here
			readline.PcItem("interact"), // same
			readline.PcItem("remove"),   // same
		),
		readline.PcItem("interact"), // Same

		// Module
		readline.PcItem("use",
			readline.PcItem("module",
				readline.PcItemDynamic(s.ListModules())),
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
				readline.PcItem("all"), // add getAgentsList here
			),
			readline.PcItemDynamic(s.GetModuleOptions()),
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

		// Compiler // ADD SPECIAL HANDLING CODE FOR MANAGING SHELL STATE HERE
		// readline.PcItem("compiler"),

		// Server
		readline.PcItem("server",
			readline.PcItem("connect"), // Add getServerList here
			readline.PcItem("list"),
		),

		// Log
		readline.PcItem("log",
			readline.PcItem("level",
				readline.PcItem("debug"),
			),
			readline.PcItem("show",
				readline.PcItem("all"),
				readline.PcItem("exploit"),
				readline.PcItem("agent"),
			),
		),

		// Module Stack
		readline.PcItem("stack",
			readline.PcItem("show"), // Add getStackList here
			readline.PcItem("pop",
				readline.PcItemDynamic(s.ListStackModules())), // Same
		),

		// Workspace
		readline.PcItem("workspace",
			readline.PcItem("list"),
			readline.PcItem("new"),
			readline.PcItem("delete"), // Same
		),

		// Agent
		readline.PcItem("cmd"),
		readline.PcItem("back"),
		readline.PcItem("download"),
		readline.PcItem("execute-shellcode",
			readline.PcItem("self"),
			readline.PcItem("remote"),
			readline.PcItem("RtlCreateUserThread"),
		),
		readline.PcItem("info"),
		readline.PcItem("kill"),
		readline.PcItem("main"),
		readline.PcItem("shell"),
		readline.PcItem("set",
			readline.PcItem("maxretry"),
			readline.PcItem("padding"),
			readline.PcItem("skew"),
			readline.PcItem("sleep"),
		),
		readline.PcItem("upload"),
	)

	var compiler = readline.NewPrefixCompleter(
		readline.PcItem("compile"),
		readline.PcItem("back"),
		readline.PcItem("help"),
		readline.PcItem("list",
			readline.PcItem("servers"),
			readline.PcItem("parameters")),
		readline.PcItem("set", readline.PcItemDynamic(s.GetCompilerOptions())),
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
func (s *Session) ListParams() func(string) (names []string) {
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
			"endpoint.default",
			// Workspace
			"workspace.name",
			"workspace.description",
			"workspace.boundary",
			"workspace.owner",
			"workspace.limit",
		}
		return sessionParams
	}
}

func (s *Session) ListWorkspaces() func(string) (names []string) {
	return func(string) []string {
		s.Send([]string{"workspace", "list"})
		workspace := <-workspaceReqs
		var list []string
		// Handle change of state here
		for _, ws := range workspace.WorkspaceInfos {
			list = append(list, ws[0])
		}
		return list
	}
}

func (s *Session) ListModules() func(string) (names []string) {
	return func(string) []string {
		s.Send([]string{"module", "list"})
		resp := <-moduleReqs
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

func (s *Session) ListStackModules() func(string) (names []string) {
	return func(string) []string {
		s.Send([]string{"stack", "list"})
		resp := <-moduleReqs
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

func (s *Session) GetModuleOptions() func(string) (options []string) {
	return func(string) []string {
		s.Send([]string{"show", "options"})
		mod := <-moduleReqs
		opts := mod.Modules[0]
		list := make([]string, 0)
		for _, opt := range opts.Options {
			list = append(list, opt.Name)
		}
		return list
	}
}

func (s *Session) GetCompilerOptions() func(string) (options []string) {
	return func(string) []string {
		s.Send([]string{"list", "parameters"})
		comp := <-compilerReqs
		opts := comp.Options
		list := make([]string, 0)
		for _, opt := range opts {
			list = append(list, opt.Name)
		}
		return list
	}
}
