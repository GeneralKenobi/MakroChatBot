package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./initialization"
	"./usercommands"

	"github.com/bwmarrin/discordgo"
)

func main() {

	session, err := initialization.Run()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Update status with a cool message
	session.UpdateStatus(0, "I Love democracy")

	defer session.Close()

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
func messageHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// User that sent the message
	user := message.Author

	// Don't respond to ourselves or other bots
	if user.ID == "xd" || user.Bot {
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
