package module

// Payload - A pseudo-module dedicated to payload generation and staging.
type Payload struct {
	// Includes a ghost implant profile (format, arch, compilation, transport-agnostic)
	// A ghost implant as bytes, for stagers
	// Some transports
	// Some security (limits, permissions, certificates, etc)
	// Some working hours
	// Some OS specifics
}

// Some of these functions will be hidden and called internally only

// Generate (with optional saving of the implant as file, or as bytes).
