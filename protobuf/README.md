## Protobuf

The `protobuf/` directory contains all protobuf message definitions and gRPC
stubs needed in Wiregost.
It is organized in two different packages:

* `client/`    - Generally refered as `clientpb` in the code, these messages/stubs should only be send from client to server, or conversely.
* `ghost/`     - Referred as `serverpb` in the `/server/` code, or `pb` in the `/ghost/` code, these messages/stubs may be sent from client to the server, or from server to agents and vice versa.
