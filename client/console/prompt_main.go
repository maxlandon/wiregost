package console

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/readline"

	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/util"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

// mainprompt - A prompt used when the user is in the main menu, with modules loaded or not.
type mainprompt struct {
	Base          string // The left-most side of the prompt (usually server address)
	CustomBase    string // A user-provided custom base
	Module        string // Contains all information on current modules/sub loaded.
	ContextPrompt string // The right-most side of the prompt (workspaces, sessions, jobs, etc.)
	OptionPrompt  string // An option is set (with error or not), action results are stored here.

	Maincallbacks map[string]func() string // All values are automatically calculated
	MainColors    map[string]string        // All values are automatically calculated
}

// Render - The main menu prompt outputs its computed content
func (m *mainprompt) Render() (prompt string) {

	// Get term width: from this will depend some formatting
	sWidth := readline.GetTermWidth()

	// Other lengths that we might need to pass around
	var bWidth int // base prompt width, after computing.
	var mWidth int // Module prompt width, after computing.
	var cWidth int // Context prompt width

	// Compute all prompts
	m.Base, bWidth = m.computeBase()
	m.Module, mWidth = m.computeModule(sWidth)
	m.ContextPrompt, cWidth = m.computeContext(mWidth, sWidth)

	// Here, we have a prompt too large for the screen, it should not.
	// Truncate module prompt with the exceding length, or context prompt.
	if bWidth+mWidth+cWidth > sWidth {
		// m.Module = truncate()
	}

	// Get the empty part and pad with it, before
	pad := m.getPromptPad(sWidth, bWidth, mWidth, cWidth)

	// Return the prompt
	prompt = m.Base + m.Module + pad + m.ContextPrompt

	return
}

// computeModule - Computes the base prompt (left-side) with potential custom prompt given.
func (m *mainprompt) computeBase() (p string, width int) {
	p += tui.RESET // Always

	if m.CustomBase != "" {
		p = m.CustomBase
	} else {
		p = "{bddg}@{server_ip} "
	}

	p += tui.RESET // Always at end

	for ok, cb := range m.Maincallbacks {
		p = strings.Replace(p, ok, cb(), 1)
	}
	for tok, color := range m.MainColors {
		p = strings.Replace(p, tok, color, -1)
	}

	width = getRealLength(p)

	return
}

// computeModule - Analyses which modules and their subtypes are currently loaded on the console.
// Because this string might be long and we still have things to print to its right-side, we may
// have to adapt the length of our module prompt output.
func (m *mainprompt) computeModule(sWidth int) (p string, width int) {
	p += tui.RESET // Always

	// If current module is empty, we query and show information on the stack.
	if cctx.Module.Info.Name == "" && cctx.Module.Info.Path == "" {
		p += "{dim} -> stack: {r}{stack_len} module"
		if m.Maincallbacks["{stack_len}"]() != strconv.Itoa(0) {
			p += "s"
		}
		goto callbacks // We compute the prompt and return
	}

	p += "{dim}=> {fw}{module_type}{fw}({r}" // First part of the module prompt

	// We first filter based on type: some modules accept 1 or more subtypes,
	// and depending on the length of each combination we adjust the output.
	switch cctx.Module.Info.Type {
	case modulepb.Type_EXPLOIT:
		// Exploits can accept one transport.
	case modulepb.Type_PAYLOAD:
		// Payloads can accept one transport.
	case modulepb.Type_POST:
	case modulepb.Type_TRANSPORT:
		// Transports can accept one payload (for stagers)
	case modulepb.Type_UNDEFINED:
	}

	// Two consoles cannot always fit side-by-side in 13"
	if (105 < sWidth) && (sWidth < 140) {
		p += "{module_path}{fw}) "
		// if cctx.Module.HasSubmodule {
		// p += "{dim} && {fw}{submod_type}({r}{submod_name}{fw})"
		// }
	}
	// Here we can have two consoles side by side. (19")
	if (140 < sWidth) && (sWidth < 175) {
		p += "{module_path}{fw}) "
		// if cctx.Module.HasSubmodule {
		// p += "{dim} && {fw}{submod_type}({r}{submod_path}{fw})"
		// }
	}
	// Here we have one big console on a screen, no restrictions.
	if 175 < sWidth {
		p += "{module_path}{fw}) "
		// if cctx.Module.HasSubmodule {
		// p += "{dim} && {fw}{submod_type}({r}{submod_path}{fw})"
		// }
	}

	p += tui.RESET // Always at end

callbacks:
	// Callbacks
	for ok, cb := range m.Maincallbacks {
		p = strings.Replace(p, ok, cb(), 1)
	}
	for tok, color := range m.MainColors {
		p = strings.Replace(p, tok, color, -1)
	}

	width = getRealLength(p)

	return
}

