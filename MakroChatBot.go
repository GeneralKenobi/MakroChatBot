package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./usercommands"

	"github.com/bwmarrin/discordgo"
)

// Configuration data
var (
	commandPrefix = "1"
	botID         = ""
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

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

	// Update status with a cool message
	discord.UpdateStatus(0, "Hello There!")

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

	if message.Content == "!Roll" {

		number, reaction := usercommands.Roll()

		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%d", number))
		discord.ChannelMessageSend(message.ChannelID, reaction)
	}

	if message.Content == "!Group1" {
		discord.ChannelMessageSend(message.ChannelID, "https://i.kym-cdn.com/photos/images/facebook/001/290/942/c31.png")
	}

	if message.Content == "!Group2" {
		discord.ChannelMessageSend(message.ChannelID, "https://vignette.wikia.nocookie.net/internet-meme/images/6/6e/Pogchamp.jpg/revision/latest?cb=20180310053228")
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
