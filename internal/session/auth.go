package session

import (
	// Standard
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	// 3rd party
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"golang.org/x/crypto/ssh/terminal"
)

// User is used to store per-shell user credentials, which will
// be used to authenticate to a server when connecting/sending requests.
type User struct {
	Name               string
	PasswordHashString string
	PasswordHash       [32]byte
	CredsFile          string
}

// NewUser instantiates a new User.
func NewUser() *User {
	user := &User{CredsFile: "~/.wiregost/client/.auth"}
	return user
}

// LoadCreds loads user credentials from a configuration file.
func (user *User) LoadCreds() (err error) {
	// Check for personal directory, exit if not present.
	credsFile, _ := fs.Expand(user.CredsFile)
	if fs.Exists(credsFile) == false {
		fmt.Println(tui.Red(" ERROR: No ID and authentication information found."))
		fmt.Println(tui.Red("        Please run the ghost_setup.go script (in the " +
			"scripts directory), for initializing and configuring the client first"))
		os.Exit(1)
	} else {
		// Load authentication parameters
		credsFile, _ := fs.Expand(user.CredsFile)
		configBlob, _ := ioutil.ReadFile(credsFile)
		json.Unmarshal(configBlob, &user)
	}
	return err
}

// Authenticate is used to perform local (shell only) authentication.
func (user *User) Authenticate() {
	attempts := 0
	fmt.Printf(tui.Bold("Password: \n"))
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	hash := sha256.Sum256(pass)
	for {
		// Success, authenticate
		if bytes.Equal(hash[:], user.PasswordHash[:]) {
			fmt.Println(tui.Green("Authenticated"))
			break
		}
		// Failure, 3 chances and then exit
		if !bytes.Equal(hash[:], user.PasswordHash[:]) {
			if attempts <= 3 {
				fmt.Println("Wrong password. Retry:")
				pass, _ = terminal.ReadPassword(int(syscall.Stdin))
				hash = sha256.Sum256(pass)
				attempts++
			}
			if attempts == 3 {
				fmt.Println(tui.Red("Authentication failure. Leaving WireGost"))
				os.Exit(1)
			}
		}
	}
}