// computeContext - Analyses the current state of various server indicators and displays them.
// Because it is the right-most part of the prompt, and that the screen might be small,
// we categorize default (assumed) screen sizes and we adapt the output consequently.
func (m *mainprompt) computeContext(mWidth, sWidth int) (p string, width int) {
	p += tui.RESET // Always

	p += "{dim}in {lv}{workspace} " // Always add the current workspace

	// Compute all values needed
	var items string

	// Half of my 13" laptop is around 115, and I have useless gaps of 5.
	// This means we have a small console space.
	if (105 < sWidth) && (sWidth < 120) {
		items = "{lv}{jobs}{fw}, {lc}{sessions}{fw}, {lv}{scans}{reset}"
	}
	// Two console cannot be side by side on a 13" laptop, but
	// we can on a 15" screen.
	if (120 < sWidth) && (sWidth < 140) {
		items = "{lv}{jobs}{fw}, {lc}{sessions}{fw}, {lv}{scans}{reset}"
	}
	// Half of my 19" monitor is around 165, again with gaps of 5.
	// Here we can have to consoles side by side.
	if (140 < sWidth) && (sWidth < 175) {
		items = "{fw}jobs({lv}{jobs}{fw}) sess({lc}{sessions}{fw}), scans({lv}{scans}{fw})"
	}

	// Here we have one big console on a screen, no restrictions.
	if 175 < sWidth {
		items = "{fw}jobs({lv}{jobs}{fw}) sess({lc}{sessions}{fw}), scans({lv}{scans}{fw})"
	}

	// Finally add the items string inside its container.
	p += fmt.Sprintf("{dim}--[{reset}%s{dim}]", items)

	p += tui.RESET // Always at end of prompt line

	// Callbacks
	for ok, cb := range m.Maincallbacks {
		p = strings.Replace(p, ok, cb(), 1)
	}
	for tok, color := range m.MainColors {
		p = strings.Replace(p, tok, color, -1)
	}

	width = getRealLength(p)

	return
}

// applyCallbacks - Replaces all callback variables within a prompt string, computes its real
// length by weeding out ASCII escape codes, and returns both.
func (m *mainprompt) applyCallbacks(in string) (p string, length int) {

	// First do items with the color-processed p string
	for ok, cb := range m.Maincallbacks {
		p = strings.Replace(in, ok, cb(), 1)
	}

	// First replace colors
	for tok, color := range m.MainColors {
		p = strings.Replace(p, tok, color, -1)
	}

	return
}

func (m *mainprompt) getPromptPad(total, base, module, context int) (pad string) {
	var padLength = total - base - module - context
	for i := 0; i < padLength; i++ {
		pad += " "
	}
	return
}

