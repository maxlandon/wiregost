package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
	"github.com/sirupsen/logrus"
)

// Shell - The Shell object implements an interactive Command Shell around a stream.
// This object is meant to provide basic functionality for manipulating a remote
// shell session such as a reverse /bin/sh one-liner: this remote stream will only
// provide primitive I/O and command execution interface. This object wraps it around
// a few methods and attempts to transform it into a more usable command-line interface.
type Shell struct {
	// Base session information, route and logging.
	*Session

	// stream - I/0 stream for this session. This can be anything: net.Conn or
	// a mux.Stream passed by a handler/route. It can be the OS Stdin/Stdout/Stderr
	// of a host, etc. Reader, Writer and Closer may even point to different streams.
	stream  io.ReadWriteCloser
	reader  *bufio.Reader // Fine-grained control over how to read the stream.
	timeout time.Duration // Default timeout when not set from user

	tokenIndex   int      // influences read loop on the stream and delimitations
	tokenSet     bool     // indicates we have set the index by echoing a test token.
	pendingToken bool     // used by the shell/builder for output concatenation.
	unwished     []string // tokens that we filter each time.
	prompt       string   // We store the remote prompt.

	builder *strings.Builder // builder - Temp storage for stream output
}

// NewShell - Instantiates a new Session that implements a command shell around a logical stream.
func NewShell(stream io.ReadWriteCloser) (sh *Shell) {

	sh = &Shell{
		New(),              // The session is interactive.
		stream,             // The underlying communication medium for this session.
		nil,                // Instantiate the reader below, in Setup().
		10 * time.Second,   // Default timeout
		0,                  // The token is by default 0. Will check in future if needs change.
		false,              // Not set yet, will be once only.
		false,              // No tokens read yet.
		[]string{},         // No tokens to filter yet
		"",                 // No prompt on remote end
		&strings.Builder{}, // Used to build the output from conn
	}

	sh.Type = serverpb.SessionType_SHELL

	sh.Log.Logger.Out = os.Stdout
	sh.Log = sh.Log.WithField("type", "shell") // Log settings

	return
}

// Setup - The shell sets up the token index for correct delimitation of command output,
// finds the remote prompt and saves it, adds unwished tokens to a list for trimming, etc.
func (sh *Shell) Setup(fullEnv bool) (err error) {

	sLog := sh.Log.WithField("stream", "setup") // Pass this log to setup functions

	sh.reader = bufio.NewReader(sh.stream)              // Initialize conn reader
	sh.getRemotePrompt(sLog)                            // Get prompt out of first output
	sh.unwished = append(sh.unwished, defaultTokens...) // Add primary unwished tokens
	err = sh.setTokenIndex(sLog, sh.timeout)            // Set token index
	if err != nil {
		sLog.Errorf("Error during command shell setup: " + err.Error())
		return err
	}

	err = sh.LoadRemoteEnvironment(sLog, fullEnv)
	if err != nil {
		sLog.Errorf("failed to load remote shell environment: %s", err.Error())
	}
	return
}

