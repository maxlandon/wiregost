package payload

import (
	modpb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/module"
	"github.com/maxlandon/wiregost/server/module/stack"
	"github.com/maxlandon/wiregost/server/transport"
)

// Module - A payload module in Wiregost is in charge of Payload setup,
// compilation and sharing for various other modules, such as a Stagers.
//
// Various subtypes might embed this Module, so that they can benefit from
// various methods while performing more specialized Payload actions: for
// instance, a Ghost Payload is a specialized type which takes care of
// setting up a Ghost implant.
type Module struct {
	// Base module. Makes this Stager a valid module in Wiregost, with full access to UI.
	*module.Base

	// Transport - The current Transport module loaded by a user for this payload. A user
	// may only load one of these at a time, but it may act on it in different ways.
	transport transport.Transport

	//Transports - One of the ways to use a Transport module with a Payload module is to
	// add various of its elements to the Payload configuration, before it is compiled.
	Transports []transport.Transport
}

// New - Instantiates a new Payload module. Called by console users and Stager modules.
func New(meta *modpb.Info) (m *Module) {

	m = &Module{
		module.New(meta),        // Base module
		nil,                     // No Transport selected yet
		[]transport.Transport{}, // No compiled transports
	}

	m.Info.Type = modpb.Type_PAYLOAD

	// Add specific fields to the Payload logger. Overwrites "module":"module" key/val pair.
	m.Log = m.Log.WithField("module", "payload")

	return
}

// AddModule - Implements the stack.Module interface. This Payload only accept a Transport.
func (m *Module) AddModule(mod stack.Module) (ok bool, err error) {

	if mod.ToProtobuf().Info.Type != modpb.Type_TRANSPORT {
		return false, nil
	}

	// Perform various checks for compatibility.

	// Then add the Transport

	return
}

// SetOption - Implements the stack.Module interface. The stager must perform various
// checks because the provided option might be in payload subtypes.
func (m *Module) SetOption(opt *modpb.Option) (err error) {

	if option, found := m.Opts[opt.Name]; found {
		err = option.Set(opt.Value)
	}

	// We came here, there might be an error but we must
	// check for other categories first: transport.

	return
}

// Run - Execute the main function of this Payload module.
func (m *Module) Run(cmd string, args []string) (result string, err error) {

	return
}

// Cleanup - Clean any state needed for this module. This function is here more to remind
// all types embedding this module that they may override it, as a good practice of cleaning.
func (m *Module) Cleanup() (err error) {

	return
}
