
## Wiregost Protobuf Objects Library - V1

In the V1 source code directory (`proto/v1/src/`):

- `client/`     - All objects needed by the client
- `db/`         - All objects that do not belong to other directories, but need storage
- `ghost/`      - All objects, requests and responses for core implant functionality
- `module/`     - Module & option definitions, requests and responses
- `server/`     - All objects needed by the server (users, events, clients, etc...)
- `transport/`  - Transport objects as a whole (listeners, binders, routes, protocols, etc...) 

- prototool.yaml    - Protobuf project management file
- file-header.txt   - Header for `prototool create file.proto` tool

