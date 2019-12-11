
WireGost partly relies on the user's personal environment, for:
* User-specific information and configuration.
* Loading and using personal exploit/payload modules.
* Using global variables.
* Storing loot, notes, and other types of data.


The personal environment is, by default, stored in the `~/.wiregost/` directory.
The structure inside is the following:
* `/.global_variables` (stores all global variables used and set in WireGost)
* `/.history` (stores the command history)

* `/logs/`  (All logs from WireGost)
    * `/global.log`     _(aggregates all logs from WireGost)_
    * `/exploit.log`    _(exploit-related logs)_
    * `/payload.log`    _(payload-related logs)_
    * `/listeners.log`  _(listeners-specific logs)_
    * `/db.log`         _(all data service related logs)_

* `/resource/`  (All resource files that can be loaded into WireGost)
    * `example.rc` _An example resource file_

* `/modules/`   (All personal modules)
    * `/exploit/` _Exploit module directory_
    * `/payload/` _Payload module directory_

