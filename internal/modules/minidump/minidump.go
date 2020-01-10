package minidump

import (
	"fmt"
	"strconv"
)

// Parse is the initial entry point for all extended modules. All validation checks and processing will be performed here
// The function input types are limited to strings and therefore require additional processing
func Parse(options map[string]string) ([]string, error) {
	// Convert PID to integer
	if options["pid"] != "" && options["pid"] != "0" {
		_, errPid := strconv.Atoi(options["pid"])
		if errPid != nil {
			return nil, fmt.Errorf("there was an error converting the PID to an integer:\r\n%s", errPid.Error())
		}
	}

	command, errCommand := GetJob(options["process"], options["pid"], options["tempLocation"])
	if errCommand != nil {
		return nil, fmt.Errorf("there was an error getting the minidump job:\r\n%s", errCommand.Error())
	}

	return command, nil
}

// GetJob returns a string array containing the commands, in the proper order, to be used with agents.AddJob
func GetJob(process string, pid string, tempLocation string) ([]string, error) {
	return []string{"Minidump", process, pid, tempLocation}, nil
}
