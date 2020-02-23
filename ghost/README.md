Ghost 
======

This directory contains all of the Ghost implant code.

* `constants/`      - Various shared constant values 
* `evasion/`        - Evasion code 
* `handlers/`       - Code for handling C2 requests and tunnels 
* `limits/`         - Code for managing/observing imposed limits on implant (hostname, username, datetime) 
* `priv/`           - Code for managing privileges, mostly Windows-related now. 
* `procdump/`       - Code for process dumping 
* `proxy/`          - Code to detect/use/set proxies on the target 
* `ps/`             - Manage/view target processes 
* `shell/`          - Code for getting/using system shells on target 
* `syscalls/`       - Library for using Windows syscalls 
* `taskrunner/`     - Code for process injection/migration 
* `transports/`     - Wires the implant to the C2 server 
* `version/`        - Implant version 
* `winhttp/`        - Code for using Windows HTTP APIs (used by proxy package) 
* `dllmain.c`       - Code used by implants compiled as shared libraries 
* `ddlmain.h`       - Headers for implants as shared libraries 
* `dllmain.go`      - (Secondary) Entrypoint for implants as shared libraries 
* `ghost.go`        - Entrypoint for all implants 
