package internal

import "fmt"

func LogVerbose(on bool, format string, a ...interface{}) {
	if on {
		fmt.Printf(format+"\n", a...)
	}
}
