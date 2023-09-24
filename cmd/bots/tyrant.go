package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type Tyrant struct {
	botToken           string
	appID              string
	publicKey          string
	guildID            string
	stoppedRoleID      snowflake.ID
	removeCommands     bool
	registeredCommands []discord.ApplicationCommand
	client             *bot.Client
}

var tyr Tyrant

func initTyrant() *Tyrant {
	stoppedRoleId, err := snowflake.Parse(os.Getenv("TYRANT_STOPPED_ROLE_ID"))
	if err != nil {
		log.Fatal("error getting stopped role: ", err)
	}

	removeCommands, err := strconv.ParseBool(os.Getenv("REMOVE_COMMANDS"))
	if err != nil {
		removeCommands = false
	}

	tyr = Tyrant{
		botToken:       os.Getenv("TYRANT_TOKEN"),
		appID:          os.Getenv("TYRANT_APP_ID"),
		publicKey:      os.Getenv("TYRANT_PUBLIC_KEY"),
		guildID:        "",
		stoppedRoleID:  stoppedRoleId,
		removeCommands: removeCommands,
	}
	if tyr.botToken == "" ||
		tyr.appID == "" ||
		tyr.publicKey == "" {
		panic("missing required environment variables")
	}
	return &tyr
}

func (tyrant *Tyrant) run() (bot.Client, error) {

	client, err := disgo.New(tyrant.botToken,
		bot.WithGatewayConfigOpts(gateway.WithIntents(
			gateway.IntentGuildVoiceStates,
			gateway.IntentGuildMessages,
		)),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		log.Fatal("error while creating bot client: ", err)
		return nil, err
	}

	tyrant.client = &client
	tyrant.registeredCommands, _ = client.Rest().SetGlobalCommands(client.ApplicationID(), commands)

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

func (tyrant Tyrant) KillTyrant(client bot.Client) {
	if tyrant.removeCommands {
		log.Info("Removing commands...")

		for _, c := range tyrant.registeredCommands {
			log.Info("Unregistering: " + c.Name())
			client.Rest().DeleteGlobalCommand(client.ApplicationID(), c.ID())
		}
	}

	log.Info("Gracefully shutting down.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Debug(err)
	}

	tyrant := initTyrant()

	client, err := tyrant.run()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer client.Close(context.TODO())

	log.Infof("Tyrant is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	tyrant.KillTyrant(client)
}
