package logger

import (
	"fmt"
	"strings"
	"time"
)

// Log logs the given message on server's standard output
func Log(message string) {

	// Prefix for the message
	logPrefix := time.Now().Format("2006-01-02 15:04:05") + " - "

	fmt.Println(logPrefix + message)
}

// LogError logs the given error message on server's standard output
func LogError(err error) {
	Log(err.Error())
}

func LogCommand(guildID, channelID, command string, args []string) {

	// Create log message. Specify the guild, channel, command and arguments
	log := "Guild: " + guildID +
		", Channel: " + channelID +
		", Command: \"" + command +
		"\", args: " + strings.Join(args, ", ")

	Log(log)
}
