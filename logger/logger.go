package logger

import (
	"fmt"
)

// Log logs the given message on server's standard output
func Log(message string) {

	fmt.Println(message)
}

// LogError logs the given error message on server's standard output
func LogError(err error) {
	Log(err.Error())
}
