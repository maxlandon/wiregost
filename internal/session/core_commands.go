package session

// This file contains all core command handlers and their registering function.

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/maxlandon/wiregost/internal/session/core"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

// Exit Session
func (s *Session) exitHandler(args []string, sess *Session) error {
	// Should also empty the Module Stack, and potentially save it if possible.
	s.Active = false
	s.Input.Close()
	return nil
}

// Show WireGost version, and potentially other version info
func (s *Session) versionHandler(args []string) {

}

// Resource commands: Load resource, and Make resource.
func (s *Session) resourceLoadHandler(args []string, sess *Session) error {
	name := ""
	// If resource.load: 1 command, 1 argument
	if len(args) == 2 {
		name = str.Trim(args[1])
	}

	filestr, _ := fs.Expand(s.Config.ResourceDir + name)
	file, _ := os.Open(filestr)
	if filepath.Ext(filestr) != ".rc" {
		fmt.Errorf("%sError: file must be a configuration (.rc) file.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.Run(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("%sError parsing resource command: %s%s\n", tui.RED, scanner.Text, tui.RESET)
		log.Fatal(err)
	}

	return nil
}

func (s *Session) resourceMakeHandler(args []string, sess *Session) error {
	name := ""
	nb := 0
	// If resource.load: 1 command, 2 arguments
	if len(args) == 3 {
		name = str.Trim(args[1])
		nb, _ = strconv.Atoi(str.Trim(args[2]))
	} else {
		fmt.Printf("%sError: missing arguments.%s\n", tui.RED, tui.RESET)
		return nil
	}

	// Check if resource already exists
	file, _ := fs.Expand(s.Config.ResourceDir + name)
	if fs.Exists(file) {
		fmt.Printf("%sError: resource file already exists. Cannot overwrite it.%s\n", tui.RED, tui.RESET)
		return nil
	}

	// If not, create it
	res, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	// Find file length, for subsequent selection
	hist, _ := os.Open(s.Config.HistoryFile)
	hlength := 0
	scanner := bufio.NewScanner(hist)
	for scanner.Scan() {
		hlength += 1
	}
	hist.Close()

	// Select last x commands in file and output them to resource file
	hist, _ = os.Open(s.Config.HistoryFile)
	defer hist.Close()
	scan := bufio.NewScanner(hist)
	for scan.Scan() {
		if hlength <= nb {
			res.Write([]byte(scan.Text() + "\n"))
		}
		hlength -= 1
	}
	fmt.Printf("%sResource file created and filed with last %s commands.%s\n", tui.GREEN, strconv.Itoa(nb), tui.RESET)

	return nil
}

// Navigation mode: Vim or Emacs.
func (s *Session) setModeHandler(args []string, sess *Session) error {
	filter := ""
	if len(args) == 2 {
		filter = str.Trim(args[1])
	}

	if filter == "" {
		if s.Input.IsVimMode() {
			println("Current mode: " + tui.Yellow("Vim"))
		} else {
			println("Current mode: " + tui.Yellow("Emacs"))
		}
	}

	if filter == "vim" {
		s.Input.SetVimMode(true)
		println("Switched mode: " + tui.Yellow("Vim"))
	}
	if filter == "emacs" {
		s.Input.SetVimMode(false)
		println("Switched mode: " + tui.Yellow("Emacs"))
	}
	return nil
}

// Shell Command handler
func (s *Session) shellHandler(args []string, sess *Session) error {
	out, err := core.Shell(args[0])
	if err == nil {
		fmt.Printf("%s\n", out)
	}
	return err
}

func (s *Session) changeDirHandler(args []string, sess *Session) error {
	raw := str.Trim(args[1])
	dir, _ := fs.Expand(raw)
	os.Chdir(dir)
	s.CurrentDir = dir

	return nil
}

// Register all handlers defined above
func (s *Session) registerCoreHandlers() {
	// Exit
	s.addHandler(NewCommandHandler("exit",
		"^(q|quit|e|exit)$",
		"Close the session and exit.",
		s.exitHandler),
		readline.PcItem("exit"))

	// Shell
	s.addHandler(NewCommandHandler("! COMMAND",
		"^!\\s*(.+)$",
		"Execute a shell command and print its output.",
		s.shellHandler),
		readline.PcItem("!"))

	// Input Mode
	s.addHandler(NewCommandHandler("mode INPUT",
		"^(mode|\\?)(.*)$",
		"Change navigation and input mode (Vim or Emacs)",
		s.setModeHandler),
		readline.PcItem("mode", readline.PcItemDynamic(func(prefix string) []string {
			prefix = str.Trim(prefix[4:])
			modes := []string{"vim", "emacs"}
			return modes
		})))

	// Load resource
	s.addHandler(NewCommandHandler("resource.load",
		"^(resource.load|\\?)(.*)$",
		"Load resource file.",
		s.resourceLoadHandler),
		readline.PcItem("resource.load", readline.PcItemDynamic(func(prefix string) []string {
			prefix = str.Trim(prefix[13:])
			files, _ := ioutil.ReadDir(s.Config.ResourceDir)

			var resources []string
			for _, file := range files {
				if filepath.Ext(file.Name()) == ".rc" {
					resources = append(resources, file.Name())
				}
			}
			return resources
		})))

	// Make resource
	s.addHandler(NewCommandHandler("resource.make",
		`^(resource.make)\s+([^\s]+)\s+(.+)$`,
		"Make resource file.",
		s.resourceMakeHandler),
		readline.PcItem("resource.make"))

	// Change directory
	s.addHandler(NewCommandHandler("cd DIRECTORY",
		"^(cd|\\?)(.*)$",
		// "^(cd)\\s*(.+)$",
		"Change the current directory within WireGost",
		s.changeDirHandler),
		readline.PcItem("cd"))
}
