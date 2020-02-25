## RPC 

The `rpc` package contains the RPC APIs. They are responsible for most of the server-side logic.

They are invoked remotely by a client console, connected via the `transport` package.

---
#### Setup 
* `rpc.go`              - Mapping RPC handlers to protobuf Message types and tunnel handlers.
* `rpc-tunnels.go`      - Functions for tunnel creation/destruction when client needs them. 

#### Main Menu Console Commands
* `jobs.go`             - Manage jobs 
* `stack.go`            - Module stack management 
* `module.go`           - Module management (run, set options, etc...) 
* `users.go`            - Manage users 
* `profiles.go`         - Manage implant profiles 
* `ghosts.go`           - Manage generated implants (builds)
* `sessions.go`         - Manage connected implants 

#### Implant Menu Console Commands
* `agent-info.go`       - Implant/target information 
* `filesystem.go`       - Target filesystem management
* `priv.go`             - Target privileges management 
* `proc.go`             - Target processes management 
* `execute.go`          - Shellcode/Assembly/DLL/payload execution and injection.

