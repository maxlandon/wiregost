## Console

The `console` package contains the code for the console 'per se'.

* `console.go`      - Console instantiation and setup (using chzyer/readline package)
* `config.go`       - Console configuration loading/saving/setup 
* `prompt.go`       - Code for dynamic display/refresh of the prompt, with variable substitution.
* `exec.go`         - Dispatch and execution of commands.
* `connect.go`      - Function for connecting to server and binding commands/events.
* `context.go`      - Initialize shell context/variables and makes them available to commands.
* `events.go`       - Handles events coming from the server.
