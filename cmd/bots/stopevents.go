package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

// Stops the User
func stopCommand(e *events.ApplicationCommandInteractionCreate,
	u discord.User, d int) {

	dur, _ := time.ParseDuration(fmt.Sprintf("%dm", d))
	targetMember, err := e.Client().Rest().GetMember(
		*e.GuildID(),
		u.ID,
	)
	if err != nil {
		e.Client().Logger().Error("error retrieving member to stop: ", err)
	}

	stopUser(e, targetMember)
	_, err = e.Client().Rest().CreateMessage(
		e.Channel().ID(),
		discord.NewMessageCreateBuilder().SetContent("<@"+u.ID.String()+
			"> has been put in timeout for "+fmt.Sprintf("%d", d)+" minutes").Build(),
	)
	if err != nil {
		e.Client().Logger().Error("error creating message: ", err)
	}

	time.Sleep(dur)

	unstopUser(e, targetMember)
	_, err = e.Client().Rest().CreateMessage(
		e.Channel().ID(),
		discord.NewMessageCreateBuilder().SetContent("<@"+u.ID.String()+
			"> is now unmuted. Try to be less annoying please...").Build(),
	)
	if err != nil {
		e.Client().Logger().Error("error creating message: ", err)
	}
}

// Plays sound effect when stopping user
func playStopSoundEffect(e *events.ApplicationCommandInteractionCreate,
	ch *snowflake.ID, closeChan chan os.Signal) {

	conn := e.Client().VoiceManager().CreateConn(*e.GuildID())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := conn.Open(ctx, *ch, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer closeCancel()
		conn.Close(closeCtx)
	}()

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}

	writeOpus(conn.UDP())
	closeChan <- syscall.SIGTERM
}

func writeOpus(w io.Writer) {
	file, err := os.Open("audio/shush-up.opus")
	if err != nil {
		panic("error opening file: " + err.Error())
	}
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenBuf [4]byte
	for range ticker.C {
		_, err = io.ReadFull(file, lenBuf[:])
		if err != nil {
			if err == io.EOF {
				_ = file.Close()
				return
			}
			panic("error reading file: " + err.Error())
		}

		// Read the integer
		frameLen := int64(binary.LittleEndian.Uint32(lenBuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, file, frameLen)
		if err != nil && err != io.EOF {
			_ = file.Close()
			return
		}
	}
}

// Initiate stopping. This will stop
func stopUser(e *events.ApplicationCommandInteractionCreate, m *discord.Member) {

	vs, connected := e.Client().Caches().VoiceState(m.GuildID, m.User.ID)
	s := make(chan os.Signal, 1)
	update := discord.MemberUpdate{}

	if connected {
		mute := true
		update.Mute = &mute
		go playStopSoundEffect(e, vs.ChannelID, s)
	}

	// Check if User is already stopped, if so, we don't modify nickname
	if !strings.HasSuffix(*m.Nick, " [Stopped]") {
		nick := *m.Nick + " [Stopped]"
		update.Nick = &nick
	}

	err := e.Client().Rest().AddMemberRole(*e.GuildID(), m.User.ID, tyr.stoppedRoleID)
	if err != nil {
		e.Client().Logger().Error("error setting stopped role: ", err)
	}

	_, err = e.Client().Rest().UpdateMember(*e.GuildID(), m.User.ID, update)
	if err != nil {
		e.Client().Logger().Error("error updating member: ", err)
	}
}

func unstopUser(e *events.ApplicationCommandInteractionCreate, m *discord.Member) {
	update := discord.MemberUpdate{}
	err := e.Client().Rest().RemoveMemberRole(*e.GuildID(), m.User.ID, tyr.stoppedRoleID)
	if err != nil {
		e.Client().Logger().Error("error unsetting stopped role: ", err)
	}

	_, connected := e.Client().Caches().VoiceState(m.GuildID, m.User.ID)

	if connected {
		mute := false
		update.Mute = &mute
	}

	nick := strings.TrimSuffix(*m.Nick, " [Stopped]")
	update.Nick = &nick

	_, err = e.Client().Rest().UpdateMember(*e.GuildID(), m.User.ID, update)
	if err != nil {
		e.Client().Logger().Error("error updating member: ", err)
	}
}
