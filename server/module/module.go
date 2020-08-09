package module

import (
	"fmt"

	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/sirupsen/logrus"
)

// Module - The base object inherited by all Module types in Wiregost. Any component in Wiregost that
// needs to be exposed to a User Interface (Console or RPC) can expose itself by embedding this type.
//
// As well, because the Module must provide scoped behavior when being run in a stack binary (which
// runs the author's module Run() func). This object therefore has methods that can be used by the
// binary stack for handling it correctly.
// Some of this object' methods will make use of gRPC client calls, to sync with their Stack peer.
type Module struct {
	Info     *modulepb.Info               // Base Information
	Commands map[string]*modulepb.Command // Module subcommands, can be added by embedders
	Opts     map[string]Option            // Module options, with utility methods on them

	// Logging: the module pushes messages to consoles through its logger
	// This should prove useful to identify precisely the source of events/
	// messages, because this logger will be derived by container types.
	// Ex (in Exploit subtype):
	// e.Log = e.Log.WithFields("exploit")
	Log *logrus.Entry

	// Owner - The user of this module instance. This is important as it will enable
	// us to further check permissions at various stages of a module execution.
	Owner *dbpb.User

	// We need methods for interacting either with :
	// the stack module equivalent when this module is in the server
	// the server module equivalent when this module is on the stack
	HasStackPeer bool // Does this module has an equivalent instance on a stack binary ?
}

// New - Instantiates a new module object, and populate it with the base information provided.
//
// This function is called twice: on server and on stack binary. The server should call it
// for the according entry in the result of ToProtobuf(), called by the stack binary.
func New(meta *modulepb.Info) (m *Module) {
	return &Module{Info: meta}
}

// information - Registers the module's metadata information
func (m *Module) information(meta *modulepb.Info) (err error) {
	m.Info = meta
	return
}

// AddCommand - Adds a command to the module, which is used by the user like 'run check'.
// @ name           - The name of the command
// @ description    - A description for this command
// @ hasPayload     - True if this command triggers a remote exploit, requiring a payload.
//
// This function is called twice: on server and on stack binary. The server should call it
// for every entry in the result of ToProtobuf(), called by the stack binary.
func (m *Module) AddCommand(name, description string, hasPayload bool) {
	// Skip command if name is empty. Do not raise errors

	// Trim spaces, don't want bad surprises because of typos.
}

// ToProtobuf - Returns the module with its information, options and commands. Used to share
// state back and forth between consoles, the server and stacks. All categories are separately
// parsed, because this amazing module structure dictates (at least) that.
func (m *Module) ToProtobuf() (mod *modulepb.Module) {

	// Base information
	mod = &modulepb.Module{Info: m.Info}

	// Commands
	for name, cmd := range m.Commands {
		mod.Commands[name] = cmd
	}

	// Options
	for name, opt := range m.Opts {
		mod.Options[name] = opt.toProtobuf()
	}

	return
}

// SetupLog - An exported function that will be called by the Stack module object. This logger will have
// hooks linked to a gRPC service, for sending back status to user consoles and logging files.
func (m *Module) SetupLog() (logger *logrus.Entry) {

	// When this function is called, we should have determined if the module has a remote side:
	// If yes, we add a special hook for logging over gRPC, or we make a completely different
	// logrus.Logger in a specialized function setupLogRemote() (logger *logrus.Logger)
	// to which we then add fields.

	// Add fields for userID or ClientID

	// Add WithFields("module") if needed

	return
}

// Asset - Find the path of an asset in the module source directory.
// It is exported, because authors will need to access/use non-Go files.
func (m *Module) Asset(path string) (filePath string, err error) {
	return
}

// AskUserConfirm - Some actions performed by the module might require user permission,
// because it has some context awareness needs that a computer is obviously unable to have.
func (m *Module) AskUserConfirm() (ok bool) {
	return
}

// PreRunChecks - All checks for commands, options, etc. are done in this function.
// Does not concern various compatibility and perms checks done by drivers, it is
// merely utility checking code to avoid useless troubles.
func (m *Module) PreRunChecks(cmd string) (err error) {

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
func (m *Module) CheckRequiredOptions() (err error) {
	return
}

// CheckCommand - Verifies the command run by the user (like 'run exploit' or 'run check') is valid.
func (m *Module) CheckCommand(command string) (err error) {
	// If we don't have commands it means there is no need for them, so no error
	if len(m.Commands) == 0 || m.Commands == nil {
		return nil
	}

	// Else check for it
	for _, cmd := range m.Commands {
		if cmd.Name == command {
			return nil
		}
	}
	return fmt.Errorf("invalid command: %s", command)
}

// Run method (called by driver) using the provided module subcommand (run 'check_vuln')o
// This will trigger the appropriate function in the Stack module itself, while preserving
// the ability to use drivers such as ExploitDriver for managing the full process of an exploit.
func (m *Module) run() (err error) {
	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Module) Cleanup() (err error) {

	return
}
