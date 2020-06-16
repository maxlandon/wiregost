
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


### The `generic.Ghost` type

Thus, Wiregost has a server-side object for each connected implant, and this server side object is a type that automatically implements
the `Session` interface, because it embeds the `Ghost` type you will find in `ghosts/generic/ghost.go`.

This type has two main roles:
- It implements the `Session` interface, so that any type embedding it will automatically considered a valid ghost session.
- It provides all cross-platform functionality, such as filesystem methods. These methods, generally implemented with the Go standard library,
  will work on any Operating System and architecture.

This architecture, in addition, permits an implant API that is OS-specific, but transport-agnostic.


### Types of implants in Wiregost

The aim of Wiregost, in a first round, is two provide 3 different OS-specific types:
- `ghosts/darwin/ghost.go`
- `ghosts/windows/ghost.go`
- `ghosts/linux/ghost.go`

You will notice in the code that all types are called `Ghost`. That is not an issue though, because Go knowns these types are different
as they belong to different packages. They also embed the `generic.Ghost` type, so that they automatically provide generic methods
such as filesystem manipulation.

In each of these packages, we will define OS-specific server-side methods for executing functions on the remote ghost implant.


### Details on each implant type

All server-side code (thus, requests) pertaining to a certain Operating System/Architecture will live into one of these subdirectories:
- `linux/` for Linux post methods
- `darwin/` for MacOS post methods
- `windows/` for Windows post methods.

Note that each of the ghost types declared in these packages may have specific attributes, so that it is easy to access all the "ecosystem"
pertaining to a target. Because these ghost types always include the `generic.Base` type, they are still valid server-side ghost objects.

Go check these subdirectories for more ! They have their own documentation.
