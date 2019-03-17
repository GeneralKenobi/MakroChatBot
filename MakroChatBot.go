package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/generalkenobi/makrochatbot/initialization"
	"github.com/generalkenobi/makrochatbot/logger"
)

func main() {

	session, err := initialization.Run()

	if err != nil {
		logger.LogError(err)
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
