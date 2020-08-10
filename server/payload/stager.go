package payload

import (
	modpb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/module"
	"github.com/maxlandon/wiregost/server/module/stack"
	"github.com/maxlandon/wiregost/server/transport"
)

// Stager - A module used to blend in a Payload module into a Transport module.
// Therefore, the Stager module has access to all the methods of both modules,
// which he can invoke along a payload staging process, over the network.
type Stager struct {
	// Base module. Makes this Stager a valid module in Wiregost, with full access to UI.
	*module.Module

	// The Transport module is used to convey necessary information to either the server
	// or implants that may need to start listeners, or setup for this. It also fournishes
	// a set of supplementatry fields for working with the network.
	Transport *transport.Module

	// The payload module is here to provide all methods necessary to generate, setup and
	// use a stager, no matter its architecture, OS, etc.
	Payload *Module
}

// NewStager - Instantiates a new Stager module. Called by Exploit modules.
func NewStager(meta *modpb.Info) (m *Stager) {

	m = &Stager{
		module.New(meta), // Populate base module
		nil,              // Don't know which one we're using yet
		nil,              // Don't know which one we're using yet.
	}

	// A stager, although being in the Payload package, is at first a Transport type, because most
	// of the "live" job is to initiate a connection over the network and interact with it.
	m.Info.Type = modpb.Type_TRANSPORT

	// Add specific fields to the Stager logger. Overwrites "module":"module" key/val pair.
	m.Log = m.Log.WithField("module", "stager")

	return
}

// AddModule - Implements the stack.Module interface. This Stager checks various elements
// and then rejects the provided module. It is even possible that this Module might be
// checked again various interfaces, for instance if it is a Payload subtype.
func (m *Stager) AddModule(mod stack.Module) (ok bool, err error) {

	// Check if module is either transport or payload

	// Check for compatibility

	return
}

// SetOption - Implements the stack.Module interface. The stager must perform various
// checks because the provided option might be in payload subtypes.
func (m *Stager) SetOption(opt *modpb.Option) (err error) {

	if option, found := m.Opts[opt.Name]; found {
		err = option.Set(opt.Value)
	}

	// We came here, there might be an error but we must
	// check for other categories first: payload or transport

	return
}

// Run - Execute the main function of this Stager.
func (m *Stager) Run(cmd string, args []string) (result string, err error) {

	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Stager) Cleanup() (err error) {

	return
}
