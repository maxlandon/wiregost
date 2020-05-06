
## Wiregost Module System

Wiregost heavily relies on modules to perform tasks. It aims to be extensible, and the following module architecture is made 
so that it is easy to write new modules, and to have a simple access to many functionalities of the framework from within these modules.

In order to provide this, on the one hand, and for users to easily develop new modules, on the other hand, this README will explain
successively what are Wiregost Modules: several things at once. Read in order, and it will be easy.

----
### Modules are an interface

In order to handle several types of modules, with potentially many different architectures, Wiregost
handles and considers all modules through an interface `Module`. This allows to plug new modules, of any
complexity, to Wiregost, as long as they implement this interface:

```go
type Module interface {
	ParseMetadata() error                         // Parse module metadata
	ToProtobuf() *modulepb.Module                 // When consoles request a copy of the module
	Run(action string) (result string, err error) // Run one of the module's functions
	Option(name string) (opt *modulepb.Option)    // Get an option of this module
	CheckRequiredOptions() (ok bool, err error)   // Check all required options have a value
	Event(event string, pending bool)             // Send an event/message back to the console running the module
	Asset(string) (filePath string, err error)    // Find the path of an asset in the module directory.
}
```

This means, for those not familiar with interfaces, that you can write a module that has 10 packages, 300 functions and 
4000 variables because it needs to perform a very complicated task, and it will still be accepted by Wiregost as a valid
module, as long as it has the above methods implemented, and you will be able to use through the console.

**Fortunately**, you won't have to define these methods yourself: they are implemented by the `Base` module type.


----
### The Base module type

In Wiregost, the most basic module implementation is the `Base` type defined in `base/base.go`. This type *implements*
the `Module` interface, and nothing more. This means:

1. All modules types inheriting from it, will automatically be considered as valid Modules, but...
2. ...This type does not offer any specific and/or easy way to interact with the various components of the framework as a whole.

A few explanations on these base methods in the snippet above:
- Module valid targets, authors, descriptions, options, etc are defined in a JSON file. `ParseMetadata()` loads it.
- Wiregost uses gRPC & Protobuf for console-server communications. To this effect, The `Base` type embeds a protobuf Module object. 
  `ToProtobuf()` handles serialization.
- `Run(action string)` is triggers the module, with an optional parameter describing which action to perform.
- `Option(name string)` is used to query module options, inspect and/or set them.
- `CheckRequiredOptions()` checks all required options have a value.
- `Event(event string, pending bool)` is the way you can push status messages back to the console. If you browse the builtin module
  methods for a given type (more on this later), you will see they are used often.
- `Asset()` is explained below.


----
### Modules are language-agnostic repositories of code

1. When browsing the documentation on how to write modules, you will see that each module is defined in a Go package, that is imported
at compile-time by Wiregost. This means that your module package can be anywhere in the world, as long as you can import it in the
`register_modules.go` file.

2. Another very cool tool that exist in Go is [packr](https://github.com/gobuffalo/packr), which allows us to store any data of any kind
into the compiled version of Wiregost Server. Please go check packr documentation for an explanation of how it works. For those who know,
this means we store any compiled and/orsource files in any language, and Wiregost will automatically have access to them because they will always
be the server's private directory.

3. The base module thus exposes a function called `Asset(path string)`, that allows to retrieve a (most-of-the-time) non-Go file that you might
need for performing tasks on a host, such as Shell script, or uploading a custom binary. Wiregost also offers, through the console, a way to edit
and save those files.

4. Once again, your module can be split into 3 directories, 3 repositories, have 300 functions and 50 packages. As long as it embeds the `Base`
module type (or, but useless though, that you reimplement it yourself). This helps, for intance, when all non-Go files are in `src/` and all Go 
code is in `scanner` (or `whatever`)


----
### Modules are implemented by Type

Wiregost already provides 4 types of modules: `Handler`, `Payload`, `Post` and `Auxiliary`. There roles/categorizations is similar to Metasploit's
module types, albeit with a few differences. 

Because each of these modules target different areas of the framework, they need to provide their own set of methods. For instance:

- A `Handler` module might have to expose functions for manipulating the various transport mechanisms available in Wiregost (and a few other utility functions)
- While a `Post` module will need to access the various core capabilities of an implant, (which are transport-agnostic, by the way). 

This is an example of how Module types must differ in the toolset they give to the module user AND the module writer.





<!-- `Handler` is the equivalent of `exploit(multi/handler)` in Metasploit. -->
