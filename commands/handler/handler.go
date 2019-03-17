package handler

import (
	"errors"
	dg "github.com/bwmarrin/discordgo"
	"github.com/generalkenobi/makrochatbot/communication"
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"github.com/generalkenobi/makrochatbot/logger"
	"strings"
)

// IllegalPrefix is the only prefix that cannot be used
const IllegalPrefix = ""

// registeredCommands contains all registered commands
var registeredCommands = make(map[string]ct.CommandHandler)

// commandPrefix is the registered command prefix
var commandPrefix string

// ParseCommand is a handler for incomming messages
func ParseCommand(session *dg.Session, message *dg.MessageCreate) {

	// Check if we're initialized
	if ok, err := isInitialized(); !ok {
		logger.LogError(err)
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

		// Construct a struct with arguments
		args := ct.CommandArgs{
			CommandName: slice[0],
			Username:    message.Author.Username,
			UserID:      message.Author.ID,
			UserArgs:    slice[1:]}

		// Log command execution
		logger.LogCommand(message.GuildID, message.ChannelID, args)

		// If there were no errors when running the command and it returned a message to send to the channel
		if output, err := function(&args); err == nil {
			if output != nil {
				// Send the produced message to the source channel
				communication.SendToChannel(message.ChannelID, output)
			}
		} else {
			// Otherwise log the error
			logger.LogError(err)
		}
	}
}

// RegisterCommand registers command - assigns a specific function to a specific string.
// When user types that string the assigned function will be executed.
// DO NOT use any prefixes when registering functions.
// Command names are case insensitive.
func RegisterCommand(name string, function ct.CommandHandler) bool {

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
