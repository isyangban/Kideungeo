package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/isyangban/Kideungeo/internal/kideungeo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	if len(DISCORD_BOT_TOKEN) == 0 {
		log.Fatal("Don't have required env variables")
	}

	bot := kideungeo.New(DISCORD_BOT_TOKEN)

	// Wait here until CTRL-C or other term signal is received.
	bot.Start()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