var (

	// MainCallbacks - All items needing calculation for the main prompt.
	MainCallbacks = map[string]func() string{
		// Local Working directory
		"{cwd}": func() string {
			cwd, _ := os.Getwd()
			return cwd
		},
		// Workspace
		"{workspace}": func() string {
			return "Finance_Department"
			// return cctx.Workspace.Name
		},
		// Local IP address
		"{local_ip}": func() string {
			addrs, _ := net.InterfaceAddrs()
			var ip string
			for _, addr := range addrs {
				network, ok := addr.(*net.IPNet)
				if ok && !network.IP.IsLoopback() && network.IP.To4() != nil {
					ip = network.IP.String()
				}
			}
			return ip
		},
		// Server IP
		"{server_ip}": func() string {
			return assets.ServerConfig.LHost
		},
		// Jobs and/or listeners
		"{jobs}": func() string {
			// return strconv.Itoa(cctx.Jobs)
			return "3"
		},
		// Nmap scans
		"{scans}": func() string {
			return "3" // Change this when available
		},
		// Sessions
		"{sessions}": func() string {
			return "12"
		},

		// Module type
		"{module_type}": func() string {
			// Safety checks, don't want problems
			if cctx.Module.Info == nil {
				return ""
			}
			switch cctx.Module.Info.Type {
			case modulepb.Type_EXPLOIT:
				return "exploit"
			case modulepb.Type_PAYLOAD:
				return "payload"
			case modulepb.Type_POST:
				return "post"
			case modulepb.Type_TRANSPORT:
				return "transport"
			case modulepb.Type_UNDEFINED:
				return "module"
			default:
				return "module"
			}
		},
		// Module path/name
		"{module_path}": func() string {
			modPath := cctx.Module.Info.Path // Alias

			// If there is no name and no path, this isn't normal, we return a blinking error
			if modPath == "" && cctx.Module.Info.Name == "" {
				return "\033[5m" + "module_path_name_not_found" + tui.RESET
			}

			// We check for duplicate module type in the path.
			path := strings.Split(modPath, "/")
			for _, modType := range modulepb.Type_name {
				if modType == path[0] {
					// If found we returned the truncated path
					return strings.Join(path[:len(path)-1], "/")
				}
			}
			// Else, return the full path
			return strings.Join(path, "/")
		},

		// Module name
		"{module_name}": func() string {
			// Safety checks, don't want problems
			if len(cctx.Module.Info.Path) == 0 && cctx.Module.Info.Name == "" {
				return "\033[5m" + "module_path_name_not_found" + tui.RESET
			}
			return cctx.Module.Info.Name
		},
		// Module subtypes, etc...
		"{submod_type}": func() string {
			return "post"
		},
		"{submod_path}": func() string {
			return "credentials/hash/"
		},
		"{submod_name}": func() string {
			return "GreatSubmodule"
		},
		// Stack data
		// Returns the number of modules on the user's stack
		"{stack_len}": func() string {
			return "3"
		},
	}

	// MainColorCallbacks - All colors and effects needed in the main menu
	MainColorCallbacks = map[string]string{
		// Base tui colors
		"{blink}": "\033[5m", // blinking
		"{bold}":  tui.BOLD,
		"{bdim}":  tui.DIM, // for Base Dim
		"{bfr}":   tui.RED, // for Base Fore Red
		"{g}":     tui.GREEN,
		"{b}":     tui.BLUE,
		"{y}":     tui.YELLOW,
		"{bfw}":   tui.FOREWHITE, // for Base Fore White.
		"{bdg}":   tui.BACKDARKGRAY,
		"{br}":    tui.BACKRED,
		"{bg}":    tui.BACKGREEN,
		"{by}":    tui.BACKYELLOW,
		"{blb}":   tui.BACKLIGHTBLUE,
		"{reset}": tui.RESET,
		// Custom colors
		"{ly}":   "\033[38;5;187m",
		"{lb}":   "\033[38;5;117m", // like VSCode var keyword
		"{db}":   "\033[38;5;24m",
		"{bddg}": "\033[48;5;237m",
		// Main theme colors
		"{fb}":  MainTheme[1],  // fore brown
		"{lc}":  MainTheme[3],  // light cream
		"{fw}":  MainTheme[7],  // forewhite, a bit darker
		"{dc}":  MainTheme[10], // dark cream
		"{r}":   MainTheme[11], // red, a bit darker
		"{lv}":  MainTheme[12], // light violet
		"{dv}":  MainTheme[14], // dark violet
		"{dim}": MainTheme[15], // Dim, a bit lighter
	}

	// MainTheme - 8 principal colors for the main menu context
	// The numbers used as keys mirrors the numbers used by the website terminal.sexy
	// when displaying the equivalent theme. This is thus an intermediary bookkeeping map.
	MainTheme = map[int]string{
		1:  util.Ctermfg130, // brown keyword
		3:  util.Ctermfg173, // Or 174 or 168/7     // cream type
		7:  util.Ctermfg250, // or 49/51            // forewhite
		10: util.Ctermfg172, // Or 173 or 210/9/8   // cream string type
		11: util.Ctermfg167, // or 160              // Red
		12: util.Ctermfg183, // Or 182, 175/6/7     // light violet keyword
		14: util.Ctermfg171, //                     // Violet
		15: util.Ctermfg240, // 41/2/3/4            // Dim
	}
)
