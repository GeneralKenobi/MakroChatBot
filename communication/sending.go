package communication

import (
	dg "github.com/bwmarrin/discordgo"
)

// Send delivers the provided messages to appropriate discord server/channel.
// It makes sure that the sending will not be interrupted - i.e. some other entity won't try to send messages at the same time resulting in mixed messages.
// To make this happen all communication has to go through this function.
func Send(session *dg.Session, channelID string, messages []string) {

	// TODO: Implement the synchronization

	// For each message
	for _, item := range messages {
		// Send it
		session.ChannelMessageSend(channelID, item)
	}

}
