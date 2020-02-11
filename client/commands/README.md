## Commands

The `commands` package contains all commands available to a Wiregost client shell. 
They are regrouped by topic/specialisation

----
#### Setup
* `command.go`              - Command/SubCommands/Arguments definitions, and command mapping.
* `register-commands.go`    - Function for registering all commands, called during console instantiation
* `shell-state.go`          - Passes the client shell context/variables to commands for read/write access

#### Core Shell 
* `core.go`                 - Core commands (local to shell) like cd, resource, exit, shell exec...
* `help.go`                 - Help printing for all commands

#### Data Service
* `workspace.go`            - Manages Wiregost workspaces
* `host.go`                 - Manages hosts
* `service.go`              - Manages services

#### Stack / Modules
* `stack.go`                - Loads/unloads modules onto the workspace's module stack
* `module.go`               - Interact with current module (run, set options, info, etc...)

#### Jobs
* `job.go`                  - Manages jobs (show, kill, kill-all)
