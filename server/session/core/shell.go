package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
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

	tokenIndex   int  // influences read loop on the stream and delimitations
	tokenSet     bool // indicates we have set the index by echoing a test token.
	pendingToken bool // used by the shell/builder for output concatenation.

	// builder - Used to construct the output of the Shell stream.
	// We clear delimiting tokens, command echoes, prompts, etc.
	builder *strings.Builder
}

// NewShell - Instantiates a new Session that implements a command shell around a logical stream.
func NewShell(stream io.ReadWriteCloser) (sh *Shell) {

	sh = &Shell{
		NewInteractive(stream), // The session is interactive.
		0,                      // The token is by default 0. Will check in future if needs change.
		false,                  // Not set yet, will be once only.
		false,                  // No tokens read yet
		&strings.Builder{},     // Used to build the output from conn
	}

	sh.Type = serverpb.SessionType_SHELL

	sh.Log.Logger.Out = os.Stdout
	sh.Log = sh.Log.WithField("type", "shell") // Log settings

	return
}

// RunCommand - Given a command and its arguments, with an optional timeout, send it through the
// shell stream, and wait for response. This should be synchronous given that we cannot run
// concurrent operations on a remote shell while separating their respective stdin/stdout.
func (sh *Shell) RunCommand(line string, timeout time.Duration) (result string, err error) {

	// Set token if needed. Persist for the entire session.
	if !sh.tokenSet {
		err = sh.setTokenIndex(timeout)
		if err != nil {
			sh.Log.Error(err)
			return
		}
	}

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

// runCommandUnix - Adapt the token for Unix remote endpoints
func (sh *Shell) runCommandUnix(line string, timeout time.Duration) (out string, err error) {

	// Log command line
	sh.Log.Infof("%scommand%s: %s%s", tui.GREEN, tui.RESET, tui.BOLD, line)

	// Token identifying this command
	token := randStringBytesRmndr()

	// Channel used by the command builder to notify timeout control
	done := make(chan struct{})

	// 1. Send command liner
	if line == "" || line == " " || line == "\n" {
		// First we check for empty commands, for checking the shell still responds
		// When the user has entered an empty input line, we use this as an automatic
		// and handy refreshment: we briefly display a spinner and then reprint prompt.
		sh.stream.Write([]byte(fmt.Sprintf("echo %s\n", token)))
	} else {
		// If not, we send the command with its token.
		sh.stream.Write([]byte(line + fmt.Sprintf(";echo %s\n", token)))
		sh.Log.Debugf("Done writing line: " + line + fmt.Sprintf(";echo %s\n", token))
	}

	// 2. Read connection.
	go func() {
		defer close(done)
		for {
			connBuf := make([]byte, 4096) // Big buffer, just in case.
			_, err = sh.stream.Read(connBuf)
			if err != nil {
				sh.Log.Error(err)
				return
			}
			// Process the output
			result, complete := sh.readUntilToken(string(connBuf), line, token, sh.tokenIndex)
			if complete {
				out = result
				break
			}
		}

	}()

	select {
	case <-done:
		return out, nil
	case <-time.After(timeout):
		close(done)
		// We still give out, in case it has something in it still.
		return out, fmt.Errorf("reading command result from conn stream timed out.")
	}
}

// runCommandWindows - Run a command on a Windows endpoint and adapt the token.
func (sh *Shell) runCommandWindows(line string, timeout time.Duration) (out string, err error) {
	return
}

// readUntilToken - The Shell's string builder is being passed temporary connection buffers
// and it processes them: finds command tokens to delimit output, trims command echos, prompts, etc.
func (sh *Shell) readUntilToken(connBuf, line, token string, tokenIndex int) (result string, complete bool) {

	// This variable is a buffer used by the builder.
	// For each token processed the builder replaces this.
	var cleared string

	// Process incomming connBuff bytes, which often contain trailing null ones.
	connBuf = string(bytes.Trim([]byte(connBuf), "\x00")) // Trim null bytes

	// 1) Check for token
	if findToken([]byte(connBuf), token) {
		// 2) Get rid of it in output
		tokenClear := strings.Split(connBuf, token)
		cleared = strings.Join(tokenClear, " ")
		sh.Log.Infof("Found token: %s%s%s ", tui.YELLOW, token, tui.RESET)

		// 3) Get rid of command echoes, prompts, etc. We take them out.
		if strings.Contains(cleared, line) {
			clear := strings.SplitN(cleared, line+";echo", 2)
			cleared = strings.Join(clear, " ")
		}

		// First add what we just processed
		sh.builder.Write([]byte(cleared))

		// 3) Depending on the token index, decide for another ride or not.
		if tokenIndex == 1 && sh.pendingToken {
			result = sh.builder.String() //
			sh.builder.Reset()           // Reset builder
			sh.pendingToken = false      // Reset token counter
			return result, true          // The output can be returned to the user console.
		}

		// Used only by setTokenIndex()
		if sh.tokenSet == false {
			sh.pendingToken = true
			result = sh.builder.String() // Return the result, just in case.
			sh.builder.Reset()           // Reset builder
			return result, true
		}

		sh.pendingToken = true       // Notify for next round
		result = sh.builder.String() // Return the result, just in case.
		return
	}

	// 2) Get rid of command echoes, prompts, etc. We take them out.
	if strings.Contains(connBuf, line) {
		clear := strings.SplitN(connBuf, line+";echo", 2) // Same command check
		cleared = strings.Join(clear, " ")
	}

	// If nothing has been cleared, we directly add the conn buffer.
	if len(cleared) == 0 {
		sh.builder.Write([]byte(connBuf))
	}

	// 2) If we reached this point, this buffer is not the end of the command output.
	//    We add it the command builder string, and we notify caller for another loop.
	sh.builder.Write([]byte(cleared))
	return "", false
}

// findToken - Check bytes in a buffer for a token
func findToken(data []byte, token string) bool {
	firstMatch := bytes.Index(data, []byte(token))
	if firstMatch == -1 {
		return false
	}
	if firstMatch > 0 {
		return true
	}
	return false
}

// NOTE: There is often a lag between a command output and the echo of the cmd we sent.
// Therefore, according to Metasploit: "If the session echoes input we don't need to echo
// the token twice. This setting will persist for the duration of the session."
// We might recalculate from time to time if needed.
func (sh *Shell) setTokenIndex(timeout time.Duration) (err error) {

	sh.Log.Infof("Setting token index for shell session output buffer control")

	// Need to tokens to test
	token := randStringBytesRmndr()
	numeric_token := rand.Int63()

	// Need two commands, accordingly
	testCmd := fmt.Sprintf("echo %d", numeric_token)
	cmd := fmt.Sprintf(testCmd+";echo %s\n", token)

	// Channel used by the command builder to notify timeout control
	done := make(chan struct{})

	// Send package
	err = sh.write(cmd, timeout)
	if err != nil {
		sh.Log.Error(err)
		return
	}

	// Read response
	go func() {
		defer close(done)
		for {
			connBuf := make([]byte, 128) // Temporary buffer
			_, err = sh.stream.Read(connBuf)
			if err != nil {
				sh.Log.Error(err)
				return
			}
			_, complete := sh.readUntilToken(string(connBuf), cmd, token, 0)
			if complete {
				break
			}
		}
	}()

	select {
	case <-done:
		// We just check if the read loop for a simple cmd has returned 1 or 2 echos of the token
		if sh.pendingToken {
			sh.Log.Infof("setTokenIndex() found 2 token echoes: setting tokenIndex to 1")
			sh.pendingToken = false
			sh.tokenIndex = 1
			sh.tokenSet = true
			sh.builder.Reset()
			return
		}
		sh.Log.Infof("setTokenIndex() found 1 token echo: setting index to 0")
		sh.pendingToken = false
		sh.tokenIndex = 0
		sh.tokenSet = true
		sh.builder.Reset()
		return

	case <-time.After(timeout):
		close(done)
		return fmt.Errorf("reading command result from conn stream timed out.")
	}
}

// write - Write to the connection with a timeout and stream error control.
func (sh *Shell) write(line string, timeout time.Duration) (err error) {

	// Channel controls
	done := make(chan struct{})
	errWrite := make(chan error, 1)

	// Asynchronous, timed writing
	go func() {
		defer close(done)
		ilength, err := sh.stream.Write([]byte(line))
		if ilength != len(line) {
			sh.Log.Warnf("length of bytes written to stream and returned"+
				" output length don't match: sendt:%d != returned:%d", len(line), ilength)
		} else if err != nil {
			errWrite <- err
		}
	}()

	// Handle errors and timeouts
	select {
	case err := <-errWrite:
		return err
	case <-time.After(timeout):
		return errors.New("write operation timed out.")
	case <-done:
		sh.Log.Infof("Done writing line: %s", line)
		return nil
	}
}

// randStringBytesRmndr - Easily produce random tokens for output buffer control.
func randStringBytesRmndr() string {

	b := make([]byte, 30)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
