package bot

import (
	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/rolerboler/bot/cmd"
	"github.com/andersfylling/rolerboler/bot/hook"
	"github.com/s1kx/unison"
)

func Run(token, prefix, clientid string) error {
	// Create bot structure
	settings := &unison.BotSettings{
		Token:         token,
		CommandPrefix: prefix,

		Commands: []*unison.Command{
			cmd.RunCommand,
			cmd.PauseCommand,
			cmd.StateCommand,
			cmd.MissingCommand,
		},
		EventHooks: []*unison.EventHook{
			hook.NewMemberHook,
		},
		Services: []*unison.Service{},
	}

	if clientid != "" {
		addBotURL := "https://discordapp.com/api/oauth2/authorize?scope=bot&permissions=0&client_id="
		logrus.Info("click to add bot: " + addBotURL + clientid)
	}

	// Start the bot
	err := unison.RunBot(settings)
	if err != nil {
		return err
	}

	return nil
}
