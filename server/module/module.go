package module

import (
	"fmt"

	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/sirupsen/logrus"
)

// module - The base object inherited by all module types in Wiregost. This object is, at its core,
// a client object that synchronises state with an equivalent Stack module, and can perform some
// base operations on it: getting information, setting options, running checks, running it.
// Some of this object's methods will make use of gRPC client calls, to sync with their Stack peer.
// The type is unexported, but its methods are: The user needs to use its behavior, not its state.
type module struct {
	info *modulepb.Info // Info (protobuf object)

	// We have commands, separated from info and options
	commands map[string]*modulepb.Command

	// We don't need options, because they are used only by the Stack module peer.
	Opts map[string]*Option

	// Logging: the module pushes messages to consoles through its logger
	Log *logrus.Entry

	// We need methods for interacting either with :
	// the stack module equivalent when this module is in the server
	// the server module equivalent when this module is on the stack
	// Can be done at the level up, by the manager handling these modules
}

// information - Registers the module's metadata information
func (m *module) information(meta *modulepb.Info) (err error) {
	m.info = meta
	return
}

// AddCommand - Adds a command to the module, which is used by the user like 'run check'.
// - name           - The name of the command
// - description    - A description for this command
// - hasPayload     - True if this command triggers a remote exploit, requiring a payload.
func (m *module) AddCommand(name, description string, hasPayload bool) {
	// Skip command if name is empty. Do not raise errors

	// Trim spaces, don't want bad surprises because of typos.
}

// ToProtobuf - Returns the module with its information, options and commands. Used to share
// state back and forth between consoles, the server and stacks. All categories are separately
// parsed, because this amazing module structure dictates (at least) that.
func (m *module) ToProtobuf() (mod *modulepb.Module) {

	// Base information
	mod = &modulepb.Module{Info: m.info}

	// Commands
	for name, cmd := range m.commands {
		mod.Commands[name] = cmd
	}

	// Options
	for name, opt := range m.Opts {
		mod.Options[name] = opt.info
	}

	return
}

// SetupLog - An exported function that will be called by the Stack module object. This logger will have
// hooks linked to a gRPC service, for sending back status to user consoles and logging files.
func (m *module) SetupLog() (err error) {
	return
}

// Asset - Find the path of an asset in the module source directory.
// It is exported, because authors will need to access/use non-Go files.
func (m *module) Asset(path string) (filePath string, err error) {
	return
}

// AskUserConfirm - Some actions performed by the module might require user permission,
// because it has some context awareness needs that a computer is obviously unable to have.
func (m *module) AskUserConfirm() (ok bool) {
	return
}

// PreRunChecks - All checks for commands, options, etc. are done in this function.
// Does not concern various compatibility and perms checks done by drivers, it is
// merely utility checking code to avoid useless troubles.
func (m *module) PreRunChecks(cmd string) (err error) {

	err = m.CheckCommand(cmd)
	if err != nil {
		return err
	}

	err = m.CheckRequiredOptions()
	if err != nil {
		return err
	}

	return
}

// CheckRequiredOptions - Checks that all required options have a value. The advantage of this
// function is that it can check for options that have been registered by module subtypes as
// well, as they all lie in the same Opts object.
func (m *module) CheckRequiredOptions() (err error) {
	return
}

// CheckCommand - Verifies the command run by the user (like 'run exploit' or 'run check') is valid.
func (m *module) CheckCommand(command string) (err error) {
	// If we don't have commands it means there is no need for them, so no error
	if len(m.commands) == 0 || m.commands == nil {
		return nil
	}

	// Else check for it
	for _, cmd := range m.commands {
		if cmd.Name == command {
			return nil
		}
	}
	return fmt.Errorf("invalid command: %s", command)
}

// Run method (called by driver) using the provided module subcommand (run 'check_vuln')o
// This will trigger the appropriate function in the Stack module itself, while preserving
// the ability to use drivers such as ExploitDriver for managing the full process of an exploit.
