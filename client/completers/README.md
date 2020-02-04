## Completers

The `completers` package contains all completers used in the console. It works as follows:
A main/root completer is registered during console instantiation, and provides completion only
for root commands.

Each time a root command needs further completion (for its subcommands, arguments, or data_service objects),
the root completer calls specialized functions contained in their appropriate source files.

* `completer.go`        - Main console completer, calling other specialized completer functions
* `workspace.go`        - Workspace subcommands and objects
* `hosts.go`            - Hosts subcommands and objects
