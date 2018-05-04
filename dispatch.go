package hdsbrute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
		if cmd.Auth == nil || len(cmd.Auth) == 0 || hasPermission(b, m, cmd) {
			cmd.Exec(b, b.Session, m, words[1:])
		} else {
			_, err := b.Session.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command!")
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
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
