package main

import (
	"github.com/maxlandon/wiregost/internal/cli"
)

func main() {
	// In Merlin,
	// go cli.Shell()       Ending right away. Why ?
	cli.NewSession()
	// session := cli.NewSession()
	// session.Shell()
}
