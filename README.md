
#                      <center>WireGost - Golang Exploitation Framework</center> 
______

![Demo](./.github/images/console-greet.png)
![Sessions-Interact](./.github/images/sessions-interact.png)


The grounds for the WireGost exploitation framework project are:
* **Go is syntactically dead simple** and has a simple C-like object-model: This will lower the number of abstraction-layers-caused headaches.
* Consequently, Go code is easily readable and maintainable. Go is strongly typed, which make Gophers winners on all fronts.
* Go **compiles itself cross-platform**: Code compiled on a Linux machine will run _on virtually any architecture and operating system_.
  The implications for payload efficiency are wide-ranging.
* Go standard library includes what is probably the most advanced networking stack at the moment. Again, the implications for both framework
  complexity and payload modularity are significant.

**The name:** I recently looked a [video](https://www.youtube.com/watch?v=T8aXx3K_lKY) where the notorious 
[Kevin Mitnick](https://en.wikipedia.org/wiki/Kevin_Mitnick) is interviewed by an attractive journalist about the usual security 
and pricacy issues. Boring questions, boring answers (so boring that everyone comments on this instead of saying obsenities on 
the girl... very surprising), but he remembered me the name of his book: _"Ghost in the Wires"_. Just on point. Thank you Mr. Mitnick.

______
## Documentation 

The documentation for WireGost is also available on the [Wiki](https://github.com/maxlandon/wiregost/wiki) of this repository.
You will find everything needed to install, setup and use Wiregost C2 Server and Console.

______
## Sub-Repository Tools

#### EffectiveCouscous
[Maltego](https://www.paterva.com/web7/buy/maltego-clients/maltego-ce.php) is not really a tool for computer exploitation. 
Maltego is a software that allows to graph various kinds of networks (computers, social, criminal, and many others), in a
versatile, flexible, automated and efficient way, . It can be used for any activity having an investigative character. 
And since computer security is the cute child of a chessboard, a magic labyrinth and an escape game, Maltego is your best friend.
[EffectiveCouscous](https://github.com/maxlandon/EffectiveCouscous) is an attempt at interfacing various security tools with Maltego.

Maltego will act as a GUI interface with visualization, inference and discovery capabilities. All data from WireGost will be used
by Maltego.


## Projects that have inspired/motivated/been outright copy-pasted

#### Sliver
[Sliver](https://github.com/BishopFox/sliver) is a post-exploitation/implant framework written in Go. It seems to be the most advanced 
framework written in Go at the moment. Therefore, most of WireGost codebase is exactly the same than Sliver. We then need to address
**a huge thanks** and **deep and sincere excuses** to the BishopFox team, because my code is mostly theirs, and I have shamelessly change the "Slivers"
everywhere with "Ghost". (I have a good excuse, though: it was the most efficient way to force myself going everywhere in their code base.)
Thanks a lot, because I've learned a ton of things from it, and I'm really admirative of such tools, I would be totally unable to produce a iota
of it on my own. 

#### Merlin
[Merlin](https://github.com/Ne0nd0g/merlin) is also a post exploitation framework written in Go. It emphasizes on the use of HTTP/2 for C2
communications. It also includes a Javascript agent, post-exploitation modules (mostly in PowerShell) usable **a-la-metasploit**. Downsides are
only one server capability, and no multi-client capacity either.

#### All the others

Computer security is as large a subject as computers alone. It goes the same for the number of tools related to it.
I would gladly pay for another 30 lives so I can discover them all, but I don't have God's SWIFT account number, and again, I'm
poor as hell. If, in the context of this project, some of them are worth so much that it would be criminal not to include 
them in this list, I will add them.


______
## TO DO & ROADMAP

#### TO DO

**Console**
* Check codebase linting
* Add completers for:
    - Session names for interacting
    - Workspace names for setting implant options
    - All command options/filters in the implant menu
    - Filesystem completion in the implant menu
    - Fix the completion for help commands, depending on menu context
* Commands for:
    - Deleting generated ghost implants
    - Deleting profiles
    - Deleting users
* Less hacky option filters for many commands, and better command help for these
* Add/Rewrite help for:
    - Execute-Assembly command
    - Add examples to many command helps
* Fix commands:
    - Fix resource make (number of lines saved)
    - Fix resource load (which is not refreshing the shell state/context)
    - Fix `stack pop` which is not pulling the next module from the server.
* Config file/content for implant prompt/completions/etc...

**C2 Server**
* Persistent module stacks
* Persistent listeners
* Fix connect/disconnect detections from the server
* Add workspace/host settings to implant modules + implant registration
* Fix .pentest/path for Data Service Env loading
* Help for MSF listeners / eventually a separate module.
* Check all proc/priv/execute commands.
* Check why obfuscated implants cannot be generated at the same time without messing the namespace up

**Data Service**
* Change Certificates location/use/storage, etc... (Potentially merge with the Server) + code to handle this.
* Move config file to the Server config directory
* Try to make the Server not depending too much on the Data Service (if possible)

**Documentation**
* Add PostgreSQL install/setup to Required Dependencies
* Pages for:
    - Canaries commands
    - Websites commands
    - Ghosts commands
    - Implant config
    - Console troubleshooting
    - Priv Commands
    - Proc Commands
    - Execute Commands
    - Agent shell command
    - Post Modules

Writing Modules:
    - Modules Overview (Payload & Post)

Data Service:
    - Config
    - Usage
    - Systemd
    - Host Commands

**Code Repository**
* Update all READMEs:
    - If they miss files in their lists
    - If they are not accurate


______
## Warmest Thanks
* The **Golang Project**.
* **BishopFox** for their Sliver framework, with which I've learned a lot.
* The **Merlin** project, with which I learned a lot too !
