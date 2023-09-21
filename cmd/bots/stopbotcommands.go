package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
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

		go stopUser(e, u, d)
	}
}

// Stops the User
func stopUser(e *events.ApplicationCommandInteractionCreate,
	u discord.User, d int) {

	dur, _ := time.ParseDuration(fmt.Sprintf("%dm", d))
	m, err := e.Client().Rest().GetMember(
		*e.GuildID(),
		u.ID,
	)

	if err != nil {
		e.Client().Logger().Error("error retrieving member to stop: ", err)
	}

	stoppedRole, err := snowflake.Parse(stopBot.stoppedRoleID)
	if err != nil {
		e.Client().Logger().Error("error retrieving stopped role: ", err)
	}

	// Slightly confusing and seemingly unnecessary pointer usage in this library
	mute := true
	update := discord.MemberUpdate{
		Mute: &mute,
	}

	// Check if User is already stopped, if so, we don't modify nickname
	if !strings.HasSuffix(*m.Nick, " [Stopped]") {
		nick := *m.Nick + " [Stopped]"
		update.Nick = &nick
	}

	err = e.Client().Rest().AddMemberRole(*e.GuildID(), u.ID, stoppedRole)
	if err != nil {
		e.Client().Logger().Error("error setting stopped role: ", err)
	}

	_, err = e.Client().Rest().UpdateMember(*e.GuildID(), u.ID, update)
	if err != nil {
		if strings.HasPrefix(err.Error(), "40032:") {
			_, err = e.Client().Rest().UpdateMember(*e.GuildID(), u.ID, discord.MemberUpdate{
				Nick: update.Nick,
			})
			if err != nil {
				e.Client().Logger().Error("error updating member: ", err)
			}
		} else {
			e.Client().Logger().Error("error updating member: ", err)
		}
	}
	_, err = e.Client().Rest().CreateMessage(
		e.Channel().ID(),
		discord.NewMessageCreateBuilder().SetContent("<@"+u.ID.String()+
			"> has been put in timeout for "+fmt.Sprintf("%d", d)+" minutes").Build(),
	)
	if err != nil {
		e.Client().Logger().Error("error creating message: ", err)
	}

	// Waits the duration the command was given
	// Considered Channels for this, but I don't see
	// the benefit compared to simple wait
	time.Sleep(dur)

	err = e.Client().Rest().RemoveMemberRole(*e.GuildID(), u.ID, stoppedRole)
	if err != nil {
		e.Client().Logger().Error("error unsetting stopped role: ", err)
	}

	// Super annoying
	mute = false
	nick := strings.TrimSuffix(*m.Nick, " [Stopped]")
	update.Mute = &mute
	update.Nick = &nick

	_, err = e.Client().Rest().UpdateMember(*e.GuildID(), u.ID, update)
	if err != nil {
		if strings.HasPrefix(err.Error(), "40032:") {
			_, err = e.Client().Rest().UpdateMember(*e.GuildID(), u.ID, discord.MemberUpdate{
				Nick: update.Nick,
			})
			if err != nil {
				e.Client().Logger().Error("error updating member: ", err)
			}
		} else {
			e.Client().Logger().Error("error updating member: ", err)
		}
	}

	_, err = e.Client().Rest().CreateMessage(
		e.Channel().ID(),
		discord.NewMessageCreateBuilder().SetContent("<@"+u.ID.String()+
			"> is now unmuted. Try to be less annoying please...").Build(),
	)
	if err != nil {
		e.Client().Logger().Error("error creating message: ", err)
	}
}
