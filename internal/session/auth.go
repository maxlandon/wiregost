package session

// This file contains all objects and functions related to client-side authentication,
// ONLY for normal shell usage, not at client initialization (first time shell is used).

//		- New Authentication information
//		- Loading Authentication information
//		- Sending Authentication information to the server.

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"golang.org/x/crypto/ssh/terminal"
)

type Auth struct {
	UserName            string
	PasswordHashString  string
	UserToken           string
	UserCertificatePath string
	PasswordHash        [32]byte
	UserAuthFile        string
}

func NewAuth() *Auth {
	auth := &Auth{
		// Default settings are loaded, will be overwritten if
		// config file is found during DoClientAuth()
		UserName:            "",
		PasswordHashString:  "",
		UserToken:           "",
		UserCertificatePath: "",
		UserAuthFile:        "~/.wiregost/.auth",
	}

	return auth
}

func (auth *Auth) LoadAuth(sess *Session) (err error) {

	// Check for personal directory, exit if not present.
	authFile, _ := fs.Expand(auth.UserAuthFile)
	if fs.Exists(authFile) == false {
		fmt.Println(tui.Red(" ERROR: No ID and authentication information found."))
		fmt.Println(tui.Red("        Please run the ghost_setup.go script (in the scripts directory), for initializing and configuring the client first"))
		os.Exit(1)
	} else {
		// Load authentication parameters
		fmt.Println(tui.Dim("Authentication parameters found."))
		path, _ := fs.Expand(auth.UserAuthFile)
		configBlob, _ := ioutil.ReadFile(path)
		json.Unmarshal(configBlob, &auth)
		fmt.Println(tui.Dim("Authentication file loaded."))
	}

	return err
}

// If User identifiers are already in config, prompt for password
func (auth *Auth) DoClientAuth(sess *Session) (result string) {
	fmt.Println()
	fmt.Printf(tui.Bold(tui.Yellow("Connection: "))+"%s\n", auth.UserName)
	attempts := 0

	fmt.Printf(tui.Bold("Password: \n"))
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	hash := sha256.Sum256(pass)

	for {
		// Success, authenticate
		if bytes.Equal(hash[:], auth.PasswordHash[:]) {
			fmt.Println(tui.Green("Authentication success"))
			fmt.Println()
			return "success"
		}
		// Failure, 3 chances and then exit
		if !bytes.Equal(hash[:], auth.PasswordHash[:]) {
			if attempts <= 3 {
				fmt.Println("Wrong password. Retry:")
				pass, _ = terminal.ReadPassword(int(syscall.Stdin))
				hash = sha256.Sum256(pass)
				attempts += 1
			}
			if attempts == 3 {
				fmt.Println(tui.Red("Authentication failure. Leaving WireGost"))
				return "failure"
			}
		}
	}

	return "success"
}
