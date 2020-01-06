package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
)

func (s *Session) handleNotifications(notif messages.Notification) {
	// Determine type of notification and unwrap
	switch notif.Type {
	case "workspace":
		switch notif.Action {
		case "switch":
			s.currentWorkspace = notif.Workspace
			s.CurrentWorkspaceID = notif.WorkspaceID
		case "delete":
			fmt.Printf("%s[!]%s Workspace %s deleted from another client shell. Falling back to default.",
				tui.BOLD, tui.RESET, s.currentWorkspace)
			s.currentWorkspace = "default"
			s.CurrentWorkspaceID = notif.FallbackWorkspaceID
			// Refresh shell
			s.refresh()
		}
	case "module":
		switch notif.Action {
		case "pop":
			if s.currentModule == notif.PoppedModule {
				s.currentModule = notif.FallbackModule
				if s.currentModule == "" {
					s.Shell.Config.AutoComplete = s.getCompleter("main")
				}
				// Refresh shell
				s.refresh()
			}
		}
	}

}
