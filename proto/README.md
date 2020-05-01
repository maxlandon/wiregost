
## Wiregost Protobuf Objects Library

In this directory are stored all Protobuf objects used in Wiregost. 
There are several aims to be fulfilled with Protobuf in this project:

- Formalize the set of capacities of Wiregost and its implants in a language-agnostic way.
- Separate functionality based on Protobuf API version.
- Enable development for implants in other languages.
- Enable easier RPC development, given the wide RPC support of Protocol Buffers.
- Implement a language-agnostic data model oriented toward security purposes. 

----
### Protobuf management

For managing the Protobuf file repository, we use [uber/prototool](https://github.com/uber/prototool).
This tool allows an easier way of dealing with compilation, linting (coding standards), default fields, etc...
Thanks to the `prototool.yaml` configuration file present in each version source directory (`wiregost/proto/v1/src`, `wiregost/proto/v2/src`),
we can easily manage, on a per-version basis:

- Language output & import options
- Coding standards used (ex: forbidden patterns in fields)
- Use of output plugins (for tags, generation, etc...)

This tool will further ease development of Wiregost server/implant functionality, cross-language, and with-easy-to reproduce builds.
As you will see in their documentation, there is support at least for Vim, which works well. 
(You'll need [ALE](https://github.com/dense-analysis/ale), but the whole is very simple to install).

----
### Wiregost Protobuf Objects versions

- `v1`      - First iteration build around GORM DB usage, with a Nmap-like object model for DB many objects

----
### Example Directory Structure (V1)

In the V1 source code directory (`proto/v1/src/`):

- `client/`     - All objects needed by the client
- `db/`         - All objects that do not belong to other directories, but need storage
- `ghost/`      - All objects, requests and responses for core implant functionality
- `module/`     - Module & option definitions, requests and responses
- `server/`     - All objects needed by the server (users, events, clients, etc...)
- `transport/`  - Transport objects as a whole (listeners, binders, routes, protocols, etc...) 

----
- `prototool.yaml`    - Protobuf project management file
- `file-header.txt`   - Header for `prototool create file.proto` tool


----

### Generated Language Files

Generated code for any language in any version is output in `proto/version/gen/language/`. 
For example, Go generated source V1 is output to `proto/v1/gen/go/`.

----

### Go Tags 

We use struct tags for easier parsing in various cases:
- XML Parsing used for Nmap scans
- GORM tags. They are used by the ORM layer we use (GORM).

We make use of [https://github.com/favadi/protoc-go-inject-tag] for injecting tags directly on the 
generated `.pb.go` files. We place comments like  `// @inject_tag: gorm"not null" xml:"name"` in our
Protobuf source files, and the binary parses them and outputs the corresponding struct tag.

When running `make tags` from the repo root, it will recursively generate struct tags
for each .pb.go file in the `proto/` directory, with the `generate-go-tags.sh` script.
