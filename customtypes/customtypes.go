package customtypes

import dg "github.com/bwmarrin/discordgo"

// CommandHandler is a type definition of a function that can be used as a command handler.
// First return value is a message to send to the channel from which the command was invoked.
// Second return value is an error, if it's not nil then the message from first return value won't be sent to the source channel.
type CommandHandler func(*CommandArgs) (*dg.MessageSend, error)

// CommandArgs is a struct for arguments that are passed to command handlers
type CommandArgs struct {

	// Name of the invoked command
	CommandName string

	// Name of the user that invoked the command
	Username string

	// ID of the user that invoked the command
	UserID string

	// Arguments passed by the user
	UserArgs []string
}

// Config contains data necessary to configure the bot
type Config struct {

	// Token to use when connecting to the server
	Token string

	// Token to use when connecting to the server
	CommandPrefix string

	// Time period (in seconds) between two subsequent platform checks
	PlatformMonitoringPeriod int
}
