package sheet

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

const (
	sgSheetEnvKey = "SG_SHEET"
	ctSheetEnvKey = "CT_SHEET"
)

// SheetCommand is the command to get the SG Sheet link
var SGSheetCommand = hdsbrute.Command{
	Cmd:     "sheet",
	HelpStr: "Provides the link for the google spreadsheet for Star Grazers",
	Exec:    sgHandleFunc,
}

// SheetCommand is frosty, duh (courtesy of TngB)
var CTSheetCommand = hdsbrute.Command{
	Cmd:     "ctsheet",
	HelpStr: "Provides the link for the google spreadsheet for CometTrans",
	Exec:    ctHandleFunc,
}

func sgHandleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	handleFunc(b, s, m, query, sgSheetEnvKey)
}

func ctHandleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	handleFunc(b, s, m, query, ctSheetEnvKey)
}

func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string, sheetEnvKey string) {
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
