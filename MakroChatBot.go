package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Configuration data
var (
	commandPrefix = "1"
	botID         = "NTUzOTgxOTg2MTE3Nzc5NDc3.D2WZGg.OVHt_avuWXlWfLw_p3xPe31zr58"
)

func main() {

	// Create a new session for Discord
	discord, err := discordgo.New("Bot " + botID)

	// Check if there was an error
	panicOnError("error creating discord session", err)

	// Add handler for incomming messages
	discord.AddHandler(commandHandler)

	// Try to open connection, check if there were errors
	err = discord.Open()
	panicOnError("Error opening connection to Discord", err)
	defer discord.Close()

	// Wait here for control signal that closes the bot
	fmt.Println("Running, press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Printf("Finishing...")
}

// If error is not nil prints the message and error and panics
func panicOnError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

// Handles incomming messages
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// User that sent the message
	user := message.Author

	// Don't respond to ourselves or other bots
	if user.ID == botID || user.Bot {
		return
	}

	// Simple echo for testing purposes
	if message.Content == "xd123" {
		discord.ChannelMessageSend(message.ChannelID, "xd123")
	}
}

// Status update, maybe I'll do it later
//discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
//	err = discord.UpdateStatus(0, "A friendly helpful bot!")
//	if err != nil {
//		fmt.Println("Error attempting to set my status")
//	}
//	servers := discord.State.Guilds
//	fmt.Printf("Running on %d servers", len(servers))
//})
