package communication

import (
	"errors"
	dg "github.com/bwmarrin/discordgo"
	"sync"
)

// session contains the active Discord session. It is used only to send messages in this package.
// It has to be closed by main function.
var session *dg.Session

// sessionMutex is used to synchronize uses of Discors session
var sessionMutex sync.Mutex

// ConnectToDiscord attempts to create and open a Discord session.
// If there are no problems a discordgo Session is returned and error is nil.
// If there were errors the returned session pointer is nil and error contains information about the error.
// If session was already created by calling this method then new session won't be created and an error will be returned.
func ConnectToDiscord(token string) (*dg.Session, error) {

	// Make sure the session wasn't yet created
	if session != nil {
		return nil, errors.New("Discord session is already open")
	}

	// Create a new session for Discord
	s, err := dg.New("Bot " + token)

	// If the process failed return an error
	if err != nil {
		return nil, errors.New("Can't create Discord session. Details: " + err.Error())
	}

	// Assign the session to package variable for future use
	session = s

	// Try to open connection, check if there were errors
	err = session.Open()

	// If so, return an error
	if err != nil {
		return nil, errors.New("Can't open websocket connection to Discord. Details: " + err.Error())
	}

	// If everything went well return the sessionn and no error
	return session, nil
}
