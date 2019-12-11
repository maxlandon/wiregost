package session

// This file contains the command parsing functions.
// All inputs to the shell are stripped and parsed here,
// then they will be sent for processing to their respective
// command handlers.

import (
	"strings"

	"github.com/evilsocket/islazy/str"
)

func ParseCommands(line string) []string {
	args := []string{}
	buf := ""

	singleQuoted := false
	doubleQuoted := false
	finish := false

	line = strings.Replace(line, `""`, `"<empty>"`, -1)
	line = strings.Replace(line, `''`, `"<empty>"`, -1)
	for _, c := range line {
		switch c {
		case ';':
			if !singleQuoted && !doubleQuoted {
				finish = true
			} else {
				buf += string(c)
			}

		case '"':
			if doubleQuoted {
				// finish of quote
				doubleQuoted = false
			} else if singleQuoted {
				// quote initiated with ', so we ignore it
				buf += string(c)
			} else {
				// quote init here
				doubleQuoted = true
			}

		case '\'':
			if singleQuoted {
				singleQuoted = false
			} else if doubleQuoted {
				buf += string(c)
			} else {
				singleQuoted = true
			}

		default:
			buf += string(c)
		}

		if finish {
			buf = strings.Replace(buf, `<empty>`, `""`, -1)
			args = append(args, buf)
			finish = false
			buf = ""
		}
	}

	if len(buf) > 0 {
		buf = strings.Replace(buf, `<empty>`, `""`, -1)
		args = append(args, buf)
	}

	cmds := make([]string, 0)
	for _, cmd := range args {
		cmd = str.Trim(cmd)
		if cmd != "" || (len(cmd) > 0 && cmd[0] != '#') {
			cmds = append(cmds, cmd)
		}
	}

	return cmds
}
