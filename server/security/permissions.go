package security

import (
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
)

// CheckCorePermissions - Verifies that implant core functionality is usable by the user requesting it.
func CheckCorePermissions(ghost *ghostpb.Ghost, user *dbpb.User) (ok bool, err error) {
	return
}
