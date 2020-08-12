package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/evilsocket/islazy/tui"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"

	"github.com/maxlandon/wiregost/client/assets"
	cliCtx "github.com/maxlandon/wiregost/client/context"
	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
)

// Authenticate - Perform full authentication process with server
func Authenticate(conn *grpc.ClientConn) (cli clientpb.ConnectionRPCClient, client *clientpb.Client) {

	// Register ConnectionRPC client to connection
	ConnectionRPC = clientpb.NewConnectionRPCClient(conn)
	cli = ConnectionRPC
	md := cliCtx.SetMetadata()

	// Send authentication request (loop 5 several attempts)
	var counter int
	for {
		if counter < 5 {
			// Prompt, store and send password (as a hash)
			req := &clientpb.AuthenticationRequest{}
			req.Username = assets.ServerUser
			req.Password = PromptUserPassword()
			req.MD = &md

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

			return cli, res.Client
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
