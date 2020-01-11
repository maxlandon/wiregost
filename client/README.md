### Client

The `client/` package contains the code necessary to run a client shell.
This code never imports any code from the `server/`.

* `assets/`     - Static asset files and management code (ex: client config)
* `commands/`   - Command implementations
* `constants/`  - Various shared constant values
* `core/`       - Client state management
* `events/`     - Handles events from Wiregost
* `help/`       - Console help
* `spin/`       - Console spinner library
* `transport/`  - Wires the client to the server (gRPC connection code)
* `version/`    - Version information
* `main.go`     - Entrypoint
