package hdsbrute

import (
	"fmt"

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
