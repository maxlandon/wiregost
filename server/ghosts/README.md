
## Wiregost Implants API 

### Implants are living behind an Interface

Similar to the [module structure](https://github.com/maxlandon/wiregost/tree/v1.0.0/server/modules/README.md), implants in Wiregost
are first and foremost defined through the following interface:

```go
type Session interface {
	ID() (id uint32)                          // Session ID
	OS() (os string)                          // Session Operating system
	Owner() (owner *serverpb.User)            // User owning the session
	Permissions() (perms ghostpb.Permissions) // Who has the right to use implant
	Request()                                 // Function for sending a message to implant (transport-agnostic)
}
```

Behind this interface will live server-side objects representing connected implants. These server side objects are able to track
state from the remote ghost, or to send requests to it.

The aim of this interface is to allow for implants that have a different set of methods: thus, a different set of capabilities.
For instance, a `Ghost` type representing a connected ghost running on a Windows target (in `windows/ghost.go`), will have 
a set of functions for manipulating the registry.


### How is this similar to Metasploit implants API ?

In Metasploit, there is no magic as to why we can write post-modules:
Each (remote) implant has an appropriate method for such or such task, and its equivalent server-side object has an equivalent method
that will trigger the remote method (yes, an RPC). Thus, because the server-side object offers many functionalities in Metasploit, 
people can write non-trivial post-exploitation modules. In addition, these methods belong to some type like `MSF::Post::Windows_x86`.


### The ghosts.Base type

Thus, Wiregost has a server-side object for each connected implant, and this server side object is a type that automatically implements
the `Session` interface, because it embeds the `Base` type you will find in `ghosts/base.go`.

This type merely implements the `Session` interface, and nothing more. This means you can embed this type in any other custom type, and it
will automatically be considered a ghost session. This, in addition, permits an implant API that is OS-specific, but transport-agnostic.


### Types of implants in Wiregost

The aim of Wiregost, in a first round, is two provide 3 different OS-specific types:
- `ghosts/darwin/ghost.go`
- `ghosts/windows/ghost.go`
- `ghosts/linux/ghost.go`

You will notice in the code that all types are called `Ghost`. That is not an issue though, because Go knowns these types are different
because they belong to a different package.

In each of these packages, we will define server-side methods for executing functions on the remote ghost implant.
