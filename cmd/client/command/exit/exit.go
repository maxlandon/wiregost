package exit

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/internal/client/console"
)

// ExitCmd - Exit the console
func ExitCmd(cmd *cobra.Command, con *console.Client, args []string) {
	con.Println("Exiting...")
	os.Exit(0)
}

// Commands returns the `exit` command.
func Command(con *console.Client) []*cobra.Command {
	return []*cobra.Command{{
		Use:     "exit",
		Short:   "Exit the program",
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			ExitCmd(cmd, con, args)
		},
	}}
}
