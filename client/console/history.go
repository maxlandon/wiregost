package console

import (
	"context"

	"github.com/maxlandon/wiregost/client/connection"
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
)

var (
	// ClientHist - Client console history
	ClientHist = &ClientHistory{LinesSinceStart: 1}
	// UserHist - User history
	UserHist = &UserHistory{LinesSinceStart: 1}
)

// This file manages all command history flux for this console. The user can request
// 2 different lists of commands: the history for this console only (identified by its
// unique ID) with Ctrl-r, or the history for all the user's consoles, with Ctrl-R.

// ClientHistory - Writes and queries only the Client's history
type ClientHistory struct {
	LinesSinceStart int // Keeps count of line since session
}

// Write - Sends the last command to the server for saving
func (h *ClientHistory) Write(s string) (int, error) {

	res, err := connection.ConnectionRPC.AddToHistory(context.Background(),
		&clientpb.AddCmdHistoryRequest{Line: s, Client: cctx.Client})
	if err != nil {
		return 0, err
	}

	if !res.Doublon {
		h.LinesSinceStart++
	}
	return h.LinesSinceStart, nil
}

// GetLine returns a line from history
func (h *ClientHistory) GetLine(i int) (string, error) {

	res, err := connection.ConnectionRPC.GetHistory(context.Background(),
		&clientpb.HistoryRequest{
			AllConsoles: false,
			Index:       int32(i),
			Client:      cctx.Client,
		})
	if err != nil {
		return "", err
	}
	h.LinesSinceStart = int(res.HistLength)

	return res.Line, nil
}

// Len returns the number of lines in history
func (h *ClientHistory) Len() int {
	return h.LinesSinceStart
}

// Dump returns the entire history
func (h *ClientHistory) Dump() interface{} {
	return nil
}

// UserHistory - Only in charge of queries for the User's history
type UserHistory struct {
	LinesSinceStart int // Keeps count of line since session
}

func (h *UserHistory) Write(s string) (int, error) {
	h.LinesSinceStart++
	return h.LinesSinceStart, nil
}

// GetLine returns a line from history
func (h *UserHistory) GetLine(i int) (string, error) {

	res, err := connection.ConnectionRPC.GetHistory(context.Background(),
		&clientpb.HistoryRequest{
			AllConsoles: true,
			Index:       int32(i),
			Client:      cctx.Client,
		})
	if err != nil {
		return "", err
	}
	h.LinesSinceStart = int(res.HistLength)

	return res.Line, nil
}

// Len returns the number of lines in history
func (h *UserHistory) Len() int {
	return h.LinesSinceStart
}

// Dump returns the entire history
func (h *UserHistory) Dump() interface{} {
	return nil
}
