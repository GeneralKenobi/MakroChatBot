package logger

import (
	ct "github.com/GeneralKenobi/MakroChatbot/customtypes"
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

// LogCommand logs information about command being invoked by a user
func LogCommand(guildID, channelID string, args ct.CommandArgs) {

	// Create log message. Specify the guild, channel, command, user and arguments
	log := "Guild: " + guildID +
		", Channel: " + channelID +
		", Command: \"" + args.CommandName +
		"\", User: " + args.Username + " (ID: " + args.UserID + ")" +
		", args: " + strings.Join(args.UserArgs, ", ")

	Log(log)
}
