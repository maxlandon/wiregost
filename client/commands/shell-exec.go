package commands

// ShellExec - Executes the input line through a system shell
func ShellExec(args []string) error {
	return nil
}

// BinaryExec - Looks up for the binary name, and execute it if found in $PATH
func BinaryExec(executable string, args []string) (res string, err error) {
	return
}

func inputIsBinary(args []string) bool {
	return false
}
