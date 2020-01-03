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
			s.CurrentWorkspaceId = notif.WorkspaceId
		case "delete":
			fmt.Printf("%s[!]%s Workspace %s deleted from another client shell. Falling back to default.",
				tui.BOLD, tui.RESET, s.currentWorkspace)
			s.currentWorkspace = "default"
			s.CurrentWorkspaceId = notif.FallbackWorkspaceId
			// Refresh shell
			s.Refresh()
		}
	}

}
