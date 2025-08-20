package internal

import "fmt"

// LogVerbose is a helper function used to give the user detailed output when verbose flag is used.
func LogVerbose(on bool, format string, a ...any) {
	if on {
		fmt.Printf(format+"\n", a...)
	}
}
