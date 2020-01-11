### Protobuf

The `protobuf/` directory contains all protobuf message definitions and gRPC
stubs needed in Wiregost.
It is organized in two different packages:

* `client/`     - Generally refered as `clientpb` in the code, these messages/stubs should only be send from client to server, or conversely.
* `server/`     - Referred as `serverpb` in the `/server/`, or `pb` in the `/ghost/` code, these messages/stubs may be sent from client to server, or from server to agents, and vice versa.

Each package contains both the protobuf message definitions and their respective gRPC generated stub files.
