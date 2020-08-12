package module

import (
	"fmt"

	"github.com/sirupsen/logrus"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/maxlandon/wiregost/server/log"
)

// Module - The module interface is used by a stack to interact with a module type embedding
// a module.Module type that implements this interface. Some of the interface functions
// are reimplemented by these module subtypes, for adapting workflow to their task.
type Module interface {
	// Run - All modules have, at least from their module.Module base type, a Run() method.
	Run(cmd string, args []string) (result string, err error)
	// ToProtobuf - The module.Module can push info, opts & commands with it.
	ToProtobuf() (mod *pb.Module)
	// SetOption - The module has the ability to set its options.
	SetOption(opt *pb.Option) (err error)
	// AddModule - Some modules may be able to combine with other module types.
	// This method leaves them with how to handle their babies.
	AddModule(m Module) (ok bool, err error)
	// SetupLog - Called by module Stacks (server-side and stack-side), for
	SetupLog(remote bool, cli *clientpb.Client, rpc serverpb.EventsClient) (logger *logrus.Entry)
}

// Base - The base object inherited by all Base types in Wiregost. Any component in Wiregost that
// needs to be exposed to a User Interface (Console or RPC) can expose itself by embedding this type.
//
// As well, because the Base must provide scoped behavior when being run in a stack binary (which
// runs the author's module Run() func). This object therefore has methods that can be used by the
// binary stack for handling it correctly.
// Some of this object' methods will make use of gRPC client calls, to sync with their Stack peer.
type Base struct {
	Info     *pb.Info               // Base Information
	Commands map[string]*pb.Command // Module subcommands, can be added by embedders
	Opts     map[string]Option      // Module options, with utility methods on them

	// Logging: the module pushes messages to consoles through its logger
	// This should prove useful to identify precisely the source of events/
	// messages, because this logger will be derived by container types.
	// Ex (in Exploit subtype):
	// e.Log = e.Log.WithFields("exploit")
	Log *logrus.Entry

	// Owner - The user of this module instance. This is important as it will enable
	// us to further check permissions at various stages of a module execution.
	Client *clientpb.Client

	// We need methods for interacting either with :
	// the stack module equivalent when this module is in the server
	// the server module equivalent when this module is on the stack
	HasStackPeer bool // Does this module has an equivalent instance on a stack binary ?
}

// New - Instantiates a new module object, and populate it with the base information provided.
//
// This function is called twice: on server and on stack binary. The server should call it
// for the according entry in the result of ToProtobuf(), called by the stack binary.
func New(meta *pb.Info) (m *Base) {

	m = &Base{
		Info:     meta,                     // Base information
		Commands: map[string]*pb.Command{}, // Empty commands
		Opts:     map[string]Option{},      // Module options, with utility methods on them
		// Client:   client,                         // The client/user having loaded the module.
		// We don't put this here, it bothers us for all module subtypes instantiations.
	}

	return
}

// AddCommand - Adds a command to the module, which is used by the user like 'run check'.
// @ name           - The name of the command
// @ description    - A description for this command
// @ hasPayload     - True if this command triggers a remote exploit, requiring a payload.
//
// This function is called twice: on server and on stack binary. The server should call it
// for every entry in the result of ToProtobuf(), called by the stack binary.
func (m *Base) AddCommand(name, description string, hasPayload bool) {
	// Skip command if name is empty. Do not raise errors

	// Trim spaces, don't want bad surprises because of typos.
}

// ToProtobuf - Returns the module with its information, options and commands. Used to share
// state back and forth between consoles, the server and stacks. All categories are separately
// parsed, because this amazing module structure dictates (at least) that.
func (m *Base) ToProtobuf() (mod *pb.Module) {

	// Base information
	mod = &pb.Module{
		Info:     m.Info,
		Commands: map[string]*pb.Command{},
		Options:  map[string]*pb.Option{},
	}

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
// The function is exported: it is only called by module Stacks (server-side and stack-side), for
// setting up the module's respective logging infrastructures.
func (m *Base) SetupLog(remote bool, cli *clientpb.Client, rpc serverpb.EventsClient) (logger *logrus.Entry) {

	m.Client = cli // Can do that here: this func is called quick after New()
	m.Log = log.ModuleLogger(remote, cli, rpc)

	return
}

// Asset - Find the path of an asset in the module source directory.
// It is exported, because authors will need to access/use non-Go files.
func (m *Base) Asset(path string) (filePath string, err error) {
	return
}

// AskUserConfirm - Some actions performed by the module might require user permission,
// because it has some context awareness needs that a computer is obviously unable to have.
func (m *Base) AskUserConfirm() (ok bool) {
	return
}

// PreRunChecks - All checks for commands, options, etc. are done in this function.
// Does not concern various compatibility and perms checks done by drivers, it is
// merely utility checking code to avoid useless troubles.
func (m *Base) PreRunChecks(cmd string) (err error) {

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
func (m *Base) CheckRequiredOptions() (err error) {
	return
}

// CheckCommand - Verifies the command run by the user (like 'run exploit' or 'run check') is valid.
func (m *Base) CheckCommand(command string) (err error) {
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
func (m *Base) run() (err error) {
	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Base) Cleanup() (err error) {

	return
}
