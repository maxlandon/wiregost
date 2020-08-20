
## Ghost Session

In Wiregost, a `Ghost` is the equivalent of Metasploit's `Meterpreter`, but with a different skillset.
Because this implant may run on different Operating Systems and Architectures with varying capacities,
many objects may live behind a core `Ghost` Session object.

This `ghost` package and its subpackages implement the core of a Ghost implant *client*, often reaching
over other packages of the framework in order to provide its full skillset: networking and routing, for
instance, is not really part of the core `ghost` client, but the latter might have to make use of (or 
more generally, interact with) these `transport` objects.

Generally and when possible, the aim was to mimic the structure of a `Meterpreter` object, when the wheel
was already round. However, aspects like implant concurrency management (`channels` in Metasploit, and here 
also), and handled in a fundamentally different way, as experienced Go and/or Ruby will acknowledge.


---
### Directory Contents

- `client.go`       - Holds the `Client`, main object through which we control a ghost, and its core methods. 
- `transport.go`    - Code for managing and using the implant's transports, themselves working with RPC stubs if needed.
- `request.go`      - Methods for sending RPC requests to implants, through their logical connection.
- `channels.go`     - Implant concurrency management (implants can have multiple `channels`).
