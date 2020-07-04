
## Wiregost Routing System 

The `route/` package contains all code necessary to build, store, modify and use network routes for implant and non-implant communication.
It aims to be more sophisticated, although more maintainable than Metasploit's routing system: 
- More protocols are supported.
- Routes can be modified and stored for further reuse.
- Many (pluggable) filters and selection strategies are used to build routes that are efficient in all cases. (More on this later)

----
### Directory Contents

- `chain.go`        -
- `node.go`         -
- `nodegroup.go`    -
- `server.go`       -
- `handler.go`      -
- `listener.go`     -
- `variables.go`    -
