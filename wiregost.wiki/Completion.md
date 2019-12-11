
Command completion should take a prominent place into WireGost, as many types of input are boring
and tedious to perform, such as IP addresses, MAC, hostnames, and so many others.
As well, command completion is a wonderful way to provide, simultaneously:
* Info about entities (current state, etc)
* Help to commands
* Help, current and default value to parameters (global, or module-specific)


Module command completion in Metasploit has a flaw: it only outputs _all_ modules satifying the
command pattern. Which is way too much, because you don't want to have 1800 exploits, 400 payloads
or 1000 auxiliary modules displayed. Module command completion should be changed to just show, for
instance, the options until the next `/`.
* This can be solved by using go-prompt completion, with a path completer for exploits, payloads, etc...


