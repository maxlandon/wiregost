package console

// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/moloch--/asciicast"
	"golang.org/x/exp/slog"
	"golang.org/x/term"

	"github.com/maxlandon/wiregost/internal/client/assets"
	"github.com/maxlandon/wiregost/internal/proto/clientpb"
	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
)

const (
	// ANSI Colors
	Normal    = "\033[0m"
	Dim       = "\033[2m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Orange    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	Bold      = "\033[1m"
	Clearln   = "\r\x1b[2K"
	UpN       = "\033[%dA"
	DownN     = "\033[%dB"
	Underline = "\033[4m"

	// Other special colors
	LightRed = "\033[38;5;210m"

	// background colors
	BgDarkGray  = "\033[100m"
	BgRed       = "\033[41m"
	BgGreen     = "\033[42m"
	BgOrange    = "\033[43m"
	BgLightBlue = "\033[104m"

	// Info - Display colorful information
	Info = Bold + Cyan + "[*] " + Normal
	// Warn - Warn a user
	Warn = Bold + Red + "[!] " + Normal
	// Debug - Display debug information
	Debug = Bold + Purple + "[-] " + Normal
	// Woot - Display success
	Woot = Bold + Green + "[$] " + Normal
	// Success - Diplay success
	Success = Bold + Green + "[+] " + Normal
)

// Logger is an io.Writer that sends data to the server.
type Logger struct {
	name   string
	Stream rpcpb.Core_ClientLogClient
}

func (l *Logger) Write(buf []byte) (int, error) {
	err := l.Stream.Send(&clientpb.ClientLogData{
		Stream: l.name,
		Data:   buf,
	})
	return len(buf), err
}

// ClientLogStream requires a log stream name, used to save the logs
// going through this stream in a specific log subdirectory/file.
func (con *Client) ClientLogStream(name string) (*Logger, error) {
	stream, err := con.Rpc.ClientLog(context.Background())
	if err != nil {
		return nil, err
	}
	return &Logger{name: name, Stream: stream}, nil
}

func (con *Client) setupLogger(writers ...io.Writer) {
	logWriter := io.MultiWriter(writers...)
	jsonOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	con.jsonHandler = slog.NewJSONHandler(logWriter, jsonOptions)

	// Log all commands before running them.
	con.App.PreCmdRunLineHooks = append(con.App.PreCmdRunLineHooks, con.logCommand)
}

// logCommand logs non empty commands to the client log file.
func (con *Client) logCommand(args []string) ([]string, error) {
	if len(args) == 0 {
		return args, nil
	}
	logger := slog.New(con.jsonHandler).With(slog.String("type", "command"))
	logger.Debug(strings.Join(args, " "))
	return args, nil
}

func (con *Client) setupAsciicastRecord(logFile *os.File, server io.Writer) {
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		x, y = 80, 80
	}

	// Save the asciicast to the server and a local file.
	destinations := io.MultiWriter(logFile, server)

	encoder := asciicast.NewEncoder(destinations, x, y)
	encoder.WriteHeader()

	// save existing stdout | MultiWriter writes to saved stdout and file
	out := os.Stdout
	mw := io.MultiWriter(out, encoder)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout,
	// stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	go io.Copy(mw, r)
}

func getConsoleLogFile() *os.File {
	logsDir := assets.GetConsoleLogsDir()
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("%s.log", dateTime))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Could not open log file: %s", err)
	}
	return logFile
}

func getConsoleAsciicastFile() *os.File {
	logsDir := assets.GetConsoleLogsDir()
	dateTime := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("asciicast_%s.log", dateTime))
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Could not open log file: %s", err)
	}
	return logFile
}

//
// -------------------------- [ Logging ] -----------------------------
//
// Logging function below differ slightly from their counterparts in client/log package:
// These below will print their output regardless of the currently active menu (server/implant),
// while those in the log package tie their output to the current menu.

func (con *Client) Printf(format string, args ...any) {
	logger := slog.NewLogLogger(con.jsonHandler, slog.LevelInfo)
	logger.Printf(format, args...)

	con.printf(format, args...)
}

// Println prints an output without status and immediately below the last line of output.
func (con *Client) Println(args ...any) {
	logger := slog.New(con.jsonHandler)
	format := strings.Repeat("%s", len(args))
	logger.Info(fmt.Sprintf(format, args))
	con.printf(format, args...)
}

// PrintInfof prints an info message immediately below the last line of output.
func (con *Client) PrintInfof(format string, args ...any) {
	logger := slog.New(con.jsonHandler)

	logger.Info(fmt.Sprintf(format, args...))

	con.printf(Clearln+Info+format, args...)
}

// PrintSuccessf prints a success message immediately below the last line of output.
func (con *Client) PrintSuccessf(format string, args ...any) {
	logger := slog.New(con.jsonHandler)

	logger.Info(fmt.Sprintf(format, args...))

	con.printf(Clearln+Success+format, args...)
}

// PrintWarnf a warning message immediately below the last line of output.
func (con *Client) PrintWarnf(format string, args ...any) {
	logger := slog.New(con.jsonHandler)

	logger.Warn(fmt.Sprintf(format, args...))

	con.printf(Clearln+"⚠️  "+Normal+format, args...)
}

// PrintErrorf prints an error message immediately below the last line of output.
func (con *Client) PrintErrorf(format string, args ...any) {
	logger := slog.New(con.jsonHandler)

	logger.Error(fmt.Sprintf(format, args...))

	con.printf(Clearln+Warn+format, args...)
}

// PrintEventInfof prints an info message with a leading/trailing newline for emphasis.
func (con *Client) PrintEventInfof(format string, args ...any) {
	logger := slog.New(con.jsonHandler).With(slog.String("type", "event"))

	logger.Info(fmt.Sprintf(format, args...))

	con.printf(Clearln+"\r\n"+Info+format+"\r", args...)
}

// PrintEventErrorf prints an error message with a leading/trailing newline for emphasis.
func (con *Client) PrintEventErrorf(format string, args ...any) {
	logger := slog.New(con.jsonHandler).With(slog.String("type", "event"))

	logger.Error(fmt.Sprintf(format, args...))

	con.printf(Clearln+"\r\n"+Warn+format+"\r", args...)
}

// PrintEventSuccessf a success message with a leading/trailing newline for emphasis.
func (con *Client) PrintEventSuccessf(format string, args ...any) {
	logger := slog.New(con.jsonHandler).With(slog.String("type", "event"))

	logger.Info(fmt.Sprintf(format, args...))

	con.printf(Clearln+"\r\n"+Success+format+"\r", args...)
}
