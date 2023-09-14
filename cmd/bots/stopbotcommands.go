package main

import (
	"fmt"
	"log"
	"time"

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
			msgformat := "You have Stopped a User! " +
				"They must have fucking sucked\n"

			var u *discordgo.User
			var dur int64
			if opt, ok := optionMap["user"]; ok {
				u = opt.UserValue(nil)
			}

			if opt, ok := optionMap["duration"]; ok {
				dur = opt.IntValue()
			} else {
				dur = 5
			}

			go stopUser(s, i, u, dur)

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

func stopUser(s *discordgo.Session, i *discordgo.InteractionCreate,
	u *discordgo.User, dur int64) {

	d, err := time.ParseDuration(fmt.Sprintf("%dm", dur))
	if err != nil {
		log.Println(err)
	}

	g, _ := s.Guild(i.GuildID)
	m, err := s.GuildMember(g.ID, u.ID)
	if err != nil {
		log.Println(err)
	}

	// s.GuildMemberNickname(m.GuildID, m.User.ID, "STOP "+m.Nick)
	s.ChannelMessageSend(i.ChannelID, "<@"+m.User.ID+"> has been put in timeout for "+fmt.Sprintf("%d", dur)+" minutes")
	time.Sleep(d)
	s.ChannelMessageSend(i.ChannelID, "<@"+m.User.ID+"> timeout has ended. Be less annoying next time...")
}
