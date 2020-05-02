
## Wiregost Database System 

This directory contains all code needed by consoles, the server and implants for storing data and interacting with it.
At its core, the Database system in Wiregost relies on PostgreSQL (thus, Wiregost data model is an ORM), and gRPC/Rest
for communication.

The aim of Wiregost is to provide a simple and extensible data model. When designing it, I relied heavily upon Metasploit
data model (useful for things like Credentials), the Nmap object model (the output of an XML nmap scan), and integrated
both within Wiregost. It is then enhanced by adding other entities needed for the core functionality (listeners, implants, etc...)

However, the complexity of an object model oriented toward security discipline is non-negligible nowadays. Meanwhile, it must remain
simple enough so that exploitation/post/research/any software can contribute to and interact with it easily.


----
### Complexity of queries and Data structuring efficency

Because Wiregost's Data is meant to be a service shared with tools than Wiregost, the Data Model has to be structured enough so that
each tool can:
- Search for entities easily
- Add new entities that are not redundant
- Update them ONLY if we know *whether we need to add OR update*.
- Delete them only if we clean all associated data

Achieving these four goals at the same time involves that we devise a structured model first, and then that we expose only a few methods
per entity (ex: a host), like get/add/update/delete. Therefore, Wiregost's Database system is a repository that is supposed to store data
in a smart way.


----
### Directory Content 

- `data_service.go`     - Entrypoint for the Dabase functionality. Functions in this file are called by the server.
- `assets.go`           - Code for managing the DB directory, its configuration file, etc...
- `schema.go`           - Automatically migrates all objects available to the PostgreSQL DB with Go types.
- `models/`             - All query functions, with a file for each object. Clients don't use them
- `remote/`             - The client functions, those actually called by server and consoles.

