
## Wiregost Protobuf Objects Library - V1

In this directory are stored all Protobuf objects used in Wiregost. 
There are several aims to be fulfilled with Protobuf in this project:

- Formalize the set of capacities of Wiregost and its implants in a language-agnostic way.
- Separate functionality based on Protobuf API version.
- Enable development for implants in other languages.
- Enable easier RPC development, given the wide RPC support of Protocol Buffers.

----
### Protobuf management

For managing the Protobuf file repository, we use [uber/prototool](https://github.com/uber/prototool).
This tool allows an easier way of dealing with compilation, linting (coding standards), default fields, etc...
Thanks to the `prototool.yaml` configuration file present in each version directory (`wiregost/proto/v1`, `wiregost/proto/v2`),
we can easily manage, on a per-version basis:

- Language output & import options
- Coding standards used (ex: forbidden patterns in fields)
- Use of output plugins (for tags, generation, etc...)

This tool will further ease development of Wiregost server/implant functionality, cross-language, and with-easy-to reproduce builds.
As you will see in their documentation, there is support at least for Vim, which works well. 
(You'll need [ALE](https://github.com/dense-analysis/ale), but the whole is very simple to install).


----
### Directory contents

- `client/`     - All objects needed by the client
- `db/`         - All objects that do not belong to other directories, but need storage
- `ghost/`      - All objects, requests and responses for core implant functionality
- `module/`     - Module & option definitions, requests and responses
- `server/`     - All objects needed by the server (users, events, clients, etc...)
- `transport/`  - Transport objects as a whole (listeners, binders, routes, protocols, etc...) 


----

### Files

- prototool.yaml    - Protobuf project management file
- file-header.txt   - Header for `prototool create file.proto` tool
