package mainmenu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Exit struct{}

// Execute - Run
func (e *Exit) Execute(args []string) (err error) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (Y/y): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}

	fmt.Println()
	return
}
