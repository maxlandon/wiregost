package core

// Unix-specific shell functions.

func Shell(cmd string) (string, error) {
	return Exec("/bin/sh", []string{"-c", cmd})
}
