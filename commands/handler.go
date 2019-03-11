package commands

import "strings"

// IllegalPrefix is the only prefix that cannot be used
const IllegalPrefix = ""

// registeredCommands contains all registered commands
var registeredCommands = make(map[string]func([]string))

// commandPrefix is the registered command prefix
var commandPrefix string

// ParseCommand attemts to parse an input string, if it's recognized as a command that command will be executed
// Input is the message sent by a user, userID is the ID of the user that generated the message
func ParseCommand(input, userID string) {

	if !isPrefixRegistered() {
		// TODO: log prefix unregistered error
		return
	}

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

		// Organize the arguments into one slice - userID always comes first
		args := append([]string{userID}, slice[1:]...)

		// If successful, log the event, execute the command and pass the remaining substrings
		function(args)
	}
}

// RegisterCommand registers command - assigns a specific function to a specific string.
// When user types that string the assigned function will be executed.
// DO NOT use any prefixes when registering functions.
func RegisterCommand(name string, function func([]string)) bool {

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

// Returns true if commandPrefix was correctly registered
func isPrefixRegistered() bool {
	return commandPrefix != IllegalPrefix
}
