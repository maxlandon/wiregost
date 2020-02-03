## Log 

The `log` package provides wrappers around Logrus loggers. Several loggers are used, depending on use case.

* `log.go`      - Generic functions/hooks available to all loggers.
* `server.go`   - Server logger, managing events like clients, compilation, etc.
* `ghost.go`    - Ghost-specific logger, saving all events for a given ghost in the appropriate directory.
* `audit.go`    - Records client connections/disconnections
