package module

// Transport - Not the equivalent of a handler in Metasploit... but much more.
// Transports are a module in Wiregost because they might be used in several contexts:
// - With payload modules, for adding transports configuration to implants that will be compiled.
// - With running implants, to which we can still add transports (removing them and cycling
//   through them is done via command line)
//
// It is the only type of module that has no driver, because
// the working process of handlers is more straightforward.
type Transport struct {
	// Needs a profobuf object

	// Needs to keep track of a connection that should be the same
	// whether it is a physical one or not: that transport might
	// reach far into Wiregost's routing system if its a Bind type transport.
	// This means we might have to store (or make use of) route chains at some point here.

	// Needs some logging

	// Needs some working hours
}

// NewTransport - Instantiate a Transport with basic information.
func NewTransport() (t *Transport) {
	t = &Transport{}
	// Needs an ID/UUID
	return
}

// Init - Initializes a Transport and performs basic setup
func (t *Transport) Init() (err error) {

	// Setup logging

	// Setup working hours

	return
}

// Start - The transport starts with all provided settings.
func (t *Transport) Start() (err error) {
	return
}

// startRemote - The transport is a listener (reverse), started on a ghost implant.
// Therefore, instead of starting the listener here, send requests to the concerned implant,
// and handle all details necessary with it. At the end, we don't know anymore (sort of) that
// the Transport had to do this.
func (t *Transport) startRemote() (err error) {
	return
}
