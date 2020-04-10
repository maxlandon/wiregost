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

package console

import (
	"github.com/maxlandon/wiregost/client/commands/canaries"
	"github.com/maxlandon/wiregost/client/commands/core"
	"github.com/maxlandon/wiregost/client/commands/filesystem"
	"github.com/maxlandon/wiregost/client/commands/ghosts"
	"github.com/maxlandon/wiregost/client/commands/jobs"
	"github.com/maxlandon/wiregost/client/commands/module"
	"github.com/maxlandon/wiregost/client/commands/profiles"
	"github.com/maxlandon/wiregost/client/commands/server"
	"github.com/maxlandon/wiregost/client/commands/sessions"
	"github.com/maxlandon/wiregost/client/commands/stack"
)

// RegisterCommands - Registers all commands to the parser
// Needs to be here to avoid circular import, and because init funcs don't work.
func RegisterCommands() {

	// Main Context ------------------------------------------------------------
	// Core
	core.RegisterCd() // cd
	core.RegisterLs() // ls

	// DB

	// Server
	server.RegisterServer()        // server
	server.RegisterServerConnect() // connect

	// Module
	module.RegisterModuleUse()          // use
	module.RegisterModuleInfo()         // info
	module.RegisterModuleShowOptions()  // options
	module.RegisterModuleSetOption()    // set
	module.RegisterToListener()         // to-listener
	module.RegisterModuleParseProfile() // parse-profile
	module.RegisterModuleToProfile()    // to-profile
	module.RegisterModuleRun()          // run
	module.RegisterModuleBack()         // back

	// Stack
	stack.RegisterStack()    // stack
	stack.RegisterStackUse() // use
	stack.RegisterStackPop() // pop

	// Profiles
	profiles.RegisterProfiles()       // profiles
	profiles.RegisterProfilesDelete() // delete

	// Ghosts
	ghosts.RegisterGhosts() // ghosts

	// Canaries
	canaries.RegisterCanaries() // canaries

	// Jobs
	jobs.RegisterJobs()        // jobs
	jobs.RegisterJobsKill()    // kill
	jobs.RegisterJobsKillAll() // kill-all

	// Sessions
	sessions.RegisterSessions()          // sessions
	sessions.RegisterSessionsInteract()  // interact
	sessions.RegisterSessionsKill()      // kill
	sessions.RegisterSessionsKillAll()   // kill-all
	sessions.RegisterSessionBackground() // background

	// Ghost Context ------------------------------------------------------------
	// Filesystem
	filesystem.RegisterGhostCd()       // cd
	filesystem.RegisterGhostLs()       // ls
	filesystem.RegisterGhostCat()      // cat
	filesystem.RegisterGhostPwd()      // pwd
	filesystem.RegisterGhostRm()       // rm
	filesystem.RegisterGhostMkdir()    // mkdir
	filesystem.RegisterGhostDownload() // download
	filesystem.RegisterGhostUpload()   // upload

	// Info

	// Priv

	// Proc

	// Execute
}
