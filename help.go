package hdsbrute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DisplayHelp displays help about the provided command
func DisplayHelp(b *Brute, m *discordgo.MessageCreate, c []*Command) {
	if c == nil || len(c) == 0 {
		message := "Supported commands:\n"
		for _, cmd := range b.Commands {
			message = fmt.Sprintf("%s\n`%s%s`%s", message, b.Prefix, cmd.Cmd[0], getAliases(cmd.Cmd, b.Prefix))
		}

		message = fmt.Sprintf("%s\n\nUse `%shelp [command]` for more info about the concrete command", message, b.Prefix)

		_, err := b.Session.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
		}

		return
	}

	for _, cmd := range c {
		if cmd != nil {
			DisplayCommandHelp(b, m, cmd)
		}
	}
}

// DisplayCommandHelp displays the help for a concrete command
func DisplayCommandHelp(b *Brute, m *discordgo.MessageCreate, cmd *Command) {
	if cmd.HelpFunc != nil {
		cmd.HelpFunc(b, b.Session, m)
		return
	}

	if cmd.HelpStr != "" {
		_, err := b.Session.ChannelMessageSend(m.ChannelID, cmd.HelpStr)
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
		}
		return
	}

	_, err := b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No usage info defined for `%s`. Here - grab some beers while waiting for Mak to add them :beers:", cmd.Cmd))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}

func getAliases(cmd []string, prefix string) string {
	if len(cmd) > 1 {
		var aliases []string
		for _, c := range cmd[1:] {
			aliases = append(aliases, fmt.Sprintf("%s%s", prefix, c))
		}
		return fmt.Sprintf(" (Aliases: `%s`)", strings.Join(aliases, ", "))
	}
	return ""
}
