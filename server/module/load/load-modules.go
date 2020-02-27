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

package load

// ********************* Adding Modules ****************************//

// Each time a module is added to Wiregost, a line like these below
// has to be added to the LoadAllModules() function.

// The import path to the module has to be added to the imports as well.

import (
	// Core
	. "github.com/maxlandon/wiregost/server/core"

	// Payloads
	dns "github.com/maxlandon/wiregost/modules/payload/multi/single/reverse_dns"
	https "github.com/maxlandon/wiregost/modules/payload/multi/single/reverse_https"
	mtls "github.com/maxlandon/wiregost/modules/payload/multi/single/reverse_mtls"
	multi "github.com/maxlandon/wiregost/modules/payload/multi/single/reverse_multi_protocol"
	http_stager "github.com/maxlandon/wiregost/modules/payload/multi/stager/reverse_http"
	https_stager "github.com/maxlandon/wiregost/modules/payload/multi/stager/reverse_https"
	tcp_stager "github.com/maxlandon/wiregost/modules/payload/multi/stager/reverse_tcp"

	// Post
	mimipenguin "github.com/maxlandon/wiregost/modules/post/linux/x64/bash/credentials/MimiPenguin"
	minidump "github.com/maxlandon/wiregost/modules/post/windows/x64/go/credentials/minidump"
)

// LoadAllModules - Load all modules in the modules directory.
func LoadModules() {

	// Payload -------------------------------------------------------------//
	AddModule("payload/multi/single/reverse_dns", dns.New())
	AddModule("payload/multi/single/reverse_mtls", mtls.New())
	AddModule("payload/multi/single/reverse_https", https.New())
	AddModule("payload/multi/single/reverse_multi_protocol", multi.New())
	AddModule("payload/multi/stager/reverse_tcp", tcp_stager.New())
	AddModule("payload/multi/stager/reverse_http", http_stager.New())
	AddModule("payload/multi/stager/reverse_https", https_stager.New())

	// Post ----------------------------------------------------------------//
	AddModule("post/windows/x64/go/credentials/minidump", minidump.New())
	AddModule("post/linux/x64/bash/credentials/MimiPenguin", mimipenguin.New())

	// Exploit -------------------------------------------------------------//

	// Auxiliary -----------------------------------------------------------//
}
