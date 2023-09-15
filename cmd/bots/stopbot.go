package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type StopBot struct {
	botToken           string
	appID              string
	publicKey          string
	guildID            string
	stoppedRoleID      string
	intents            discordgo.Intent
	removeCommands     bool
	registeredCommands []*discordgo.ApplicationCommand
}

func initStopBot() *StopBot {
	removeCommands, err := strconv.ParseBool(os.Getenv("REMOVE_COMMANDS"))
	if err != nil {
		removeCommands = false
	}
	sb := StopBot{
		botToken:       os.Getenv("STOP_BOT_TOKEN"),
		appID:          os.Getenv("STOP_BOT_APP_ID"),
		publicKey:      os.Getenv("STOP_BOT_PUBLIC_KEY"),
		guildID:        "",
		stoppedRoleID:  os.Getenv("STOP_BOT_STOPPED_ROLE_ID"),
		intents:        discordgo.IntentDirectMessages,
		removeCommands: removeCommands,
	}
	return &sb
}

var sb *StopBot

func (sb StopBot) run() (*discordgo.Session, error) {

	if sb.botToken == "" ||
		sb.appID == "" ||
		sb.publicKey == "" {
		return nil, errors.New("ENV Variables are not set properly")
	}

	s, err := discordgo.New("Bot " + sb.botToken)
	if err != nil {
		log.Println(err.Error())
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	sb.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	err = s.Open()
	if err != nil {
		return nil, errors.New("error opening Discord session for stop bot")
	}

	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, sb.guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		sb.registeredCommands[i] = cmd
	}

	log.Println("Stop Bot is now running")
	return s, nil
}

func (sb StopBot) KillStopBot(s *discordgo.Session) {
	if sb.removeCommands {
		log.Println("Removing commands...")

		c, _ := s.ApplicationCommands(sb.appID, sb.guildID)
		for _, cmd := range c {
			log.Println("Unregistering: " + cmd.Name)
			s.ApplicationCommandDelete(s.State.User.ID, sb.guildID, cmd.ID)
		}
	}

	log.Println("Gracefully shutting down.")
}

func RunStopBot() {
	sb = initStopBot()

	s, err := sb.run()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer s.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-sc

	sb.KillStopBot(s)
}
