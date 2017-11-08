package cmd

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/s1kx/unison"
)

var MissingCommand = &unison.Command{
	Name:        "missing",
	Description: "check how many is missing a role",
	Action:      MissingCommandAction,
	Deactivated: false,
	Permission:  unison.NewCommandPermission(),
}

func MissingCommandAction(ctx *unison.Context, m *discordgo.Message, content string) error {
	missingRole := strings.ToLower(strings.TrimSpace(content))
	didNotAddRoleToMemberSum := 0

	// if theres more than one role specified, quit
	if strings.ContainsAny(missingRole, " ") {
		err := errors.New("Given role had more than one word: `" + missingRole + "`")
		return err
	}

	if !PermissionRole(ctx, m) {
		err := errors.New("Member " + m.Author.Username + " do not have role permissions")
		return err
	}

	// Find the channel that the message came from.
	c, err := ctx.Discord.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	// check that role argument exists as a guild role
	validRoleArg := false
	missingRoleID := ""
	guildRoles, err := ctx.Bot.Discord.GuildRoles(c.GuildID)
	for _, guildRole := range guildRoles {
		if strings.ToLower(guildRole.Name) == missingRole {
			missingRoleID = guildRole.ID
			validRoleArg = true
			break
		}
	}
	if !validRoleArg {
		errMsg := "Argument given is a none existant role on the server: `" + missingRole + "`"
		return errors.New(errMsg)
	}

	// go through every member...
	done := false
	afterMemberID := ""
	for !done {
		logrus.Info("Getting up to 1000 members")
		members, err := ctx.Bot.Discord.GuildMembers(c.GuildID, afterMemberID, 1000)
		logrus.Info("Got " + strconv.Itoa(len(members)) + " members")
		if err != nil {
			return err
		}

		// go through every member to check for missing roles
		for _, member := range members {
			logrus.Info("Checking roles of member `" + member.User.String() + "`")
			// if the member does not have the role, add it
			lacking := true
			for _, role := range member.Roles {
				if role == missingRoleID {
					lacking = false
					break
				}
			}

			// if missing the role, add it
			if lacking {
				logrus.Info("Adding role to member `" + member.User.String() + "`")
				err := ctx.Bot.Discord.GuildMemberRoleAdd(c.GuildID, member.User.ID, missingRoleID)
				if err != nil {
					logrus.Error("Could not add role `" + missingRole + "` with roleID `" + missingRoleID + "` to member `" + member.User.String() + "` because of error: " + err.Error())
					logrus.Error("guildID `" + c.GuildID + "`, userID `" + member.User.ID + "`")
					didNotAddRoleToMemberSum++
				}
			}

			// check the ratelimit and add a needed timeout per request
			bucket := ctx.Bot.Discord.Ratelimiter.GetBucket(discordgo.EndpointGuildMemberRole(member.GuildID, "", ""))
			if bucket.Remaining == 0 {
				logrus.Info("Sleeping..")
				time.Sleep(ctx.Bot.Discord.Ratelimiter.GetWaitTime(bucket, bucket.Remaining))
				logrus.Info("Awake!")
			}

		}

		// check if there are any more left
		if len(members) < 1000 {
			done = true
		}
	}

	response := "Added role to remaining members."
	if didNotAddRoleToMemberSum > 0 {
		response += " But " + strconv.Itoa(didNotAddRoleToMemberSum) + " members did not get the role due to errors."
	}
	ctx.Bot.Discord.ChannelMessageSend(m.ChannelID, response)

	return nil
}
