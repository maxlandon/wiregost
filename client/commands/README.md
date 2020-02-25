## Commands

The `commands` package contains all commands available to a Wiregost client shell. 
They are regrouped by topic/specialisation

----
#### Setup
* `command.go`              - Command/SubCommands/Arguments definitions, and command mapping.
* `register-commands.go`    - Function for registering all commands, called during console instantiation
* `shell-state.go`          - Passes the client shell context/variables to commands for read/write access

#### Core Shell & Main Menu
* `core.go`                 - Core commands (local to shell) like cd, resource, exit, shell exec...
* `help.go`                 - Help printing for all commands
* `server.go`               - Server connection and management 
* `users.go`                - Manage users of a Wiregost server 
* `stack.go`                - Loads/unloads modules onto the workspace's module stack
* `module.go`               - Interact with current module (run, set options, info, etc...)
* `job.go`                  - Manages jobs (show, kill, kill-all)
* `profile.go`              - Implant profiles management 
* `ghosts.go`               - Manage generated implants (builds)
* `sessions.go`             - Manage and interact with connected implants 

#### Implant Menu
* `agent-help.go`       - Help commmands for implant menu
* `agent-info.go`       - Implant/target information commands 
* `filesystem.go`       - Target filesystem management 
* `priv.go`             - Target privileges management 
* `proc.go`             - Target processes management 
* `execute.go`          - Shellcode/Assembly/DLL/payload execution and injection.

#### Data Service
* `workspace.go`            - Manages Wiregost workspaces
* `host.go`                 - Manages hosts
* `service.go`              - Manages services

