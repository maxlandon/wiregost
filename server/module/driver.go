package module

// Driver - A "driver of drivers": running it will run all modules gathered between an ExploitDriver
// and a PayloadDriver. One of them may be empty, because payloads can be compiled and transports
// launched without an Exploit, and an Exploit might not necessarily require a Payload for some actions.
type Driver struct {
	// Sub-drivers
	edriver *ExploitDriver

	pdriver *PayloadDriver
}

// NewDriver - Instantiates a new Driver each time a new module is loaded
func NewDriver() (d *Driver) {
	d = &Driver{
		edriver: NewExploitDriver(),
		pdriver: NewPayloadDriver(),
	}

	return
}

// Run module methods
// We could have a general Driver, handling and synchronizing both ExploitDrivers and PayloadDrivers.
// This would clean a bit the code of the ExploitDriver, things we could mutualize, etc.

// It would be a bit the inverse equivalent of the msf/lib/base/simple/exploit.rb file, in which a base Exploit
// module creates its own driver, and synchronises it then with the Exploit's payload:
// The role of the Manager would be to a subdriver to both, the Exploit would not create and handle its own driver.
func (d *Driver) Run() (err error) {

	// Depending on the type of module present on the driver, we have to check and branch into new funcs.

	// If Exploit is not nil in edriver, setup the edriver

	// If PayloadDriver is ready, and requested exploit command requires a Payload.

	// Validate Exploit settings, if needed

	// Set keep handler to PayloadDriver if Exploit has multiple targets

	// Validate Payload settings

	// Here, Metasploit would do things with target index. Depending on how
	// we plan to handle them, we might need something here.

	// Here, Metasploit wires the modules (Exploit and Payload) to the console UI.

	// If exploit is set:
	// RUN THE EXPLOIT
	d.runExploit()

	// We should get back a Job_ID, especially if we ran the exploit command in background

	// Here Metasploit handles potential errors.

	return
}

// runExploit - This function takes care of running an Exploit module: here, we set up and synchronize
// two different drivers: an ExploitDriver and a PayloadDriver. The ExploitDriver should have an Exploit loaded,
// while the PayloadDriver can be empty, if the Exploit function triggered does not require a network handler.
func (d *Driver) runExploit() (result string, err error) {

	// Clone the Exploit and Payload if needed, so that no changes affect them during the process.

	// Instantiate an ExploitDriver with an Exploit module instance in it.

	// Instantiate a PayloadDriver

	return
}
