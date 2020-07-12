
## Implant RPC Infrastructure 

Given the aim of providing multiple transports to implants for their communication, we need an intermediate layer of RPC code, in order to provide:
- Transport-agnostic functionality.
- Easy-to-extend core functionality AND tranport capacity.

There are some good RPC frameworks in Go out there, and each of them may support different sets of transport protocols.
Therefore, to provide an easily extendable RPC layer, we split the code per underlying RPC framework, unless the infrastructure
is made "in-house", and therefore will be stored in `custom/`.

In any case, all RPC stubs make use of Protobuf for message serial/deserialization.

----
### Directory Contents 

- `grpc/`       - GRPC is based on TCP/HTTP2.
- `rpcx/`       - RPCX supports TCP/QUIC/HTTP/KCP.
- `custom/`     - RPC stubs used by custom transports (mTLS, DNS).
