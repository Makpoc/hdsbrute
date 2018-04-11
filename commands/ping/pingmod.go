package ping

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
	"github.com/makpoc/hdsbrute/commands"
)

var memberRoles []string
var userAPI commands.UserAPI

const cmd = "pingmod"

// UserInfoCommand ...
var Command = hdsbrute.Command{
	Cmd:     cmd,
	HelpStr: "TODO",
	Init: func(brute *hdsbrute.Brute) error {
		userAPI = commands.NewUserApi(brute.Config.BackendURL, brute.Config.Secret)
		envRoles, ok := os.LookupEnv("MEMBER_ROLES")
		if ok {
			for _, role := range strings.Split(envRoles, ",") {
				memberRoles = append(memberRoles, role)
			}
		}

		log.Println("UserInfo initialized")
		return nil
	},
	Exec: handleFunc,
}

func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	if len(query) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usage: %s%s @role", b.Prefix, cmd))
		return
	}

	backendUsers, err := userAPI.GetUsers()
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to get users from sheet", err.Error()))
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
		}
	}

	for _, user := range backendUsers {
		user.
	}
}
