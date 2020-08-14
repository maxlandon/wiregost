package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/evilsocket/islazy/tui"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

// Shell - The Shell object implements an interactive Command Shell around a stream.
// This object is meant to provide basic functionality for manipulating a remote
// shell session such as a reverse /bin/sh one-liner: this remote stream will only
// provide primitive I/O and command execution interface. This object wraps it around
// a few methods and attempts to transform it into a more usable command-line interface.
type Shell struct {
	// Interactive - Provides core methods to manipulate a stream
	*Interactive

	// tokenIndex - Used in metasploit to control if the token used to delimit command
	// outputs has to be run each time. Token set is a checker to avoid useless rerun
	tokenIndex int
	tokenSet   bool

	// connBuffer - We need, for each command sent, to delimit its according output
	// without either blocking or without loosing bytes when reading the stream.
	// Accordingly, we use a Reader that can buffer data coming from conn, and we
	// can work on it peacefully.
	connBuffer *bufio.Reader
}

// NewShell - Instantiates a new Session that implements a command shell around a logical stream.
func NewShell(stream io.ReadWriteCloser) (sh *Shell) {

	sh = &Shell{
		NewInteractive(stream),  // The session is interactive.
		0,                       // The token is by default 0. Will check in future if needs change.
		false,                   // Not set yet, will be once only.
		bufio.NewReader(stream), // We are ready to buffer and manipulate the stream.
	}

	sh.Type = serverpb.SessionType_SHELL

	sh.Log = sh.Log.WithField("type", "shell") // Log settings

	return
}

// SetupLog - A Shell session sets up a few things in the log: it logs all commands
// sent to the remote endpoint, as well as their output. This can be used as a buffer,
// in case user wants to go back in history.
func (sh *Shell) SetupLog() (err error) {
	return
}

// RunCommand - Given a command and its arguments, with an optional timeout, send it through the
// shell stream, and wait for response. This should be synchronous given that we cannot run
// concurrent operations on a remote shell while separating their respective stdin/stdout.
func (sh *Shell) RunCommand(line string, timeout time.Duration) (result string, err error) {

	// 1st and foremost, check for user input, if there are some escape codes we
	// need to handle.

	// Here, ultimately we should have already saved a host to DB for this session
	// even if there is close to zero information about it other than suppositions
	// from parent exploits/transports/payloads...
	// Instead of having a loose sh.Platform string which could be anything, we
	// will think of db.Host next time we see this message.

	// Waiting, only Unix...
	return sh.runCommandUnix(line, timeout)
}

// LoadRemoteEnvironment - Primary method for retrieving the target environment variables, and assigning
// them to this session, therefore available to both console completions and further session code.
func (sh *Shell) LoadRemoteEnvironment(full bool) (err error) {

	// user

	// Environment variables

	// network interfaces

	// processes

	// Executables

	// FULL ------------------

	// users

	// groups

	// directory tree

	return
}

// write - Write to the connection with a timeout and stream error control.
func (sh *Shell) write(line string, timeout time.Duration) (err error) {

	// Channel controls
	doneWrite := make(chan struct{}, 1)
	errWrite := make(chan error, 1)

	// Asynchronous, timed writing
	go func() {
		ilength, err := sh.stream.Write([]byte(line))
		if ilength != len(line) {
			sh.Log.Warnf("length of bytes written to stream and returned"+
				" output length don't match: sendt:%d != returned:%d", len(line), ilength)
		} else if err != nil {
			errWrite <- err
		} else {
			doneWrite <- struct{}{}
		}
	}()

	// Handle errors and timeouts
	select {
	case err := <-errWrite:
		return err
	case <-time.After(timeout):
		return errors.New("write operation timed out.")
	case <-doneWrite:
		return nil
	}
}

