package wsmap

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

// helpFunc is the function called to display help/usage info
func helpFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := buildHelpString(fmt.Sprintf("%s%s", b.Prefix, mapCmd))
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}

// buildHelpString builds the help and usage string for the map
func buildHelpString(cmdWithPrefix string) string {
	var commands []string
	commands = append(commands, fmt.Sprintf("**%s** displays the map for the current WS with coordinates overlayed on top and optional sector highlighting.", cmdWithPrefix))
	commands = append(commands, "")
	commands = append(commands, "**Usage:**")
	commands = append(commands, fmt.Sprintf("``%s [[color] [coord...]...] message``", cmdWithPrefix))
	commands = append(commands, "")
	commands = append(commands, "**Examples:**")
	commands = append(commands, fmt.Sprintf("``%s green d2 @player1 @player2 form defensive line``", cmdWithPrefix))
	commands = append(commands, fmt.Sprintf("``%s orange e1 e2 e3 red d3 @squadA attack @miner1 Void D3 to take out teleport nodes``", cmdWithPrefix))
	commands = append(commands, "")
	commands = append(commands, "**Color Codes:**")
	commands = append(commands, "Green - Defense")
	commands = append(commands, "Orange - Offense")
	commands = append(commands, "Pink - Mining")
	commands = append(commands, "Red - Void")
	commands = append(commands, "Yellow - Crunch")
	commands = append(commands, "Warn - RED ALERT!!!")

	finalMessage := strings.Join(commands, "\n")
	return strings.TrimRight(finalMessage, "\n")
}
