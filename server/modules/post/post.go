// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package post

import (
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
	"github.com/maxlandon/wiregost/server/ghosts"
	"github.com/maxlandon/wiregost/server/modules/base"
	"github.com/maxlandon/wiregost/server/security"
)

// Post - A module dedicated to post-exploitation activities
type Post struct {
	*base.Module
	Session *ghosts.Ghost // Session is an interface accepting different implants
}

// NewPost - Instantiates a new post, and handles base module instantiation
func NewPost() (post *Post) {
	post = &Post{&base.Module{}, nil}

	// Default Information filling
	post.Info.Type = modulepb.Type_POST

	return
}

// GetSession - Returns the Session corresponding to the Post "Session" option.
func (m *Post) GetSession(id uint32) (err error) {

	requested := ghosts.Ghosts.Get(id)
	ghost := requested.Core.Info()

	// We check permissions here and now, as we cannot pass
	// the module's context to each implant method call in module
	_, err = security.CheckCorePermissions(ghost, m.Client.User)
	if err != nil {
		return
	}
	m.Session = requested

	return
}

// CheckTarget - Check various elements about the target, for ensuring module can be run
func (m *Post) CheckTarget() (err error) {

	return
}
