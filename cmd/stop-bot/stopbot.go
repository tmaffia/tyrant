package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	botToken  = os.Getenv("STOP_BOT_TOKEN")
	appID     = os.Getenv("STOP_BOT_APP_ID")
	publicKey = os.Getenv("STOP_BOT_PUBLIC_KEY")
)

func main() {
	if botToken == "" ||
		appID == "" ||
		publicKey == "" {

		fmt.Println("ENV Variables are not set properly")
		return
	}

	log.Println(botToken)
	log.Println(appID)
	log.Println(publicKey)

	s, err := discordgo.New("Bot " + botToken)

	if err != nil {
		log.Println(err.Error())
	}

	s.AddHandler(message)

	s.Identify.Intents = discordgo.IntentDirectMessages

	err = s.Open()

	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	log.Println("Stop Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
