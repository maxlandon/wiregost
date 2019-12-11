package main

// This file is the WieGost client shell executable.

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/maxlandon/wiregost/internal/session"

	"github.com/evilsocket/islazy/log"
)

func main() {
	sess, err := session.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sess.Close()

	if err = sess.Start(); err != nil {
		log.Fatal("%s", err)
	}

	// Eventually start the interactive session.
	for sess.Active {
		line, err := sess.ReadLine()
		if err != nil {
			if err == io.EOF || err.Error() == "Interrupt" {
				if exitPrompt() {
					sess.Run("exit")
					os.Exit(0)
				}
				continue
			} else {
				log.Fatal("%s", err)
			}
		}

		for _, cmd := range session.ParseCommands(line) {
			if err = sess.Run(cmd); err != nil {
				log.Error("%s", err)
			}
		}
	}
}

func exitPrompt() bool {
	var ans string
	fmt.Printf("Are you sure you want to quit this session? y/n ")
	fmt.Scan(&ans)

	return strings.ToLower(ans) == "y"
}
