package hdsbrute

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// DisplayHelp displays help about the provided command
func DisplayHelp(b *Brute, s *discordgo.Session, m *discordgo.MessageCreate, c []*Command) {
	if c == nil || len(c) == 0 {
		message := "Supported commands:\n"
		for _, c := range b.Commands {
			message = fmt.Sprintf("%s\n`%s%s`", message, b.Prefix, c.Cmd)
		}

		message = fmt.Sprintf("%s\n\nUse `%shelp [command]` for more info about the concrete command", message, b.Prefix)

		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	for _, cmd := range c {
		DisplayCommandHelp(b, s, m, cmd)
	}
}

// DisplayCommandHelp displays the help for a concrete command
func DisplayCommandHelp(b *Brute, s *discordgo.Session, m *discordgo.MessageCreate, cmd *Command) {
	if cmd.HelpFunc != nil {
		cmd.HelpFunc(b, s, m)
		return
	}

	if cmd.HelpStr == "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No usage info defined for `%s`. Here - grab some beers while waiting for Mak to add them :beers:", cmd.Cmd))
	}

	s.ChannelMessageSend(m.ChannelID, cmd.HelpStr)
}
