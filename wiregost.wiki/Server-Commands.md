As explained in other parts of the Documentation, WireGost works according to a client-server model.
Therefore, some commands and parameters are specific to Server control and|or information.

These commands and parameters are:

#### Commands

* `server.show`                         _Show client-side and server-side General & Authentication parameters_ 
* `server.generate_token`               _Ask the server to generate a new client-side token for having access to the server' services (commands, db)_
* `server.connect`                      _Connect to one of the saved WireGost servers_
* `server.add`                          _Add server based on the current value of parameters_

#### Parameters

* `server.address`                      _IP or resolved address of the server_
* `server.port`                         _Port on which the server is listening_
* `server.name`                         _Name under which this server will be saved and displayed_
* `server.default`                      _Make this server the default server on which to attempt connection when client is launched_
* `server.certificate`                  _Path to server certificate needed for connection (not required)_


Admin commands are only available in the client, if the client's user is registered in the WireGost server's administrators.

* `server.admin.show`                   _Show registered users (and show if they are active)_
* `server.admin.add_user (name)`        _Register a new user (password will be sent during first connection)_
* `server.admin.delete_user (name)`     _Delete one or more of the registered users_


