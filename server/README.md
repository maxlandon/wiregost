
## Wiregost Server

The `server/` directory contains all the code used by the server. All interactions with other components,
such as the Database, happen through client code generally contained in other root directories, like `db/`,
or `proto/` for Protobuf generated code.

---
### Directory Contents

- `assets/`     - All configuration, directories and setup functions, as well as Packer packed assets.
- `c2/`         - Code for custom C2 connections/RPC handlers, such as Sliver's DNS, MTLS, and HTTPS handlers.
- `certs/`      - Certificate infrastructure for Wiregost. Generates and check them for ghosts, user consoles, etc.
- `clients/`    - All code for connecting consoles to the server, authenticate them and register RPC functionality.
- `events/`     - Various events occur all around Wiregost server, implants and consoles. This handles and dispatches them.
- `generate/`   - The server can compile obfuscated ghost implants and pre-set obfuscated user consoles, here.
- `ghosts/`     - All connected and/or pending ghost implants have a server-side object that is used by consoles and modules.
- `jobs/`       - All concurrent jobs in the server, such as listeners, routers, proxies, binders, compilers, etc...
- `log/`        - Wiregost has a powerful local/remote logging infrastructure, with different files for different areas.
- `modules/`    - Module objects definitions and their base methods, for several types (exploit, base, post, etc...).
- `transport/`  - All the routing infrastructure and ghost connections are stored in this package.
- `version/`    - Version code.
