## Templates

The `templates` directory contains a template bundled module for each type of module available in Wiregost.
The core structure of all modules is the same, however each module may only use a subset of this structure.

----
## Package structure

* `module.go`       - Module struct, implementing a module properties and options (here to avoid circular dependencies)
* `exploit`         - Exploit module template
* `payload`         - Payload module template 
* `post`            - Post module template 
* `auxiliary`       - Auxiliary module template 

