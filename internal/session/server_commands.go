package session

import "github.com/chzyer/readline"

// This file contains all server command handlers and their registering function.

// Register User
func (s *Session) registerUserHandler(args []string, sess *Session) error {
	s.ServerManager.RegisterUserToServer(s.User)
	return nil
}

// Register all handlers defined above
func (s *Session) registerServerHandlers() {
	//Register User
	s.addHandler(NewCommandHandler("server.register",
		"server.register",
		"Register user to server, during first-time connection",
		s.registerUserHandler),
		readline.PcItem("server.register"))
}
