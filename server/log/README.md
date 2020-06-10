
## Wiregost Logging & Debugging Infrastructure 


Naturally, a complete and flexible logging infrastructure is necessary to any implant framework.
Wiregost aims to provide it for all components, with a clear segregration of logs based on their subjects.

The logging infrastructure in Wiregost is built with the [logrus](https://github.com/sirupsen/logrus) package.

In Wiregost, each of the following components has a dedicated logger:
- The server
- Each ghost implant
- Client consoles
- Client Users


-----
### Server Logging

The server also has different files for logging different types of events:

- `clients.log`     - Records all client consoles connections, and all the certificate loading/verification process related to it.
- `compilation.log` - Records all events happening during the compilation of an implant.
- `listeners.log`   - All implant connections/disconnections from the server, and certificate loading/verification related to it.
- `server.log`      - All other logged events (asset/module unpacking).


-----
### Implants Logging

Each implant in Wiregost has its own log file, located in the corresponding ghost directory. The logging system for implants
aim to be more flexible than that of the server, given that it will generate more traffic (implants push events back to the server,
mostly for debugging purposes).

Therefore, it is possible to enable/disable the logging that an implant might push back to the server. There are two types of debugging
for implants: `local` (print events directly on the console running the implant, on the target) or `remote` (push back events to the server).


-----
### User logging

Wiregost also provides per-user logging: every request made by one of a user's consoles is logged, and this includes:
- Commands sent to an implant
- Modules used (no matter the type of module)


-----
### Centralized Logging

You may have thought, when reading about Wiregost logging infrastructure, that such a segregration of log events might impair on the overall
clarity of log/events streams (in other words, the ease for an auditor to build a global, chronological picture of ALL events happening in the
Wiregost system).

Therefore, some of the logs that have been evocated above will also be centralized in a fewer number of log files, in order to give a more 
"streamlined" view of events in Wiregost. These centralized log files are the following:

- `Here, define the files we need`


-----
### Directory contents

- `clients.go`      - Logs all consoles connections/disconnections, for all users.
- `ghost.go`        - Logging infrastructure for implant events.
- `user.go`         - Logs all actions from a user, or events related to him.
- `server.go`       - Logging infrastructure for various events (compilations, conns, etc)
- `hooks.go`        - Functions needed to implement `logrus` log package.
