
#                      <center>WireGost - Golang Exploitation Framework</center> 
______

![Demo](./pics/welcomeToWireGost.png)


The grounds for the WireGost project are:
* **Go is syntactically dead simple** and has a balanced object-model: the functional logic is no longer separated from code logic.
  This will lower the number of abstraction-layers-caused headaches.
* Consequently, Go code is highly readable and maintainable. Go is strongly typed, which make Gophers winners on all fronts.
* Go **compiles itself cross-platform**: Code compiled on a Linux machine will run _on virtually any architecture and operating system_.
  The implications for payload efficiency are wide-ranging.
* Go standard library includes what is probably the most advanced networking stack at the moment. Again, the implications for both framework
  complexity and payload modularity are huge.
* Go, as a result of the above and for other reasons, acts pretty well as a system programming language. 

These grounds seem to be even more solid when considering another widely-known tool:
* The [Metasploit Framework](https://github.com/rapid7/metasploit-framework) is wonderful, by many regards, but...
* Ruby has a cryptic object model, which somehow becomes a nightmare when the need arises to write an exploit or its related libraries.
  And if the exploit touches upon some other complex language like C++, you now have two nightmares, instead of one. Go will not solve
  the problem entirely but 1) some strong typing won't be a luxury, 2) being much simpler, it will relieve the user from abstraction-caused 
  headakes.

**The name:** I recently looked a [video](https://www.youtube.com/watch?v=T8aXx3K_lKY) where the notorious 
[Kevin Mitnick](https://en.wikipedia.org/wiki/Kevin_Mitnick) is interviewed by an attractive journalist about the usual security 
and pricacy issues. Boring questions, boring answers (so boring that everyone comments on this instead of saying obsenities on 
the girl... very surprising), but he remembered me the name of his book: _"Ghost in the Wires"_. Just on point. Thank you Mr. Mitnick.

______
## Table of Contents 

The documentation for WireGost is also available on the [Wiki](https://github.com/maxlandon/wiregost/wiki) of this repository.

It is **strongly advised to fully read the Overview section first**, as it explains the structure
and functionning of WireGost, and as this will be important in order to understand the installation
process. Other than this, everything should be straightforward.

### Overview
* [**General Architecture**](https://github.com/maxlandon/wiregost/wiki/General-Architecture)
* [**Code Structure**](https://github.com/maxlandon/wiregost/wiki/Code-Structure)
* [**Environment**](https://github.com/maxlandon/wiregost/wiki/Personal-Environment)


### Usage
* [**Installation**](https://github.com/maxlandon/wiregost/wiki/Installation)
* [**Base Usage**](https://github.com/maxlandon/wiregost/wiki/Base-Usage)
* [**Commands**](https://github.com/maxlandon/wiregost/wiki/Commands)
    * [Core](https://github.com/maxlandon/wiregost/wiki/Core-Commands)
    * [Help](https://github.com/maxlandon/wiregost/wiki/Help-Commands)
    * [Server](https://github.com/maxlandon/wiregost/wiki/Server-Commands)
    * [Log](https://github.com/maxlandon/wiregost/wiki/Log-Commands)
    * [Chat](https://github.com/maxlandon/wiregost/wiki/Chat-Commands)
    * [Workspace](https://github.com/maxlandon/wiregost/wiki/Workspace-Commands)
    * [Stack](https://github.com/maxlandon/wiregost/wiki/Stack-Commands)


### Other 
* [**Required Specs**](https://github.com/maxlandon/wiregost/wiki/Required-Specs)
* [**Ideas**](https://github.com/maxlandon/wiregost/wiki/Ideas)
* [**To Do**](https://github.com/maxlandon/wiregost/wiki/To-Do)

______
## Current Landscape and Tools

As always when devising security-oriented tools, a little list of base tools is needed.
Especially when your practical experience of computer security is limited, as is mine.
This list might be somehow temporary, but serves more as a source of inspiration and ideas.

#### Metasploit

The [Metasploit Framework](https://github.com/rapid7/metasploit-framework) is a wonderful framework. It is even more than that, it 
is a teacher. It is so good a teacher that all the planet knows its face, the face of its -meterpreter- kids, and the look of its 
classroom. It cannot go anywhere without raising the (AV) crowds. That's nice when you are a politician, but that's bad when 
your from intelligence.

Metasploit emphasizes modularity: modularity of exploits, of payloads, of modules, of protocols, and so on.
As far as this modularity is concerned, Metasploit needs to be found back in WireGost. It is a requirement.
Also, and since almost a year, Metasploit makes use of a HTTP service for interacting with its database.
This is obviously very useful, especially when considering a GUI interface making use of this data. See **Maltego** below for 
ideas about such a GUI interface.

#### Sliver

[Sliver](https://github.com/BishopFox/sliver) is an post-exploitation/implant framework written in Go. It seems to be the most advanced 
framework written in Go at the moment. Therefore, most of WireGost codebase will be inspired of Sliver (if not outright copied from it).

#### Cobalt Strike, Canvas

I actually have nothing to say on these ones: I'm poor as hell.

#### PowerShell Empire

[PowerShell Empire](https://www.powershellempire.com) is a rather new Post-Exploitation Framework built, as its name tells, in 
PowerShell and Python. It is dedicated to post-exploitation on Windows and Apple hosts. It introduces various things that do not 
exist in Metasploit Meterpreter agents:

* An asynchronous communication model (Meterpreters are _partially HTTP-based_ and synchronous, 
  which means the connection is constantly opened. Empire agents are _fully HTTP based_, which means 
  the connection is only opened for the time of sending the message.)
* Renewed Mimikatz modules.
* Improved and enlarged persistence capabilities.
* Some kill-switch controls over non-stealthy modules. (Opsec safe/unsafe)
* A bunch of other modules.

#### Maltego

[Maltego](https://www.paterva.com/web7/buy/maltego-clients/maltego-ce.php) is not really a tool for computer exploitation. 
Maltego is a software that allows to graph various kinds of networks (computers, social, criminal, and many others), in a
versatile, flexible, automated and efficient way, . It can be used for any activity having an investigative character. 
And since computer security is the cute child of a chessboard, a magic labyrinth and an escape game, Maltego is your best friend.

Maltego will act as a GUI interface with visualization, inference and discovery capabilities. All data from WireGost will be used
by Maltego.

#### All the others

Computer security is as large a subject as computers alone. It goes the same for the number of tools related to it.
I would gladly pay for another 30 lives so I can discover them all, but I don't have God's SWIFT account number, and again, I'm
poor as hell. If, in the context of this project, some of them are worth so much that it would be criminal not to include 
them in this list, I will add them.

______

## The Requirements

[This is the list of requirements for WireGost](https://github.com/maxlandon/wiregost/wiki/Requirements). 
This list will be updated as ideas appear, appear to be good, or appear to be bad.

______

## Warmest Thanks
* To the **Metasploit** team and all of the contributors to the project, thanks to whom I -and so many others- have learned countless
  things about computers, before learning about security.
* The same applies to contributors of **PowerShell Empire**, **Maltego**, and so many other tools.
* The **Golang Project**.
* The thousands of people who have contributed and still do, for free, to help us being a little less ignorant and continuously
  more curious.
