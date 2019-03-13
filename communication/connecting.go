package communication

import (
	"errors"
	dg "github.com/bwmarrin/discordgo"
)

// ConnectToDiscord attempts to create and open a Discord session.
// If there are no problems a discordgo Session is returned and error is nil.
// If there were errors the returned session pointer is nil and error contains information about the error.
func ConnectToDiscord(token string) (*dg.Session, error) {

	// Create a new session for Discord
	session, err := dg.New("Bot " + token)

	// If the process failed return an error
	if err != nil {
		return nil, errors.New("Can't create Discord session. Details: " + err.Error())
	}

	// Try to open connection, check if there were errors
	err = session.Open()

	// If so, return an error
	if err != nil {
		return nil, errors.New("Can't open websocket connection to Discord. Details: " + err.Error())
	}

	// If everything went well return the sessionn and no error
	return session, nil
}
