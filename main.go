package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/brattonross/discordbot/command"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Bot ...
type Bot struct {
	Commands *command.Handler
	Config   *Config
	Discord  *discordgo.Session
}

// Config represents the configuration for the bot.
// This can be edited in the config.json file in the bot's base directory.
type Config struct {
	DiscordToken string   `json:"discordToken"`       // Discord Bot user token.
	Prefixes     []string `json:"prefixes,omitempty"` // The prefix for a command, defaults to "!" if not set.
}

// NewBot creates a new instance of the bot.
func NewBot() *Bot {
	return &Bot{
		Commands: command.NewHandler(),
		Config:   &Config{Prefixes: []string{"!"}},
	}
}

func main() {
	bot := NewBot()

	// Read and apply config
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": "./config.json",
		}).Fatal("failed to read config file")
		return
	}
	err = json.Unmarshal(b, bot.Config)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": "./config.json",
		}).Fatal("unable to unmarshal config file")
		return
	}

	// Setup commands
	bot.Commands.AddCommand("bee", &command.Bee{})

	// Create Discord session
	bot.Discord, err = discordgo.New("Bot " + bot.Config.DiscordToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to create new discordgo session")
		return
	}

	bot.Discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		log.Info("parsing message: ", m.Content)

		args := strings.Split(m.Content, " ")
		if len(args) < 1 || len(args[0]) < 1 {
			return
		}
		prefix := string(args[0][0])

		found := false
		for _, p := range bot.Config.Prefixes {
			if p == prefix {
				found = true
				break
			}
		}
		if !found {
			return
		}

		// Create a new context for this message
		ctx := &command.Context{
			Prefix:  prefix,
			Command: args[0][1:],
			Args:    args[1:],
			Message: m,
			Session: s,
		}

		if cmd := bot.Commands.Command(ctx.Command); cmd != nil {
			cmd.Execute(ctx)
		}
	})

	err = bot.Discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed open discordgo session")
		return
	}

	fmt.Println("Successfully opened Discord session!")
	fmt.Println("Press Ctrl+C to quit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Discord.Close()
}
