
## Wiregost Networking Stack 

One of the main aims of Wiregost, through its Go implementation, is to provide advanced networking and routing capacities to an implant framework.
Therefore, the `c2/` package and all sub-packages in it constitute the routing mechanism used by the Server. This will allow:

- A large suite of network/transport/application protocols usable for traffic routing.
- Advanced implant traffic routing (ghost implant communications)
- Advanced non-implant traffic routing (everything talking TCP/UDP or anything else will be able to make use of Wiregost routes)
- An advanced routing management system: pluggable selection strategies permission controls, to build routes optimized for certain purposes.
- Advanced routing tree management for Wiregost users.

Below is an overview of the Routing System in Wiregost. All of these capacities will be further explained and detailed either in the 
subpackages' READMEs, in the Wiki documentation or in the code.

----
### Directory Contents

**Package files**
- `client.go`
- `route.go`
- `router.go`
- `selector.go`

**Subpackages**
- `route/`          - *Route*  All route building, storage and usage is done in this package. It has its own README
- `selector/`       - *Route*  Pluggable selection strategies, filters, permission controls, etc are in this package.
- `connector/`      - *Network*Code for handling the application layer of a connection.
- `transporter/`    - *Network*Code for handling the transport layer of a connection (handshake, etc...)
