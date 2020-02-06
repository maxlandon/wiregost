## RPC 

The `rpc` package contains the RPC APIs. They are responsible for most of the server-side logic.

They are invoked remotely by a client console, connected via the `transport` package.

* `rpc.go`              - Mapping RPC handlers to protobuf Message types and tunnel handlers.
* `rpc-tunnels.go`      - Functions for tunnel creation/destruction when client needs them. 
