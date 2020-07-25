
## Wiregost C2 Infrastructure

The `c2/` directory contains the code for C2 communications. It has several aims in sight:

- Provide a single, simplified interface to an implant communications, no matter the underlying protocol.
- Provide multiple logical connections over a single physical (reverse) connection, for stealthy routing.

---
### Multiplexing: Solving Pivoting and Routing Problems

One big aim of Wiregost is to provide cutting-edge routing capabilities to implants, while preserving
stealth and core post-exploitation functionality. Therefore, many things have to be assumed (or assured):

- Implants will most of the time reverse call the server (bind requests will be dropped by firewalls/NATs).
- Consequently, all communications that we want to send through have to "seem" being initiated from the ghost implant,
  or be handled in the same physical connection opened when it registered.
- Therefore, and until a better, wider system is devised, all protocols used for pure implant requests/responses must
  be multiplexable, at the request of either end.

The Go language and libraries allow us to multiplex connections for various transport/application protocols, such as
TCP (and therefore, SOCKs), KDP (UDP), SCTP, HTTP, Named Pipes...

This is mainly thanks to libraries like [yamux](https://github.com/hashicorp/yamux) or [smux](https://github.com/xtaci/smux), which accept and/or output abstract objects (interfaces) like `net.Listener` and `net.Conn`, for the latter it is just a logical connection along others, in *the same physical connection*.

Finally, this multiplexing system allows for a clear and logical separation of ghost implant sessions, even if they communicate over the same wire.

---
### Links between Implant Connections and the Routing System 

Because all traffic routed through our implants will have to go through the same physical connections they use for core communications, this means the `net.Conn` objects (or affiliates) used by `Sessions` in the `c2` directory will be reused by the `transport/route` package. 

Therefore, these `net.Conn` should be easily usable from outside the `c2` directory and its subpackages.

---
### Directory Contents 

- `dns/`        - DNS communications (does not support multiplexing, therefore no routes)
- `mtls/`       - TCP + TLS communications (supports multiplexing)
- `https/`      - HTTP(S) communications (supports multiplexing)
- `sctp/`       - Stream Control Transport Protocol (supports multiplexing)


- `session.go`  - Any ghost connection, no matter the protocol, is represented and managed as a `Session` object.


<!--  NOTE: ARCHITECTURE IS BUILT AROUND 4 CASES: -->
<!-- -------------------------------------------------- -->
<!--  -->
<!-- The client is the Wiregost server, talking to a in implant (server) on the target. -->
<!--     1) Target and client are both behind separate NATs -->
<!--     2) Target and client are both behind same NATs -->
<!--     3) Target is behind a NAT whereas the client isn't, and has a global IP address. -->
<!--     4) Client is behind a NAT whereas the target isn't, and has a global IP address. -->
<!--  -->

