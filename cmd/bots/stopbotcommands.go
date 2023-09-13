package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	integerOptionMinValue = 1.0
	commands              = []*discordgo.ApplicationCommand{
		{
			Name:        "stop",
			Description: "Command to put annoying idiots in timeout",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Timeout duration in minutes",
					MinValue:    &integerOptionMinValue,
					MaxValue:    60,
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			msgformat := "You learned how to use command options! " +
				"Take a look at the value(s) you entered:\n"

			var user *discordgo.User
			var duration int64

			if opt, ok := optionMap["user"]; ok {
				user = opt.UserValue(nil)
				margs = append(margs, opt.UserValue(nil).ID)
				msgformat += "> user: <@%s>\n"
			}

			if opt, ok := optionMap["duration"]; ok {
				duration = opt.IntValue()
				margs = append(margs, opt.IntValue())
				msgformat += "> duration: %d\n"
			} else {
				duration = 5
			}

			log.Println(user)
			log.Println(duration)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
	}
)
