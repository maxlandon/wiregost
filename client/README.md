## Client

The `client/` directory contains the code necessary to run a client shell.
This code never imports any code from the `server/`.

* `assets/`                 - Static asset files and management code (ex: client config)
* `commands/`               - Command implementations
* `completers/`             - Root and Command-specific completers
* `console/`                - Console code (prompt, command dispatch, connect, etc...) 
* `constants/`              - Various shared constant values
* `core/`                   - Client/Server bind state management
* `events/`                 - Handles events from Wiregost server
* `help/`                   - Console help strings
* `spin/`                   - Console spinner library
* `transport/`              - Wires the client to the server (MTLS connection code)
* `util/`                   - Various utilities needed by the client/console
* `version/`                - Version information
* `wiregost-console.go`     - Entrypoint
