package main

import (
	"errors"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type StopBot struct {
	botToken  string
	appID     string
	publicKey string
	intents   discordgo.Intent
}

func initStopBot() (StopBot, error) {
	sb := StopBot{
		botToken:  os.Getenv("STOP_BOT_TOKEN"),
		appID:     os.Getenv("STOP_BOT_APP_ID"),
		publicKey: os.Getenv("STOP_BOT_PUBLIC_KEY"),
		intents:   discordgo.IntentDirectMessages,
	}
	return sb, nil
}

func (b StopBot) message(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func (sb StopBot) run() error {

	if sb.botToken == "" ||
		sb.appID == "" ||
		sb.publicKey == "" {

		return errors.New("ENV Variables are not set properly")
	}

	s, err := discordgo.New("Bot " + sb.botToken)

	if err != nil {
		log.Println(err.Error())
	}

	s.AddHandler(sb.message)

	s.Identify.Intents = discordgo.IntentDirectMessages

	err = s.Open()

	if err != nil {
		return errors.New("error opening Discord session for stop bot")
	}

	log.Println("Stop Bot is now running")
	return nil
}

func RunStopBot() {
	sb, err := initStopBot()

	if err != nil {
		log.Fatal(err.Error())
	}

	sb.run()
}
