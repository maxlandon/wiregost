package post

import (
	pb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/module"
)

// Init - Setup a module instance, and return it to the server for loading. This func is called at startup.
func Init() (m *Module, err error) {

	// Instantiate module with parent types and functions. Do not touch this line.
	m = &Module{module.NewPost()}

	// Module Information
	m.Init(&pb.Module{
		// Base information
		Name:        "[OS][category][domain][subject]",
		Authors:     []string{"Author 1", "Author 2"},
		Credits:     []string{"Contributor 1", "Contributor 2"},
		Description: "This module is an example post module. Fill the details before implementing it.",
		Priviledged: false,
		// OS/Platform/language
		Lang: "Go",
		// Other
		Notes: "You can add some notes to this module, whether it be additional references to check or details to care for.",
	})

	// Commands
	m.AddCommand("add_user", "This command performs a user-related post-exploitation activity.", false)
	m.AddCommand("migrate_dll", "This function triggers a vulnerability allowing to migrate the implant.", true)

	// Options
	m.AddOption(
		"RHost",                           // The option name (cannot be empty)
		"Base",                            // The option category name (can be empty), used for pretty printing on consoles.
		"192.168.1.1/24",                  // If non-empty, this will serve a default value. This field is used for the option value.
		"The remote host address to dial", // Option description
		true,                              //Is this option required ? If true, the Value field cannot be empty.
	)

	return
}
