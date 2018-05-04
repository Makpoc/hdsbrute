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

	fmt.Printf("%s reporting for duty!", u.Username)

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
		b.displayHelp(m, words[1:])
		return
	}

	cmd := b.findCommand(words)
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

func (b *Brute) findCommand(query []string) *Command {
	if len(query) == 0 {
		return nil
	}
	var q = strings.TrimPrefix(query[0], b.Prefix)
	for _, c := range b.Commands {
		for _, cAndAliases := range c.Cmd {
			if strings.ToLower(q) == strings.ToLower(cAndAliases) {
				return &c
			}
		}
	}
	return nil
}

func (b *Brute) displayHelp(m *discordgo.MessageCreate, query []string) {
	var commands []*Command
	cmd := b.findCommand(query)
	if cmd != nil {
		commands = append(commands, cmd)
	}
	DisplayHelp(b, m, commands)
}

// ready will be called when the bot receives the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	go func() {
		var statuses = []string{
			"Hades' Star",
			"RS with bots",
			"RS, killing cerb",
			"cards with colossi",
			"shipments delivery",
			"shipments delivery",
			"shipments delivery",
			"but low on hydro",
			"with TM variations",
			"with sand on Mars",
		}
		rand.Seed(time.Now().Unix())
		for {
			err := s.UpdateStatus(0, statuses[rand.Intn(len(statuses))])
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			time.Sleep(20 * time.Minute)
		}
	}()
}
