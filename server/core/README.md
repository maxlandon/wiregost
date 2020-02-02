### Core 

The `core` package contains data structures and methods for managing state from agents & clients.

* `core.go`         - Small utility functions
* `clients.go`      - Containing, setting and getting client state from server-side.
* `ghosts.go`       - Containing all data structures and methods for representing the state of a Ghost implant.
* `job.go`          - Containing references to all background jobs in Wiregost
* `tunnels.go`      - Contains code for create/closing tunnels, each being a simple programmatic (no network code here) mapping between a Client and a Ghost objects.
