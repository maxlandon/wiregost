package console

import (
	ansi "github.com/acarl005/stripansi"
	"github.com/evilsocket/islazy/tui"

	"github.com/maxlandon/wiregost/client/context"
	"github.com/maxlandon/wiregost/client/util"
)

// cctx - An alias to the console state/context
var cctx = context.Context

// promptbis - The object in charge of computing prompts, refreshing and printing them.
type promptbis struct {
	Main  *mainprompt  // main menu
	Ghost *ghostprompt // ghost session menu
}

func (c *console) InitPrompt() {

	// Various values related to readline
	c.Shell.Multiline = true   // spaceship-like (two-line) prompt
	c.Shell.ShowVimMode = true // with Vim status

	PromptBis = &promptbis{
		Main: &mainprompt{
			Maincallbacks: MainCallbacks,
			MainColors:    MainColorCallbacks,
		},
		Ghost: &ghostprompt{},
	}

	c.Shell.SetPrompt(PromptBis.Render())
}

// Render - The ghost menu prompt outputs its computed content
func (m *promptbis) Render() (prompt string) {

	switch cctx.Menu {
	case context.MainMenu, context.ModuleMenu:
		prompt = m.Main.Render()
	case context.GhostMenu:
		prompt = m.Ghost.Render()
	}
	return
}

// ghostprompt - A prompt used when the user interacts with a ghost session.
type ghostprompt struct {
	Base             string // The left-side: implant name/UUID, user, currDir, etc.
	CustomGhost      string // A user-provided custom base
	ContextPrompt    string // The right-most side of the prompt (channels, pivots, workspaces...)
	StatusDisconnect string // The ghost implant is either lost or down, so we have restrictions.
	OptionPrompt     string // An option is set (with error or not), action results are stored here.

	// callbacks map[string]func() string // All values are automatically calculated
}

// Render - The ghost menu prompt outputs its computed content
func (m *ghostprompt) Render() (prompt string) {
	return
}

var (
	GhostCallbacks = map[string]string{
		// Session type
		// Ghost implant name/UUID
		// Impersonated User/UserID/GID
		// Hostname
		// implant (or implant channel) working directory
		// Current transport sheme:addres:port
		// Target OS/Arch, etc, in one nice string
	}

	// GhostColorCallbacks - All colors and effects needed in the ghost session menu
	GhostColorCallbacks = map[string]string{
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
		"{ly}": "\033[38;5;187m",
		// "{lb}":   "\033[38;5;117m", // like VSCode var keyword
		"{bddg}": "\033[48;5;237m",
		// Main theme colors
		"{fb}":  GhostTheme[1],  // fore brown
		"{lv}":  GhostTheme[2],  // light violet
		"{cg}":  GhostTheme[3],  // cream green
		"{db}":  GhostTheme[4],  // dim blue
		"{lbg}": GhostTheme[6],  // light blue gray
		"{fw}":  GhostTheme[7],  // forewhite, a bit darker
		"{dim}": GhostTheme[8],  // Dim, a bit lighter
		"{sb}":  GhostTheme[9],  // sky blue
		"{tg}":  GhostTheme[10], // turquoise green
		"{lb}":  GhostTheme[12], // blue
		"{dv}":  GhostTheme[13], // dim violet
	}

	// GhostTheme - 20 colors for the ghost implant context (terminal.sexy kasugano theme)
	GhostTheme = map[int]string{
		1:  util.Ctermfg130, // brown keyword  -
		2:  util.Ctermfg72,  // alternative skyblue with 9  -
		3:  util.Ctermfg158, //  or 159/194/195 // white-green type
		4:  util.Ctermfg25,  // or 24               // Dim blue type
		6:  util.Ctermfg103, // or 110              // blue-grey
		7:  util.Ctermfg250, // or 49/51            // forewhite
		8:  util.Ctermfg237, //  or 6/8/9        // Dim
		9:  util.Ctermfg111, // or 75/68/     // sky blue
		10: util.Ctermfg79,  // 72/86        // turquoise string type
		12: util.Ctermfg32,  // or 25/27/33    // blue keyword
		13: util.Ctermfg61,  // or 62/104   // violet dim
	}
)

// getRealLength - Some strings will have ANSI escape codes, which might be wrongly
// interpreted as legitimate parts of the strings. This will bother if some prompt
// components depend on other's length, so we always pass the string in this for
// getting its real-printed length.
func getRealLength(s string) (l int) {
	return len(ansi.Strip(s))
}
