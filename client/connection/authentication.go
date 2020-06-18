package connection

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	cliCtx "github.com/maxlandon/wiregost/client/context"
	client "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
)

// Authenticate - Perform full authentication process with server
func Authenticate(conn *grpc.ClientConn) (cli client.ConnectionRPCClient, user dbpb.User) {

	// Register ConnectionRPC client to connection
	cli = client.NewConnectionRPCClient(conn)

	// Send authentication request
	var counter int

	// Loop for several attempts
	for {
		// User has 5 allowed attempts to authenticate
		if counter < 5 {
			// Prompt, store and send password (as a hash)
			req := &client.AuthenticationRequest{}
			req.Password = PromptUserPassword()
			req.Username = assets.ServerUser
			req.MD = cliCtx.SetMetadata()

			// Send request
			res, err := cli.Authenticate(context.Background(), req, grpc.EmptyCallOption{})

			// If refused, try again (five tries)
			if res.Success == false {
				fmt.Println(tui.Red("Wrong password."))
				counter++
				continue
			}

			// If error, notify conn error and exit application
			if err != nil {
				fmt.Println(tui.Red("Error during authentication request."))
			}

			return
		}

		// If we go here, then user has tried five times unsuccessfully.
		fmt.Println(tui.Red("Authentication failure: 5 wrong attempts."))
		fmt.Println(tui.Red("Exiting application"))
		os.Exit(1)
	}
}

// PromptUserPassword - Ask the console user to authenticate
func PromptUserPassword() (password string) {

	fmt.Println("Enter user password:")
	pass, err := terminal.ReadPassword(1)
	if err != nil {
		return
	}
	password = string(pass)

	return
}

var (
	Info    = fmt.Sprintf("%s[-]%s ", tui.BLUE, tui.RESET)   // Info - All normal messages
	Warn    = fmt.Sprintf("%s[!]%s ", tui.YELLOW, tui.RESET) // Warn - Errors in parameters, notifiable events in modules/sessions
	Error   = fmt.Sprintf("%s[!]%s ", tui.RED, tui.RESET)    // Error - Error in commands, filters, modules and implants.
	Success = fmt.Sprintf("%s[*]%s ", tui.GREEN, tui.RESET)  // Success - Success events

	Infof   = fmt.Sprintf("%s[-] ", tui.BLUE)   // Infof - formatted
	Warnf   = fmt.Sprintf("%s[!] ", tui.YELLOW) // Warnf - formatted
	Errorf  = fmt.Sprintf("%s[!] ", tui.RED)    // Errorf - formatted
	Sucessf = fmt.Sprintf("%s[*] ", tui.GREEN)  // Sucessf - formatted

	RPCError     = fmt.Sprintf("%s[RPC Error]%s ", tui.RED, tui.RESET)     // RPCError - Errors from the server
	CommandError = fmt.Sprintf("%s[Command Error]%s ", tui.RED, tui.RESET) // CommandError - Command input error
	ParserError  = fmt.Sprintf("%s[Parser Error]%s ", tui.RED, tui.RESET)  // ParserError - Failed to parse some tokens in the input
	DBError      = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)      // DBError - Data Service error
)
