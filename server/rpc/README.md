
## Wiregost RPC Layer

The `rpc` directory and its subpackages play an important role in Wiregost, despite what codebase readers will think:
"A ghost object for module methods, a Session object for transports, and a RPC object for methods again, this is redundant..."

----
### Problems calling for an intermediate layer

It appears Wiregost has to deal with different needs depending on the use case:

1. A console user does not care about how the whole thing is structured, because it just types commands.
2. A module *writer* cares about what methods are exposed, and wants them to be easy to call and use.
3. Wiregost, besides, has to check and enforce various permissions for implants, no matter how they are used (console or modules).
4. As a consequence of point 3, module writers should not have to deal with permission problems.
5. Performance concerns, such as timeouts and broadlier `context.Context` objects, are on Wiregost's charge, and both module writers and console users don't want to bother with it, appart from setting optional timeouts.

----
### The RPC layer: Handling all intermediary tasks

As a result of the points detailed above, we need an intermediate RPC layer, fulfilling all of these needs.
The RPC layer is in charge of:

1. Providing a complete list of RPC stubs for all major OS and major core areas. (See below for their use case)
2. Handling performance concerns (timeouts, etc)
3. Enforcing implant permissions (core functionality or routing infrastructure)

A typical RPC stub would look like this: 

```go
func (c *Client) Ls(ctx context.Context, req corepb.LsRequest) (res corepb.Ls, err error) {

        // Fetch the metadata object in the Context object
        // This metadata gives us information on the console's user that initiated the RPC stub, or the module's user.
        in := wctx.GetMetadata(ctx)

        // We check if the permissions for this implant allow it to perform the user's request
        ok, err := security.CheckPermissions(in.User, c.Ghost.ID)
        if !ok {
                return nil, err // We return if not allowed 
        }

        // We use Protobuf for serialization
        reqBytes, err := proto.Marshal(req)

        // Send through the custom C2 Conn, and we pass the timeout to it.
        res, err = c.Session.Request(ghostpb.LsRequest, ctx.Deadline(), reqBytes)  

        return
}
```

----
### RPC Layer stubs: Their use case

#### Modules

In another package (the `ghosts/*` packages) we also declare methods for requesting implants to perform such or such task. For instance, the `LS()` function above has a sister equivalent in these packages:

```go
func (g *Ghost) Ls(path string) (ls *corepb.Ls) {
	return
}
```

This method will itself call the RPC's Ls() function. However, above method is *only meant to be used by module users/writers*: this means we want an easier programmatic access to implant functionality, while keeping away all details from the module writer. As well, the return values might be other than protobuf objects, for easier results reuse.


#### Consoles 

Console users have much different needs: they don't care about function calls, and the console directly handle loads of responses/objects for printing them later. 
As well, for things like completion, we want to quickly & repeatedly call implant funcs, with performance and security needs but no programmatic concerns.