// runCommandUnix - Adapt the token for Unix remote endpoints
func (sh *Shell) runCommandUnix(line string, timeout time.Duration) (out string, err error) {

	// Log command line
	sh.Log.Infof("%scommand%s: %s%s", tui.GREEN, tui.RESET, tui.BOLD, line)

	// Set token index if not set yet. Then create token
	if sh.tokenSet == false {
		sh.tokenSet, err = sh.setTokenIndex(timeout)
		if err != nil && !sh.tokenSet {
			return
		}
	}
	token := randStringBytesRmndr()

	// Construct the final cmd and send it.
	// Don't return from an error immediately:
	// We give the shell a chance to read output.
	line += fmt.Sprintf(";echo %s\n", token)
	err = sh.write(line, timeout)
	if err != nil {
		sh.Log.Errorf(err.Error())
	}

	return sh.readUntilToken(token, sh.tokenIndex, timeout)
}

// runCommandWindows - Run a command on a Windows endpoint and adapt the token.
func (sh *Shell) runCommandWindows(line string, timeout time.Duration) (out string, err error) {
	return
}

// readUntilToken - Accepts a token for delimiting output buffers. This is useful for
// sequencing printing on user consoles. Metasploit makes uses of tokens: they append
// to each string sent ';echo {token}' , for having retroactive output delimitation.
// We might also use an identified prompt string (for instance the first one received).
func (sh *Shell) readUntilToken(token string, wantedIndex int, timeout time.Duration) (out string, err error) {

	if timeout == 0 {
		sh.Log.Warnf("readUntilToken() was called with empty timeout: returning")
		return // Exit immediately, should have been set
	}

	// Set parts
	// var partsNeeded int
	// if wantedIndex == 0 {
	//         partsNeeded = 2
	// } else {
	//         partsNeeded = 1 + (wantedIndex * 2)
	// }

	// Read loop
	for {
		select {
		case <-time.After(timeout):
			err = errors.New("read operation timed out.")
			break // Give a shot to return the (incomplete) part and log it.
		default:
			// We read the incoming data: if we receive some and the token is not
			// found, we increment/reset the timeout: maybe the connection is slow.

			// line, err := .ReadString('\n')
			// if err != nil {
			//         return line, err
			// }
		}
	}

	// End: log output

}

// NOTE: There is often a lag between a command output and the echo of the cmd we sent.
// Therefore, according to Metasploit: "If the session echoes input we don't need to echo
// the token twice. This setting will persist for the duration of the session."
// We might recalculate from time to time if needed.
func (sh *Shell) setTokenIndex(timeout time.Duration) (set bool, err error) {

	sh.Log.Debug("Setting token index for shell session output buffer control")
	// Need to tokens to test
	token := randStringBytesRmndr()
	numeric_token := rand.Int63()

	// Need two commands, accordingly
	testCmd := fmt.Sprintf("echo %d", numeric_token)
	cmd := fmt.Sprintf(testCmd+";echo %d", token)

	// Send package
	err = sh.write(cmd, timeout)
	if err != nil {
		return
	}

	// Read response
	out, err := sh.readUntilToken(token, 0, timeout)
	if err != nil {
		return
	}

	// Check if response is indeed a number
	if nb, ok := strconv.Atoi(out); ok == nil {
		if nb != 0 && int64(nb) == numeric_token {
			sh.tokenIndex = 0
			sh.tokenSet = true // We don't need to do it anymore anyway.
			sh.Log.Infof("shell output token match with test token: %d == %d", numeric_token, nb)
		} else {
			sh.tokenIndex = 1
			sh.tokenSet = true
			sh.Log.Warnf("shell output token does not match with test token: %d == %d", numeric_token, nb)
		}
		return
	}

	return false, fmt.Errorf("could not convert token to int: %s", out)
}

// randStringBytesRmndr - Easily produce random tokens for output buffer control.
func randStringBytesRmndr() string {

	n := rand.NewSource(time.Now().Unix())

	b := make([]byte, n.Int63())
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
