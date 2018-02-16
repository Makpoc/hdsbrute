package sheet

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

const sheetEnvKey = "sheet"

// SheetCommand is frosty, duh (courtesy of TngB)
var SheetCommand = hdsbrute.Command{
	Cmd:     "sheet",
	HelpStr: "Provides the link for the google spreadsheet :spy:",
	Exec:    sheetFn,
}

func sheetFn(s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
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
