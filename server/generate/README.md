## Generate

The `generate` package is responsible for generating Ghost implant binaries, such as executables and shared libraries

* `ghost-configs.go`        - Config structs for ghost  implant generation, and their protobuf methods
* `canaries.go`             - Code for manipulating canaries in implants
* `code-rendering.go`       - Depending on an implant's C2s, the implant code will vary. This file renders the appropriate Go code
* `codenames.go`            - Provides code names for implants
* `generation.go`           - Functions for generating ghost executables and shared libraries
* `ghosts.go`               - Code for managing ghost implant builds (get/save/etc...)
* `parse-profile.go`        - Functions used mainly by modules, in order to parse profile-saved C2s
* `profiles.go`             - Code for managing profiles (get/save/etc...)
* `srcfiles.go`             - Used by rendering function to locate necessary implant source files
* `c-compiler.go`           - Function for getting the appropriate MinGW C compiler
