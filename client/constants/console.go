package constants

const (
	// Console menu contexts -----------------------------------------------

	// MainContext - Available only in main menu
	MainContext = "main"
	// ModuleContext - Available only when a module is loaded
	ModuleContext = "module"
	// GhostContext - Available only when interacting with a ghost implant
	GhostContext = "ghost"

	// MetadataKey - Used to reference the Data struct contained in the context
	MetadataKey = "wiregost"
)
