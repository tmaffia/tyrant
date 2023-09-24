package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var (
	integerOptionMinValue = 1
	integerOptionMaxValue = 60
	commands              = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "stop",
			Description: "Command to put annoying idiots in timeout",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "User",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "duration",
					Description: "Timeout duration in minutes",
					MinValue:    &integerOptionMinValue,
					MaxValue:    &integerOptionMaxValue,
					Required:    false,
				},
			},
		},
	}
)

func commandListener(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()
	if data.CommandName() == "stop" {
		err := e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("You have stopped a user! " +
				"They must have been annoying\n").
			Build(),
		)
		if err != nil {
			e.Client().Logger().Error("error on sending response: ", err)
		}

		u := data.User("user")
		d := data.Int("duration")
		if d == 0 {
			d = 5
		}

		go stopCommand(e, u, d)
	}
}
