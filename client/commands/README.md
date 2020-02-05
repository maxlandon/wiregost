## Commands

The `commands` package contains all commands available to a Wiregost client shell. 
They are regrouped by topic/specialisation

----
#### Setup
* `command.go`              - Command/SubCommands/Arguments definitions, and command mapping.
* `register-commands.go`    - Function for registering all commands, called during console instantiation

#### Core Shell 
* `core.go`                 - Core commands (local to shell) like cd, resource, exit, shell exec...
* `help.go`                 - Help printing for all commands

#### Data Service
* `workspace.go`            - Manages Wiregost workspaces (data_service)
* `host.go`                 - Manages hosts (data_service)
* `service.go`              - Manages services (data_service)
