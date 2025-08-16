package runner

import (
	"os"
	"os/exec"
)

// RunInteractive runs a command with stdio attached so SSH can propmt for passphrase
func RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=1",
		"GIT_SSH_COMMAND=ssh -o BatchMode=no",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
