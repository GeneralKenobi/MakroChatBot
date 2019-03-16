package customtypes

import dg "github.com/bwmarrin/discordgo"

// CommandHandler is a type definition of a function that can be used as a command handler.
// First return value is a message to send to the channel from which the command was invoked.
// Second return value is an error, if it's not nil then the message from first return value won't be sent to the source channel.
type CommandHandler func(*CommandArgs) (*dg.MessageSend, error)
