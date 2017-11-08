package hook

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"

	"github.com/andersfylling/rolerboler/bot/state"
	"github.com/s1kx/unison"
	"github.com/s1kx/unison/events"
)

var NewMemberHook = &unison.EventHook{
	Name:        "newMember",
	Description: "Add the member role to new users",
	OnEvent:     unison.EventHandlerFunc(newMemberAction),
	Events: []events.EventType{
		events.GuildMemberAddEvent,
	},
}

var roles map[string]string

func init() {
	roles = make(map[string]string)
}

func newMemberAction(ctx *unison.Context, ev *events.DiscordEvent) (handled bool, err error) {
	var m *discordgo.Member
	// Check event type
	switch e := ev.Event.(type) {
	case *discordgo.GuildMemberAdd:
		m = e.Member
	default:
		return false, nil
	}

	// decide how to react depending on bot state (pause, normal, etc.)
	if state.GetState(m.GuildID) == state.Normal {
		var roleID string
		// check if cached
		if val, ok := roles[m.GuildID]; ok {
			roleID = val
		} else {
			// go through each role to find one called `member`
			// TODO: make this setable using a command
			rolies, err := ctx.Bot.Discord.GuildRoles(m.GuildID)
			if err != nil {
				logrus.Error(err.Error())
				return false, err
			}
			for _, role := range rolies {
				if strings.ToLower(role.Name) == "member" {
					roleID = role.ID
					roles[m.GuildID] = roleID
					break
				}
			}
		}

		if roleID == "" {
			return false, errors.New("Unable to find a member role")
		}

		err := ctx.Bot.Discord.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
		if err != nil {
			return false, err
		}
	} else if state.GetState(m.GuildID) == state.Pause {
		// paused, don't add roles
	} else {
		// smth else, just don't add roles
	}

	return true, nil
}
