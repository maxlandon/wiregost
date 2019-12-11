As explained in other parts of the Documentation, a workspace, in WireGost, has an important role
in saving and processing state for client. This state includes a Module Stack (specific to each workspace,
for each client), the set of entities (hosts, creds, services, etc) that will be available to requests, etc...

The commands and parameters available for workspace management are the following:

#### Commands

* `workspace.show`                          _Show all workspaces_
* `workspace.switch (name)`                 _Switch to workspace_
* `workspace.delete (name)`                 _Delete workspace_

#### Parameters

* `workspace.{options} set value`           _Set workspace options_