// LoadRemoteEnvironment - Primary method for retrieving the target environment variables, and assigning
// them to this session, therefore available to both console completions and further session code.
func (sh *Shell) LoadRemoteEnvironment(log *logrus.Entry, full bool) (err error) {

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

// Cleanup - Clean any state related to this Interactive Session.
// Should call the *Session implementation at some point.
func (sh *Shell) Cleanup() (err error) {
	return
}

// Kill - Terminate the Interactive session. Cleans up the resources and
// calls the *Session Kill() implementation for deregistering the Session.
func (sh *Shell) Kill() (err error) {

	// This involves handling the way we kill the ReadWriteCloser.
	// The issue here is that we don't know anything about
	return
}

// RunCommand - Given a command and its arguments, with an optional timeout, send it through the
// shell stream, and wait for response. This should be synchronous given that we cannot run
// concurrent operations on a remote shell while separating their respective stdin/stdout.
func (sh *Shell) RunCommand(cmd string, args []string, timeout time.Duration) (result string, err error) {

	// default timeout
	if timeout == 0 {
		timeout = sh.timeout
	}

	// Already called in LoadRemoteEnvironment, but just in case.
	if !sh.tokenSet {
		err = sh.setTokenIndex(sh.Log, timeout)
		if err != nil {
			sh.Log.Error(err)
			return
		}
	}

	cmd = cmd + " " + strings.Join(args, " ") // Concat command and args.

	// Get platform/OS for correct dispatch. Use associated host/DB, or anything.
	// Waiting, only Unix...
	return sh.runCommandUnix(cmd, timeout)
}

// runCommandUnix - Adapt the token for Unix remote endpoints
func (sh *Shell) runCommandUnix(cmd string, timeout time.Duration) (out string, err error) {

	// Forge token and command, and send to remote
	token := randStringBytesRmndr()
	forgedCmd := sh.forgeCommand(cmd, token)

	if err = sh.write(forgedCmd, timeout); err != nil {
		sh.Log.Error(err)
		return
	}
	sh.Log.Debugf(tui.Green("command: ") + tui.Bold(cmd))

	// 2. Read connection.
	done := make(chan struct{})
	processed := make(chan string, 1)
	go func(chan struct{}, chan string) {
		defer close(done)
		for {
			select {
			default:
				// Read all output until one/both tokens are found
				output, err := sh.readCommandOuput(cmd, token, sh.tokenIndex)
				if err != nil {
					return // We already logged the error
				}

				// Process output
				out, err = sh.processRawLine(output, cmd, token)
				if err != nil {
					return // We already logged the error
				}
				processed <- out
				return
			case <-done:
				return
			}
		}
	}(done, processed)

	// We wait either for the response body, or a timeout.
	for {
		select {
		case out = <-processed:
			sh.Log.Debugf(tui.Dim("result: ") + tui.Bold(cmd))
			return out, nil
		case <-time.After(timeout):
			close(done)
			// We still give out, in case it has something in it still.
			return out, fmt.Errorf("reading command result from conn stream timed out")
		}
	}
}

// runCommandWindows - Run a command on a Windows endpoint and adapt the token.
func (sh *Shell) runCommandWindows(line string, timeout time.Duration) (out string, err error) {
	return
}

// forgeCommand - depending on the input, perform a few adjustements.
func (sh *Shell) forgeCommand(cmd, token string) (line []byte) {
	if cmd == "\n" || cmd == "\n " || cmd == "" || cmd == " " {
		return []byte(fmt.Sprintf("echo %s\n", token))
	}
	return []byte(cmd + fmt.Sprintf(";echo %s\n", token))
}

// readUntilToken - The Shell's string builder is being passed temporary connection buffers
// and it processes them: finds command tokens to delimit output, trims command echos, prompts, etc.
func (sh *Shell) readCommandOuput(cmd, token string, index int) (result string, err error) {

	tokenMatch := regexp.MustCompile(token)

	for {
		line := make([]byte, 4096)
		_, err = sh.reader.Read(line)
		if err != nil {
			sh.Log.Errorf("error reading the connection stream: " + err.Error())
		}

		// Count token occurences in this buffer
		switch tokens := len(tokenMatch.FindAllIndex(line, -1)); {

		// If we found both in the same buffer
		case tokens == 2:
			sh.Log.Tracef("found 2 tokens in the same buffer")
			sh.pendingToken = false
			return string(line), nil

		// If we found one in the same buffer
		case tokens == 1:
			// 1st token is found
			if !sh.pendingToken {
				sh.Log.Tracef("Found 1st token: %s%s%s ", tui.YELLOW, token, tui.RESET)
				if index == 1 {
					sh.Log.Tracef("Waiting for second token...")
					sh.pendingToken = true
					continue
				}
				sh.Log.Tracef("Token index is 0: breaking read loop")
				result = sh.builder.String() + string(line)
				sh.builder.Reset()
				return
			}
			// 2 second token was in fact found
			sh.Log.Tracef("Found 2nd token: %s%s%s%s .Breaking read loop",
				tui.YELLOW, tui.BOLD, token, tui.RESET)
			sh.pendingToken = false
			result = sh.builder.String() + string(line)
			sh.builder.Reset()
			return

		// We didn't find any in the buffer
		case tokens == 0:
			if len(line) > 0 { // Buffer is not empty, add it for next ride
				sh.builder.Write(line)
			}
			continue // Go for another loop
		}
	}
}

// processRawLine - Analyzes the output of a command on the wire and trims it from all unwished items.
func (sh *Shell) processRawLine(line string, cmd, token string) (processed string, err error) {

	processed = string(bytes.Trim([]byte(line), "\x00"))                // Clean null bytes
	processed = strings.Join(strings.Split(processed, cmd), " ")        // Put out the command
	processed = strings.Join(strings.Split(processed, token+"\n"), " ") // Token

	// All other unwished tokens
	for _, tok := range sh.unwished {
		processed = strings.ReplaceAll(processed, tok, "")
	}

	return
}

// setTokenIndex - There is often a lag between a command output and the echo of the cmd we sent.
// Therefore, according to Metasploit: "If the session echoes input we don't need to echo
// the token twice. This setting will persist for the duration of the session."
// We might recalculate from time to time if needed.
func (sh *Shell) setTokenIndex(log *logrus.Entry, timeout time.Duration) (err error) {

	log.Debugf("Setting token index for shell session output buffer control")

	// Need two tokens to test, and two commands
	token, numericToken := randStringBytesRmndr(), rand.Int63()
	testCmd := fmt.Sprintf("echo %d", numericToken)
	cmd := fmt.Sprintf(testCmd+";echo %s\n", token)

	// Send package
	err = sh.write([]byte(cmd), timeout)
	if err != nil {
		sh.Log.Error(err)
		return
	}

	// Do the process in a goroutine and time it.
	done := make(chan struct{})
	go func(<-chan struct{}) {
		defer close(done)
		select {
		default:
			// Read all output until one/both tokens are found
			output, err := sh.readCommandOuput(cmd, token, 0)
			if err != nil {
				sh.Log.Errorf("error when setting ouput token: %s", err.Error())
				return
			}

			items := strings.Split(output, "\n") // Get all lines separately
			var tokenLine string                 // check for the prompt line first

			// Edge case: the first string is the prompt, so we loop again
			if strings.TrimSpace(items[0]) == sh.prompt {
				tokenLine = items[1]
			} else {
				tokenLine = items[0]
			}

			// Depending on the outcome, we read 1/2 line on the conn, for clearing.
			_, _, err = sh.reader.ReadLine() // One is in common
			if err != nil {
				log.Warnf("could not clear extra lines from conn after token setup")
			}

			nb, err := strconv.Atoi(tokenLine)
			if nb == 0 && err != nil {
				sh.tokenIndex = 1
				_, _, err = sh.reader.ReadLine() // One is because we indeed have an echo
				if err != nil {
					log.Warnf("could not clear extra lines from conn after token setup")
				}
			} else {
				sh.tokenIndex = 0
			}
			sh.tokenSet = true
			sh.pendingToken = false // double check
		case <-done:
			return
		}
	}(done)

	select {
	case <-done:
		log.Debugf("test output token done: tokenIndex = %d", sh.tokenIndex)
		return
	case <-time.After(timeout):
		close(done)
		return fmt.Errorf("reading command result from conn stream timed out")
	}
}

// write - Write to the connection with a timeout and stream error control.
func (sh *Shell) write(line []byte, timeout time.Duration) (err error) {

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
		return errors.New("write operation timed out")
	case <-done:
		sh.Log.Debugf("Done writing line: %s", line)
		return nil
	}
}

// getRemotePrompt - Asks for and reads the first line from the stream, which should contain the prompt.
func (sh *Shell) getRemotePrompt(log *logrus.Entry) (err error) {

	// Send an empty command for refresh
	if _, err = sh.stream.Write([]byte("\n")); err != nil {
		log.Warnf("failed to write to stream when requesting prompt: %s", err.Error())
		return
	}

	// Read the first line out of the stream: its the remote prompt
	promptBytes, _, err := sh.reader.ReadLine()
	if err != nil {
		log.Warnf("failed to set prompt with first input: %s", err.Error())
		return
	}
	sh.prompt = strings.TrimSpace(strings.Trim(string(promptBytes), "\n"))
	log.Debugf("saved remote prompt string: %s", tui.Blue(sh.prompt))

	sh.unwished = append(sh.unwished, sh.prompt)
	return
}

// defaultTokens - tokens (words or patterns) that we trim from every command output
var defaultTokens = []string{";echo "}

// letterBytes - used to produce random string tokens
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randStringBytesRmndr - Easily produce random tokens for output buffer control.
func randStringBytesRmndr() string {

	b := make([]byte, 30)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
