## Modules 

The `modules` directory contains all modules available in Wiregost 

The following are available modules types:

* `exploit`         - All exploits, as well as listeners generation modules
* `payload`         - These modules are sorts of stubs to Wiregost implant generation code.
* `post`            - All modules that can be run on a Ghost agent (privesc, recon, persistence, etc...)
* `auxiliary`       - All modules that do not trigger any vulnerability (scanners, proxies, etc...)

----
## Bundled Module

The main goal of the following package structure, and the model of "bundled modules"
is to provide locality of information for each module. However, some modules may only use
a subset of the available bundle structure (ie payloads, scanners, listeners)

* `module.go`           - Core functionality of the module, with Run(), Init(), Reload() functions.
* `metadata.json`       - Module information and options, loaded during module initialization
* `lib/`                - Code used by the module's core if needed (only Go code in there)
* `docs/`               - Any documentation relevant to the module
* `src/`                - Any non-Go source code needed by the module (ie PowerShell scripts for a Post module)
* `data/`               - Mostly any non-Go executable file, or platform-specific build information, needed by the module

----
### Notes

* The `module.go` file in each template directory should always have this name, as Wiregost will look for it when loading/using the module
* This same file is a template, with core function signatures, that the user just needs to fill with module-specific logic.
