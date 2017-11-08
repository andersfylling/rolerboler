package cmd

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/rolerboler/bot/state"
	"github.com/bwmarrin/discordgo"
	"github.com/s1kx/unison"
)

var PauseCommand = &unison.Command{
	Name:        "pause",
	Description: "Get the number of servers using this bot",
	Action:      PauseCommandAction,
	Deactivated: false,
	Permission:  unison.NewCommandPermission(),
}

func PauseCommandAction(ctx *unison.Context, m *discordgo.Message, content string) error {
	if !PermissionRole(ctx, m) {
		return errors.New("Member " + m.Author.Username + " do not have role permissions")
	}
	// Find the channel that the message came from.
	c, err := ctx.Discord.State.Channel(m.ChannelID)
	if err != nil {
		logrus.Error("Could not find the massage's Channel")
		return nil
	}

	// Find the guild for that channel.
	g, err := ctx.Discord.State.Guild(c.GuildID)
	if err != nil {
		logrus.Error("Could not find the Channel's Guild")
		return nil
	}

	if state.GetState(g.ID) == state.Pause {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Cannot pause bot when already paused")
		return nil
	}

	state.SetState(g.ID, state.Pause)
	ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Setting state to pause..")
	if state.GetState(g.ID) == state.Pause {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is now Pause")
	} else {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is not pause.")
	}
	return nil
}
