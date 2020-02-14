## Server

The `server/` directory contains the Wiregost server implementation.

* `assets/`         - Static assets embedded in the server binary, and their associated methods
* `c2/`             - The server-side command and control implementation
* `certs/`          - X509 certificate generation and management code
* `core/`           - Data structures and methods managing connection state from agents/clients
* `cryptography/`   - Cryptography code and wrappers around Go's standard crypto APIs
* `encoders/`       - Data encoders and decoders
* `generate/`       - Package for generating agents and shared libraries
* `gobfuscate/`     - Compile-time obfuscation library
* `gogo/`           - Go wrappers around the Go compiler toolchain
* `handlers/`       - Methods invokable by agents without user interaction
* `log/`            - Wrappers around Logrus logging
* `module/`         - Agents' functionality is used through modules. Their functionning is defined here.
* `msf/`            - Metasploit helper functions
* `rpc/`            - RPC functions and logic, called by clients to control agents
* `transport/`      - Contains server-side code for handling communication with clients
* `users /`         - Contains server-side code for managing Wiregost users 
* `website/`        - Static content that can be hosted on HTTP(S) C2 domains
