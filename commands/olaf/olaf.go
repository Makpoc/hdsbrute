package olaf

import (
	"fmt"
	"log"

	"github.com/makpoc/hdsbrute"

	"github.com/bwmarrin/discordgo"
)

// Command is olaf
var Command = hdsbrute.Command{
	Cmd:     []string{"olaf"},
	HelpStr: "Olaf is the new frosty",
	Init: func(b *hdsbrute.Brute) error {
		log.Println("Olaf ready!")
		return nil
	},
	Exec: handleFunc,
}

// handleFunc responds to frosty command with some random stuff. In loving memory of TngB
func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	guild, err := hdsbrute.GetGuild(s, m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to get guild info")
		return
	}
	for _, e := range guild.Emojis {
		if e.Name == "kidding" {
			err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, e.APIName())
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			break
		}
	}
}
