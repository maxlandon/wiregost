package main

// The core package is the backbone of the WireGost server, Spectre.

// It will interface with various components such as authentication services,
// listeners, databases, etc.

import (
	"fmt"

	"github.com/maxlandon/wiregost/internal/server/core"
)

type Spectre struct {
	// Data
	// Services
	// UserManager
}

func main() {

	err := core.NewClientRPC()
	if err != nil {
		fmt.Println(err)
	}
}
