package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

type StopBot struct {
	botToken      string
	appID         string
	publicKey     string
	guildID       string
	stoppedRoleID string
	// intents            disgo.Intent
	removeCommands     bool
	registeredCommands []discord.ApplicationCommand
}

func initStopBot() *StopBot {
	removeCommands, err := strconv.ParseBool(os.Getenv("REMOVE_COMMANDS"))
	if err != nil {
		removeCommands = false
	}
	stopBot = StopBot{
		botToken:      os.Getenv("STOP_BOT_TOKEN"),
		appID:         os.Getenv("STOP_BOT_APP_ID"),
		publicKey:     os.Getenv("STOP_BOT_PUBLIC_KEY"),
		guildID:       "",
		stoppedRoleID: os.Getenv("STOP_BOT_STOPPED_ROLE_ID"),
		// intents:        discordgo.IntentDirectMessages,
		removeCommands: removeCommands,
	}
	return &stopBot
}

var stopBot StopBot

func (sb *StopBot) run() (bot.Client, error) {

	if sb.botToken == "" ||
		sb.appID == "" ||
		sb.publicKey == "" {
		return nil, errors.New("missing required environment variables")
	}

	client, err := disgo.New(sb.botToken,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		log.Fatal("error while creating bot client: ", err)
		return nil, err
	}

	sb.registeredCommands, _ = client.Rest().SetGlobalCommands(client.ApplicationID(), commands)

	if err != nil {
		log.Fatal("error while registering commands: ", err)
		return nil, err
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
		return nil, err
	}

	return client, nil
}

func (sb StopBot) KillStopBot(client bot.Client) {
	if sb.removeCommands {
		log.Info("Removing commands...")

		for _, c := range sb.registeredCommands {
			log.Info("Unregistering: " + c.Name())
			client.Rest().DeleteGlobalCommand(client.ApplicationID(), c.ID())
		}
	}

	log.Info("Gracefully shutting down.")
}

func RunStopBot() {
	sb := initStopBot()

	client, err := sb.run()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer client.Close(context.TODO())

	log.Infof("Stop Bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	sb.KillStopBot(client)
}
