## Module

The `module` package contains the implementation of modules in Wiregost, as well
as templates for each module type. The goal is obviously to provide maximum reusability
and modularity to Wiregost's components and agents.

Go in the `templates/` directory for further information on each bundled module type.

----
## Package Structure

* `modules.go`      - Implementation of module interfaces, (defining their usage in Wiregost)
* `load.go`         - Function for loading all modules in Wiregost
* `stack.go`        - Functions for managing modules stacks for users (loading/unloading/init ...) 
* `templates/`      - Containing module and directory templates for each type of bundled module (exploit, post, auxiliary)

----
## Bundled Module

The main goal of the following package structure, and the model of "bundled modules"
is to provide locality of information for each module. However, some modules may only use
a subset of the available bundle structure (ie payloads, scanners, listeners).
When reusing one of the bundle templates, feel free to add directories and files, and to remove the unneeded ones.

* `module.go`           - Core functionality of the module, with Run(), Init() and SetOption() functions.
* `metadata.json`       - Module information and options, loaded during module initialization
* `lib/`                - Code used by the module's core if needed (only Go code in there) (Optional)
* `docs/`               - Any documentation relevant to the module (Optional)
* `src/`                - Any non-Go source code needed by the module (ie PowerShell scripts for a Post module) (Optional)
* `data/`               - Mostly any non-Go executable file, or platform-specific build information, needed by the module (Optional)

----
### Notes

* The `module.go`file is a template, with core function signatures, that the user just needs to fill with module-specific logic.
Most of things that should be done in order to have a working module will be either explained fully in the [documentation for Writing Modules](https://github.com/maxlandon/wiregost/wiki/Writing-Modules), or in the `module.go` file of the templates
