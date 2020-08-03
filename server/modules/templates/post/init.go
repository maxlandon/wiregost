package main // Rename this with module (directory) name

import (
	"github.com/maxlandon/wiregost/server/modules/post"
)

// Init - Setup the module, and return it to the server for loading. This func is called at startup.
func Init() (m *Module, err error) {

	// Base setup
	m = &Module{post.NewPost()} // Instantiate module with all parent types and functionalities
	err = m.RegisterOptions()   // Register all options added by module author below
	i := m.Info                 // Shorthand for less verbosity when filling information below

	// Module Information
	i.Name = "MyModuleName"
	i.Path = "Must check if we still need this, or if there is a way to automate it"
	i.Authors = []string{"Author 1", "Author 2"}
	i.Credits = []string{"Contributor 1", "Contributor 2"}
	i.Description = "This module is an example post module. Fill the details before implementing it."
	i.Lang = "Go"
	i.Priviledged = false

	// Module targets
	i.Targets = []string{"windows", "linux", "darwin"}

	// Others

	return
}

// RegisterOptions - Use this function to add options to your module.
// The function is then called as part of the Init() process above.
func (m *Module) RegisterOptions() (err error) {

	// Option Registration Example
	m.NewOption(
		"RHost",                           // The option name (cannot be empty)
		"Base",                            // The option category name (can be empty), used for pretty printing on consoles.
		"192.168.1.1/24",                  // If non-empty, this will serve a default value. This field is used for the option value.
		"The remote host address to dial", // Option description
		true,                              //Is this option required ? If true, the Value field cannot be empty.
	)
	return
}
