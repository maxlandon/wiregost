
## Wiregost Implant Code

The `ghost` directory contains all code necessary to Wiregost implants. All compiled code in a ghost implant has a source here.

----
### Directory Contents 

- `core/`       - Implant entrypoints (separated by OS), and core functionality (separated by OS and/or area).
- `log/`        - Local & remote logging infrastructure for implants.
- `c2/`         - Code used by implants for traffic routing and transport protocols, usually separated by OS.
- `rpc/`        - RPC interfaces and handlers, in order to provide easy-to-code cross-protocol functionality.

