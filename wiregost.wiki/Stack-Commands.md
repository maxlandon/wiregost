The Module Stack serves as a sort of "draft table" in which are loaded some modules (exploit, payload,
auxiliary, etc...)

The interest of the Stack, as previously explained, is to persist the state of the user's drafts (modules settings
already entered, etc...).

Then, when a user switches to a given workspace, the corresponding Stack, saved by the Server, is brought along and
loaded. This allows the user to find back what he was working on for this specific workspace.

The commands available are the following:

* `stack.show {all, module_name}`         _Show all modules that have been loaded since session start_
* `stack.unload {all, module_type}`       _Unload some or all modules currently loaded_
