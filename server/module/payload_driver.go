package module

// PayloadDriver - A driver for payload construction, setup, compilation and staging.
// An instance is created each time a user loads up a payload module.
// Drivers are, in some sort, templates that also enhance the way information is accessible to console
// users, as well as the way they can modify it, maintaining granularity as well as good context scoping.
//
// It plays a role similar to ExploitDriver, thereby offering multiple possibilities
// for configuring per OS, per transport, and for potentially adding security modules.
//
// It is somehow a bit different from an ExploitDriver though:
// - While the ExploitDriver can only have one Exploit module, the PayloadDriver can have one Payload,
type PayloadDriver struct {
	// Manages a Payload module
	Payload *Payload

	// Transports. They are handled and stored by the driver because they will be used by it
	// in different cases: for compiling implants with these transports information, as well
	// starting handlers based on them, which will need the routing system to work correctly.
	// Might need some route chains as well.
	Transports []*Transport

	// Manages some handlers (they pull their configuration out of the implant's transports)

	// The payload driver might have to access Wiregost's routing system. Maybe not.

	ReadyForExploit bool // Used to signal at least one payload or Transport is ready,
	// so a remote exploit can run and will find a listener/binder back soon.
}

// NewPayloadDriver - Instantiate the driver with some core information and components.
func NewPayloadDriver() (p *PayloadDriver) {
	p = &PayloadDriver{}
	return
}
