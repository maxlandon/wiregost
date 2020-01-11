### Ghost

The `ghost/` directory contains all the code used by various agents in Wiregost.
It includes os/arch specific code needed by agents.
It also includes code needed for forwarding requests, if routes are set between agents.

* `constants/`      - Various shared constant values
* `evasion/`        - Evasion code used by agents
* `handlers/`       - Handlers used to receive commands from Wiregost
* `limits/`         - Code used to assess execution limitations (privileges, etc)
* `priv/`           - Code used for windows agents to check/manipulate privileges
* `procdump/`       - Used to dump processes on targets
* `proxy/`          - Used to detect/set/use proxies
* `forwarding/`     - Port forwarding functionality of agents
* `ps/`             - Find processes on targets
* `shell/`          - Used to spawn and interact with shell on target
* `syscalls/`       - Low-level code related to system APIs
* `taskrunner/`     - Manipulates processes/task (ex: inject shellcode)
* `transports/`     - Wires the agent to Wiregost server, also handles crypto
* `version/`        - Gather OS version/information of targets
* `winhttp/`        - Windows HTTP server related code

