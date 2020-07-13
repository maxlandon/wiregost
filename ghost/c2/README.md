
## Implant Transport & Routing Infrastructure 

The `c2/` directory contains all the code necessary to implants for handling:
- Their own connection.
- Traffic to be routed to further implants.

Because many inverse flows of data will result from this architecture, it is vital 
to minimize port opening on target machines. SO_REUSEPORT() should be used whenever
possible.


----
### Directory Contents 

- `route/`          - *Route*  All route building, storage and usage is done in this package. It has its own README
- `connector/`      - *Network* Code for handling the application layer of a connection.
- `transporter/`    - *Network* Code for handling the transport layer of a connection (handshake, etc...)
