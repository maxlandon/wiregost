package clients

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	"github.com/maxlandon/wiregost/server/assets"
)

// GetHistory - A console requests a command history line.
func (c *connectionServer) GetHistory(in context.Context, req *pb.HistoryRequest) (res *pb.History, err error) {

	fmt.Println(req)
	cli := Consoles.GetClient(req.Client.ID)
	if cli == nil {
		return nil, errors.New("server could not find your client when requesting history")
	}

	// Get path for user
	path := assets.GetUserHistoryDir(cli.User.ID, cli.User.Name)

	var data []byte // Contents of the file

	// Find file data, cut it and process it
	if req.AllConsoles {
		filename := filepath.Join(path, fmt.Sprintf("aggregated_user_history"))
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
	} else {
		filename := filepath.Join(path, fmt.Sprintf("console_%s", cli.ID))
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
	}
	lines := strings.Split(string(data), "\n") // For all returns we have a new line

	fmt.Println(lines)
	fmt.Println(lines[req.Index])
	return &pb.History{Line: lines[req.Index], HistLength: int32(len(lines))}, nil
}

// AddToHistory - A client has sent a new command input line to be saved.
func (c *connectionServer) AddToHistory(in context.Context, req *pb.AddCmdHistoryRequest) (res *pb.AddCmdHistory, err error) {

	fmt.Println("adding to hist")
	fmt.Println(req)
	cli := Consoles.GetClient(req.Client.ID)
	if cli == nil {
		return nil, errors.New("server could not find your client when requesting history")
	}

	// Filter various useless commands
	if stringInSlice(strings.TrimSpace(req.Line), uselessCmds) {
		return &pb.AddCmdHistory{Doublon: true}, nil
	}

	// Get path for user
	path := assets.GetUserHistoryDir(cli.User.ID, cli.User.Name)

	// Write to client console file
	filename := filepath.Join(path, fmt.Sprintf("console_%s", cli.ID))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.New("server could not find your client when requesting history: " + err.Error())
	}
	if _, err = f.WriteString(req.Line + "\n"); err != nil {
		return nil, errors.New("server could not find your client when requesting history: " + err.Error())
	}
	f.Close()

	// Write to aggregated_user_history file
	filename = filepath.Join(path, fmt.Sprintf("aggregated_user_history"))
	f, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.New("an error happened when writing to your history file")
	}
	if _, err = f.WriteString(req.Line + "\n"); err != nil {
		return nil, errors.New("an error happened when writing to your history file")
	}
	f.Close()

	return &pb.AddCmdHistory{}, nil
}

// A list of commands that are useless to save if they are strictly as short as in the list
var uselessCmds = []string{
	"exit",
	"cd",
	"ls",
	"ls",
	"cat",
	"pwd",
	"use",
	"clear",
	"back",
	"pop",
	"push",
	"stack",
	"config",
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
