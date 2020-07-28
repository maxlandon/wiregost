
## Custom DNS Connection 

The `dns` package contains a refactoring of Sliver's DNS C2 implementation. This refactoring
has also for aim to make this DNS protocol to implement the `net.Conn` object.

----
### Directory Contents 

- `dns_conn.go`      - Definition of the DNS connection, implementing `net.Conn`
