package hdsbrute

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Brute ...
type Brute struct {
	Prefix   string
	BotID    string
	Session  *discordgo.Session
	Commands []Command
	Config   *Config
}

// New ...
func New(token string) (*Brute, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	var brute = &Brute{}
	brute.loadConfig()

	brute.Session = s
	brute.Prefix = brute.Config.BotPrefix
	brute.Commands = []Command{}

	u, err := s.User("@me")
	if err != nil {
		return nil, fmt.Errorf("failed to query @me: %v", err)
	}

	brute.BotID = u.ID

	s.AddHandler(ready)

	return brute, nil
}

// loadConfig loads the configuration from the environment into the brute struct
func (b *Brute) loadConfig() {
	b.Config = &Config{
		BackendURL: GetEnvPropOrDefault("BACKEND_URL", "http://localhost"),
		Secret:     GetEnvPropOrDefault("SECRET", ""),
		BotPrefix:  GetEnvPropOrDefault("BOT_PREFIX", "."),
	}
}

// AddCommand ...
func (b *Brute) AddCommand(cmd Command) {
	if cmd.Init != nil {
		err := cmd.Init(b)
		if err != nil {
			log.Printf("failed to initialize command %s: %v\n", cmd.Cmd, err)
			return
		}
	}
	b.Commands = append(b.Commands, cmd)
	log.Printf("Added %s to commands list\n", cmd.Cmd)
}

// Dispatch ...
func (b *Brute) Dispatch(m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, b.Prefix) {
		return
	}

	if m.Author.ID == b.BotID {
		return
	}

	words := strings.Fields(m.Content)

	if len(words) == 0 {
		return
	}

	if strings.ToLower(words[0]) == b.Prefix+"help" {
		b.displayHelp(b.Session, m, words[1:])
		return
	}

	cmd := b.findCommand(strings.TrimPrefix(words[0], b.Prefix))
	if cmd != nil && cmd.Exec != nil {
		cmd.Exec(b, b.Session, m, words[1:])
	}
}

// Start starts the bot
func (b *Brute) Start() error {
	b.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		b.Dispatch(m)
	})

	err := b.Session.Open()
	if err != nil {
		return err
	}

	return nil
}

// Close closes the connection
func (b *Brute) Close() error {
	if b != nil {
		return b.Session.Close()
	}

	return nil
}

func (b *Brute) findCommand(query string) *Command {
	for _, c := range b.Commands {
		if strings.ToLower(query) == strings.ToLower(c.Cmd) {
			return &c
		}
	}
	return nil
}

func (b *Brute) displayHelp(s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	if len(query) == 0 {
		message := "Supported commands:\n"
		for _, c := range b.Commands {
			message = fmt.Sprintf("%s\n`%s%s`", message, b.Prefix, c.Cmd)
		}

		message = fmt.Sprintf("%s\n\nUse `%shelp [command]` for more info about the concrete command", message, b.Prefix)

		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	cmd := b.findCommand(query[0])
	DisplayHelp(b, s, m, []*Command{cmd})
}

// ready will be called when the bot receives the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	go func() {
		var statuses = []string{
			"Hades' Star",
			"RS with bots",
			"shipments delivery",
			"shipments delivery",
			"shipments delivery",
			"but low on hydro",
			"with TimeMachine",
		}
		rand.Seed(time.Now().Unix())
		for {
			err := s.UpdateStatus(0, statuses[rand.Intn(len(statuses))])
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
