package commands

import (
	comm "../communication"
	"errors"
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"strings"
)

// IllegalPrefix is the only prefix that cannot be used
const IllegalPrefix = ""

// registeredCommands contains all registered commands
var registeredCommands = make(map[string]func([]string) ([]string, error))

// commandPrefix is the registered command prefix
var commandPrefix string

// ParseCommand is a handler for incomming messages
func ParseCommand(session *dg.Session, message *dg.MessageCreate) {

	// Check if we're initialized
	if ok, err := isInitialized(); !ok {
		// TODO: Log error - Handler not initialized
		fmt.Println(err.Error())
		return
	}

	// Don't respond to any bots (including ourselvers)
	if message.Author.Bot {
		return
	}

	// The message content, converted to lower-case
	input := strings.ToLower(message.Content)

	// Check if the input string contains the registered prefix - if not, return as there is nothing we can do
	if !strings.HasPrefix(input, commandPrefix) {
		return
	}

	// Remove the prefix from the input
	input = input[len(commandPrefix):]

	// Split the input based on whitespace - each substring now contains only non-whitespace characters
	slice := strings.Fields(input)

	// If the slice resulted in no substrings - return (command was empty)
	if len(slice) < 1 {
		return
	}

	// Try to get a function matching to the first substring (which is the command name, the eventual remaining substrings are parameters)
	if function, ok := registeredCommands[slice[0]]; ok {

		// TODO: Log that a command is executed

		// Organize the arguments into one slice - username always comes first
		args := append([]string{message.Author.Username}, slice[1:]...)

		// If successful, log the event, execute the command and pass the remaining substrings
		if messages, err := function(args); err == nil {
			comm.Send(session, message.ChannelID, messages)
		}
	}
}

// RegisterCommand registers command - assigns a specific function to a specific string.
// When user types that string the assigned function will be executed.
// DO NOT use any prefixes when registering functions.
// Command names are case insensitive.
func RegisterCommand(name string, function func([]string) ([]string, error)) bool {

	// Check if there's already a command registered for that name
	if _, ok := registeredCommands[name]; ok {
		// If so, return false - it won't be overwritted
		return false
	}

	// If the name was not present, then we can register a function to it
	registeredCommands[name] = function
	return true
}

// RegisterCommandPrefix registers the given string as recognized command prefix.
// Prefix can be registered only once. Future calls to this function won't do anything.
// Returns true if registration was successful, false otherwise.
// All prefixes are legal except for the IllegalPrefix "" (empty string)
func RegisterCommandPrefix(prefix string) bool {

	// If there already is a prefix defined or the provided prefix is illegal, return false
	if commandPrefix != "" || prefix == "" {
		return false
	}

	// Otherwise assign the new value and return success
	commandPrefix = prefix
	return true
}

// isPrefixRegistered returns true if commandPrefix was correctly registered
func isPrefixRegistered() bool {
	return commandPrefix != IllegalPrefix
}

// isInitialized returns true if all necessary initialization was performed. If it's not the case returns false and assigns an error (based on first
// caught unitialized aspect)
func isInitialized() (bool, error) {

	// Check if prefix was initialized
	if !isPrefixRegistered() {
		return false, errors.New("Handler error: prefix was not initialized")
	}

	return true, nil
}
