package hdsbrute

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// HandlerFn is the function to handle commands
type HandlerFn func(*Brute, *discordgo.Session, *discordgo.MessageCreate, []string)

// Command represents a command
type Command struct {
	Cmd      []string
	Args     []string
	HelpStr  string
	HelpFunc func(*Brute, *discordgo.Session, *discordgo.MessageCreate)
	Init     func(*Brute) error
	Exec     HandlerFn
	Auth     []string // the list of roles that are allowed to use a command. Empty means everyone
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
