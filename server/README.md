
## Wiregost Server


This directory emcompasses a large domain set: the Server is in charge of handling multiple users and their consoles,
managing sessions, coordinate the routing system , use and manage module of different types, setup and compile payloads
of various types and platforms, etc.

Therefore, this directory contain packages which might themselves be package roots for their respective domains, such
as `payload/`, `session/`, `transport/`, etc.

Note: some objects in Wiregost may be considered as modules (because they allow/need UI interaction) while being in
a different directory, such as the `transport.Module` type, which is not, *in the traditional Metasploit meaning*,
a module, because Wiregost merges the notion of `Transport` and `Handler` (a handler being a transport launched now).


---
### Directory Contents

The classification is a little more involved here:

#### Core

- `assets/`     - All configuration, directories and setup functions, as well as Packer packed assets.
- `certs/`      - Certificate infrastructure for Wiregost. Generates and check them for ghosts, user consoles, etc.
- `clients/`    - All code for connecting consoles to the server, authenticate them and register RPC functionality.
- `events/`     - Various events occur all around Wiregost server, implants and consoles. This handles and dispatches them.
- `generate/`   - Some compilation code for obfuscation, compiled user consoles and other stuffs is here.
- `log/`        - Wiregost has a powerful local/remote logging infrastructure, with different files for different areas.
- `jobs/`       - All concurrent jobs in the server, such as listeners, routers, proxies, binders, compilers, etc...
- `version/`    - Version code.

#### Modules

- `module/`     - Module objects definitions and their base methods, for several types (exploit, base, post, etc...).
- `payload/`    - All payload setup/generation helpers are in this directory. Should apply to all no matter Arch/OS/Transport.
- `exploit/`    - Exploit module definitions and core method set.

#### Post-Exploitation

- `ghosts/`     - All connected and/or pending ghost implants have a server-side object that is used by consoles and modules.
- `session/`    - Session types and their core method set are defined in this directory, which goes deep.

#### Network & Transport

- `c2/`         - Code for custom C2 connections/RPC handlers, such as Sliver's DNS, MTLS, and HTTPS handlers.
- `transport/`  - All the routing infrastructure and ghost connections are stored in this package.
