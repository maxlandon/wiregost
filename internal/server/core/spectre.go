package core

// The core package is the backbone of the WireGost server, Spectre.

// It will interface with various components such as authentication services,
// listeners, databases, etc.

type Spectre struct {
	// Database
	DBManager *DBManager
	// Services
	UserManager *UserManager
	// RPC Servers
	ClientRPC *ClientRPC
}

func NewSpectre() *Spectre {
	spectre := &Spectre{
		DBManager:   NewDBManager(),
		UserManager: NewUserManager(),
		ClientRPC:   NewClientRPC(),
	}

	// Make DBManager available to all services needing it.
	spectre.RegisterDatabaseToServices()

	// Register all services to the ClientRPC.
	spectre.RegisterServicesToRPC()

	// Start RPC
	spectre.ClientRPC.Start()

	return spectre
}

func (s *Spectre) RegisterDatabaseToServices() error {
	// UserManager
	s.UserManager.database = s.DBManager

	return nil
}

func (s *Spectre) RegisterServicesToRPC() error {
	RegisterUserManagerServer(s.ClientRPC.server, s.UserManager)

	return nil
}
