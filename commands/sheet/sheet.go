package sheet

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

const (
	sheetEnvKey = "SHEET"
)

// SheetCommand is the command to get the SG Sheet link
var SheetCommand = hdsbrute.Command{
	Cmd:     "sheet",
	HelpStr: "Provides the link for the google spreadsheet for Star Grazers",
	Exec:    handleFunc,
}

func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	sheetLink, ok := os.LookupEnv(sheetEnvKey)
	if !ok {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sheet link not set in environment :poop:"))
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Here you go **%s**: %s", m.Author.Username, sheetLink))
	if err != nil {
		fmt.Println(err)
	}
}
