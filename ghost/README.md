
## Wiregost Implant Code

The `ghost` directory contains all code necessary to Wiregost implants. All compiled code in a ghost implant has a source here.

----
### Directory Contents 

- `core/`       - Implant entrypoints (separated by OS), and core functionality (separated by OS and/or area).
- `route/`      - Code used by implants for traffic routing.
- `rpc/`        - RPC interfaces and handlers, in order to provide easy-to-code cross-protocol functionality.
- `transport/`  - Boilerplate code for handling transport protocols, usually separated by OS.

