
### CLI Interface

The overall workflow needs to be similar to Metasploit. However we need better: 
* Help (focused on some set of commands, if passed as argument)
* command list (avoid long list, add information to each command)
* Host, Service, Credential display.
* Dynamic prompt, if useful
* More color, for readability.


### Data Model

* I will not reivent the good wheel of Metasploit: Keep the same base data model.
* However, add more information such as AV capacities, and many others, especially 
  in order to refine the set of opsec-safe modules to use in a given situation.

* Same DB import/export capabilities would be nice.


### Workspace

* The workflow should revolve around workspaces even more than in Metasploit:
    * A workspace is an object that stores, processes and saves state.
    * This state is comprised of a Module Stack, a global_variables file, etc.
    * Each workspace has its own subdirectory in `~/.wiregost/workspaces/`. This
      helps compartimenting generated payloads, global_variables files, db export
      files, etc...


### Framework Control & Data Exposure

* The same WireGost session should provide bultin multi-client CLI capability.
Therefore:
* Run WireGost as a server, and preferably as a systemd service,
* The RPC interface should be provided with gRPC.


### Modules

* The bulk of the work will be to reimplement all the modules in Go.
* Therefore some of them (the newer ones) will appear more useful than others, such as Evasion modules.


### Payloads

* Some degree of interoperability/integration with PowerShell Empire would be a good thing.
* Although for the moment custom WireGost payloads and handlers are not a priority at all.


### Documentation

* Metasploit _usage_ documentation is great, but as far as the codebase is concerned, it is
  inexistant. The cryptic object model of the Ruby language is not helping either. Therefore:
* Better Codebase documentation.

