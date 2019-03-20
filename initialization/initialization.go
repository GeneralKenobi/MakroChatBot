package initialization

import (
	dg "github.com/bwmarrin/discordgo"

	"github.com/generalkenobi/makrochatbot/commands/handler"
	pm "github.com/generalkenobi/makrochatbot/commands/platformmonitor"
	"github.com/generalkenobi/makrochatbot/commands/reactions"
	"github.com/generalkenobi/makrochatbot/commands/roll"
	"github.com/generalkenobi/makrochatbot/communication"
	"github.com/generalkenobi/makrochatbot/configuration"
	ct "github.com/generalkenobi/makrochatbot/customtypes"
	"github.com/generalkenobi/makrochatbot/logger"
	"math/rand"
	"time"
)

// Run performs all initization for the program.
// If everything goes well then an open discordgo session is returned - IT HAS TO BE CLOSED BEFORE CLOSING THE PROGRAM.
// If any crucial part of initization fails then the returned session will be null and error will contain information about the error.
func Run() (*dg.Session, error) {

	// Perform the necessary initialization
	config, session, err := requiredInit()

	// If something went wrong, return the error
	if err != nil {
		return nil, err
	}

	logger.Log("Core initialization complete")

	// Add handler for incomming messages
	session.AddHandler(handler.ParseCommand)

	// Register command prefix
	handler.RegisterCommandPrefix(config.CommandPrefix)

	// Call the helper function to register all commands
	registerCommands()

	logger.Log("Command handler initialized")

	// Seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	logger.Log("Random number generator seeded")

	pm.Start(10)

	// Finally return the session and no error
	return session, nil
}

// registerCommands registers all commands handled by the bot
func registerCommands() {

	handler.RegisterCommand("roll", roll.Roll)
	handler.RegisterCommand("group1", reactions.ImageReaction)
	handler.RegisterCommand("group2", reactions.ImageReaction)
	handler.RegisterCommand("subscribe", pm.CreateMonitorSubscription)
	handler.RegisterCommand("unsubscribeall", pm.RemoveAllSubscriptions)
}

// requiredInit performs crucial initialization tasks - loading config file and opening Discord session.
// If any of these fails then the program can't run.
// If everything went well the obtained Config and an opened discordgo session are returned, error in that case is nil.
// If something went wrong both config and session are nil and error contains information on what went wrong.
func requiredInit() (*ct.Config, *dg.Session, error) {

	// Get the configuration file
	config, err := configuration.GetConfig()

	// If the config file couldn't be obtained, return nil (can't initialize) and the error generated by config file fetcher
	if err != nil {
		return nil, nil, err
	}

	// Open discord session
	session, err := communication.ConnectToDiscord(config.Token)

	// If the session couldn't be opened, return nil (can't initialize) and the generated error
	if err != nil {
		return nil, nil, err
	}

	return &config, session, nil
}
