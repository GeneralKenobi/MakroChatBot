package communication

import (
	"errors"
	dg "github.com/bwmarrin/discordgo"
	"github.com/generalkenobi/makrochatbot/logger"
)

// SendToChannel delivers the provided messages to appropriate discord server/channel.
// It makes sure that the sending will not be interrupted - i.e. some other entity won't try to send messages at the same time resulting in mixed messages.
// To make this happen all communication has to go through this function.
func SendToChannel(channelID string, message *dg.MessageSend) {
	sendHelper(channelID, message)
}

// SendToUser delivers the provided messages to appropriate discord server/channel.
// It makes sure that the sending will not be interrupted - i.e. some other entity won't try to send messages at the same time resulting in mixed messages.
// To make this happen all communication has to go through this function.
func SendToUser(userID string, message *dg.MessageSend) {

	// Try to create a channel with the user
	if userChannel, err := session.UserChannelCreate(userID); err == nil {
		// If successful, use the helper to send the message to him
		sendHelper(userChannel.ID, message)
	} else {
		// Otherwise log error
		logger.LogError(errors.New("Can't create channel for communication with user (ID): " + userID))
	}
}

// sendHelper is a helper function for sending messages to Discord.
// It uses sessionMutex to gain ownership of session before sending.
// It will check to make sure message is not nil (if it is it will log an error).
// It will also log an error if sending failed.
func sendHelper(channelID string, message *dg.MessageSend) {

	// Check if the message is nil
	if message == nil {

		// TODO: Try to provide some info about caller / stack trace

		// If it is, create an error message
		errMessage := "Can't send nil message."

		logger.LogError(errors.New(errMessage))

		return
	}

	// Lock the mutex and defer the unlock
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	// Make sure the session is not nil
	if session == nil {
		logger.LogError(errors.New("Can't send message - Discord session is nil"))
		return
	}

	// Send the message
	if _, err := session.ChannelMessageSendComplex(channelID, message); err != nil {
		// If something went wrong, log the error
		logger.LogError(err)
	}
}
