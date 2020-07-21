
## Wiregost Implant Code

The `ghost` directory contains all code necessary to Wiregost implants. All compiled code in a ghost implant has a source here.
The code is separated according to the following principles: the root `ghost` directory contains all packages that make up the
backbone of a ghost implant (comms, concurrency, logging, assets, RPC, etc).
The `core` package contains two things: implants entrypoints (executables) for each major supported OS (Linux, Darwin, Windows),
in their respective `linux`, `darwin`, `windows` packages, along with all OS-specific code. All the other packages in `core` are
separated by task/topic/area, and usually contain files for all OS, with conditional compilation with `_windows` suffixes.

----
### Directory Contents 

- `core/`       - Implant entrypoints (separated by OS), and core functionality (separated by OS and/or area).
- `log/`        - Local & remote logging infrastructure for implants.
- `c2/`         - Code used by implants for traffic routing and transport protocols, usually separated by OS.
- `rpc/`        - RPC interfaces and handlers, in order to provide easy-to-code cross-protocol functionality.
- `channels`    - Implants support multiple processes (system shells & commands, routers, listeners, etc).
- `assets`      - All directories, compiled variables, configurations and other basic setup is done here.
- `profile`     - Implants continuously monitor their performance and resource usage.
- `security`    - All security objects, details, and functions used by implants (load configs, secure exits, etc).

