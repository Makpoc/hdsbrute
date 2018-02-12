package hdsbrute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Initer is an interface with Init method
type Initer interface {
	Init()
}

// Brute ...
type Brute struct {
	Prefix   string
	BotID    string
	Session  *discordgo.Session
	Commands []Command
}

// New ...
func New(prefix, token string) (*Brute, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	var brute = &Brute{}

	brute.Session = s
	brute.Prefix = prefix
	brute.Commands = []Command{}

	u, err := s.User("@me")
	if err != nil {
		return nil, fmt.Errorf("failed to query @me: %v", err)
	}

	brute.BotID = u.ID
	return brute, nil
}

// AddCommand ...
func (b *Brute) AddCommand(cmd Command) {
	err := cmd.Init()
	if err != nil {
		fmt.Printf("failed to initialize command %s: %v\n", cmd.Cmd, err)
		return
	}
	b.Commands = append(b.Commands, cmd)
	fmt.Printf("Added %s to commands list\n", cmd.Cmd)
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
		cmd.Exec(b.Session, m, words[1:])
	}
}

// Start ...
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
