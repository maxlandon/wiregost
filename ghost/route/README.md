
## Implant Routing Infrastructure 

The `route/` package contains all code for implants to route traffic through them. As said in the server `route/` directory,
the routing infrastructure of ghost implants is meant to route both implant-specific communications and non-related traffic
(such as Nmap scans, SSH traffic, etc).

Because most of route building, management and implementation is made on the server, ghost implants only have a subset of the
server's `route` package (listener, server, protocol-specific code, bypass code).

----
### Directory Contents 

- `server.go`       -
- `handler.go`      -
- `listener.go`     -
- `variables.go`    -

- `connector/`      - *Network* Code for handling the application layer of a connection.
- `transporter/`    - *Network* Code for handling the transport layer of a connection (handshake, etc...)
