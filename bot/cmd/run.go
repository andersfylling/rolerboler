package cmd

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/rolerboler/bot/state"
	"github.com/bwmarrin/discordgo"
	"github.com/s1kx/unison"
)

var RunCommand = &unison.Command{
	Name:        "run",
	Description: "Get the number of servers using this bot",
	Action:      RunCommandAction,
	Deactivated: false,
	Permission:  unison.NewCommandPermission(),
}

func RunCommandAction(ctx *unison.Context, m *discordgo.Message, content string) error {
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

	state.SetState(g.ID, state.Normal)
	ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Setting state to normal..")
	if state.GetState(g.ID) == state.Normal {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is now Normal")
	} else {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is not Normal.")
	}

	return err
}
