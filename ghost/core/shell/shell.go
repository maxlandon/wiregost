package shell

import (
	"io"

	"github.com/go-cmd/cmd"
)

// Shell - Represents a shell instance on the target system.
// Many shells can be opened and used concurrently from the same implant.
type Shell struct {
	ID      uint64
	Command *cmd.Cmd // We use the go-cmd package for better/easier process management
	Stdout  io.ReadCloser
	Stdin   io.WriteCloser
}
