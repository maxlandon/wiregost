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
	"github.com/maxlandon/wiregost/server/ghosts"
	"github.com/maxlandon/wiregost/server/modules/base"
)

// Post - A module dedicated to post-exploitation activities
type Post struct {
	*base.Module
	Session ghosts.Session // Session is an interface accepting different implants
}

// NewPost - Instantiates a new post, and handles base module intanstantion
func NewPost() (post *Post) {
	post = &Post{&base.Module{}, nil}
	return
}

// GetSession - Returns the Session corresponding to the Post "Session" option.
func (m *Post) GetSession() (err error) {
	m.Session.ID()
	return
}

// CheckTarget - Check various elements about the target, for ensuring module can be run
func (m *Post) CheckTarget() (err error) {
	return
}
