<!-------------------------------------------------------->
______
## Database

**Database**
* `db.source.connect {source_name}`     _Connect to data source_
* `db.source.show`                      _Show all data sources_
* `db.source.add`                       _Add data source based on current db.source parameters_
* `db.source.delete {source_name}`      _Delete a data source_ 
* `db.file.export`                      _Generate export file_
* `db.file.import`                      _Import Data file_

_Parameters_
* `db.file.export.format`       _Specify output format_
* `db.file.export.filename`     _Specify filename_
* `db.file.import.filepath`     _Specify import file path_

* `db.source.url`
* `db.source.certificate`
* `db.source.priv`
* `db.source.password`

<!-------------------------------------------------------->
______
## Entities

**Hosts**
* `hosts show all`              _Various search filters for hosts_
* `hosts show ip`
* `hosts show mac`
* `hosts show os`
* `hosts show name`
* `hosts show purpose`
* `hosts show filter`

**Services**
* `services show all`           _Various search filters for services_
* `services show port`
* `services show proto`
* `services show name`
* `services show info`
* `services show state`
* `services show filter`

**Creds**
* `creds show all`              _Various search filters for creds_
* `creds show type`
* `creds show priv`
* `creds show filter`

**Notes**
* `notes show all`              _Various search filters for creds_
* `notes show host`
* `notes show filter`

_And other filters_


<!-------------------------------------------------------->
______
## Listeners 

 **Listeners** show all listeners launched from modules in the session
* `listeners.show {all, listener_name}`   _Show one or more listeners_
* `listeners.kill (all, listener_name)`   _Kill some or all listeners_
* `listeners.rename current new`         _Rename listener_
* `listeners.duplicate {listener_name}`   _Duplicate listener and launch it_


<!-------------------------------------------------------->
______
## Variables

**Note**: The "advanced" option classification, inspired from Metasploit, might only be
useful as a display filter/structure, not as a command filter. Maybe these commands will
be removed

**Variables**
* `global.show`                                                 _Show all global variables_
* `global.options.{all_known_modules_options} (value)`          _Set global option_
* `global.options.{all_known_modules_options} unset `           _Unset global option_


<!-------------------------------------------------------->
______
## Modules

**Edit**
* `edit`                         _Edit current module_
* `loadpath`                     _Load code path_
* `reload (all|lib|module)`      _reload all modules_


**Exploits**
* `use multi/windows/java`                      _Load a module and make it active_
* `exploit.options.show`                        _Show module options_
* `exploit.advanced`                            _Show module options_
* `exploit.payloads.show`                       _Use compatible payload_ 
* `exploit.payloads.set path/to/payload`        _Use compatible payload_ 
* `exploit.encoders.show`                       _Use compatible encoder_
* `exploit.encoders.set path/to/encoder`        _Use compatible encoder_
* `exploit.run`                                 _Run exploit_

_Parameters_
* `exploit.options.{exploit_generic} value`     _Set module options_
* `exploit.advanced.{exploit_specific}`         _Set advanced options_

_Examples_
* `exploit.options.URL`					        _Set URL path_
* `exploit.options.target (target_list)`        _Set target type_
* `exploit.options.LHOST`                       _Set local host_
* `exploit.options.RHOST`                       _Set remote host_

**Payloads**
* `payload.load windows/x64/shell/reverse_https`		
* `payload.handler`
* `payload.generate`

_Parameters_
* `payload.options.{payload_generic}`
* `payload.advanced.{payload_specific}`

_Goes on for evasion, post, auxiliary_
