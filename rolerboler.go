package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/andersfylling/rolerboler/bot"
	"github.com/urfave/cli"
)

const (
	// EnvVarPrefix is the prefix for environment variables
	EnvVarPrefix = "ROLERBOLER"
)

var DiscordToken string
var DiscordPrefix string
var DiscordClientID string

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func webServerForHeroku() {
	port := os.Getenv("PORT")
	if port == "" {
		return // not on heroku i assume....
	}

	http.HandleFunc("/", helloWorld)
	err := http.ListenAndServe(":"+port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	go webServerForHeroku()
	// Initialize command-line application
	app := &cli.App{
		Name:  "rolerboler",
		Usage: "Discord bot to deal with assigning and revoking roles",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				EnvVars: envVarNames("DEBUG"),
				Usage:   "debug mode",
			},
		},
		Before: initApplication,
		Action: runApplication,
	}
	app.Run(os.Args)
}

func envVarName(name string) string {
	return fmt.Sprintf("%s_%s", EnvVarPrefix, strings.ToUpper(name))
}

func envVarNames(names ...string) []string {
	res := make([]string, len(names))
	for i, name := range names {
		res[i] = envVarName(name)
	}
	return res
}

var logFormatter = logrus.TextFormatter{
	FullTimestamp:   true,
	TimestampFormat: "2006-01-02 15:04:05",
}

func initApplication(c *cli.Context) error {
	debug := c.Bool("debug")

	// Configure logger.
	logrus.SetFormatter(&logFormatter)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Use environment variables as main config
	token := os.Getenv("ROLERBOLER_TOKEN")
	if token != "" {
		commandPrefix := os.Getenv("ROLERBOLER_COMMANDPREFIX")
		if commandPrefix == "" {
			logrus.Fatal("No command prefix set. env var ROLERBOLER_COMMANDPREFIX")
			return nil
		}
		logrus.Debugf("Set command prefix to %s", commandPrefix)

		clientID := os.Getenv("ROLERBOLER_CLIENTID")
		if clientID == "" {
			logrus.Info("No discord client id found in env ROLERBOLER_CLIENTID")
		} else {
			logrus.Debugf("Found client id: %s", clientID)
		}

		logrus.Debugf("Found discord token: %s", token)

		DiscordToken = token
		DiscordPrefix = commandPrefix
		DiscordClientID = clientID
	} else {
		logrus.Fatal("Error could not find environment variable for discord token ROLERBOLER_TOKEN")
		return nil
	}

	return nil
}

func runApplication(c *cli.Context) error {
	err := bot.Run(DiscordToken, DiscordPrefix, DiscordClientID)
	if err != nil {
		logrus.Fatal(err)
	}

	return err
}
