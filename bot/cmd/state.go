package cmd

import (
	"errors"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/rolerboler/bot/state"
	"github.com/bwmarrin/discordgo"
	"github.com/s1kx/unison"
)

var StateCommand = &unison.Command{
	Name:        "state",
	Description: "Get the current bot state",
	Action:      StateCommandAction,
	Deactivated: false,
	Permission:  unison.NewCommandPermission(),
}

// PermissionRole Check if user has permission to deal with roles
func PermissionRole(ctx *unison.Context, m *discordgo.Message) bool {
	authorPermissions, _ := ctx.Bot.Discord.UserChannelPermissions(m.Author.ID, m.ChannelID)

	return (authorPermissions & discordgo.PermissionManageRoles) > 0
}

func StateCommandAction(ctx *unison.Context, m *discordgo.Message, content string) error {

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
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is Pause")
	} else if state.GetState(g.ID) == state.Normal {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is Normal.")
	} else {
		ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, "Bot state is not familiar("+strconv.Itoa(int(state.GetState(g.ID)))+").")
	}

	return err
}
