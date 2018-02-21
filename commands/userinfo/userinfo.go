package userinfo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hades-api/sheet/models"
	"github.com/makpoc/hdsbrute"
)

var backendURL string
var backendSecret string

const cmd = "info"

// UserInfoCommand ...
var UserInfoCommand = hdsbrute.Command{
	Cmd:      cmd,
	HelpFunc: helpFunc,
	Init: func(brute *hdsbrute.Brute) error {
		backendSecret = brute.Config.Secret
		backendURL = brute.Config.BackendURL

		log.Println("UserInfo initialized")
		return nil
	},
	Exec: handleFunc,
}

// helpFunc is the function called to display help/usage info
func helpFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := []string{
		"**Description**:",
		fmt.Sprintf("Prints ships and modules info taken from the `%ssheet`", b.Prefix),
		"",
		"**Usage**:",
		"",
		fmt.Sprintf("`%s%s [username|mention]`", b.Prefix, cmd),
	}
	s.ChannelMessageSend(m.ChannelID, strings.Join(helpMessage, "\n"))
}

// handleFunc handles requests for the tz command
func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	if len(query) == 0 {
		helpFunc(b, s, m)
		return
	}

	url := fmt.Sprintf("%s/api/v1/users", backendURL)
	if backendSecret != "" {
		url = fmt.Sprintf("%s?secret=%s", url, backendSecret)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get map. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get user info"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get map. Status was: %v\n", resp.Status)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get user info"))
		return
	}

	var users models.Users
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		log.Printf("Failed to decode UserInfo response. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get user info - %s", err.Error()))
		return
	}

	// TODO: - implement as query parameter to the backend!
	userArg := strings.Join(query[0:], " ")
	user, err := getUserFromArg(s, m, userArg)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user: %s", userArg))
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
		}
		return
	}

	for _, u := range users {
		if strings.ToLower(user.Username) == strings.ToLower(u.Name) {
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(u, user.AvatarURL("")))
			if err != nil {
				log.Printf("Failed to send TimeZones message: %v\n", err)
			}
			return
		}
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user in sheet database: %s", userArg))
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
	}
}

func getUserFromArg(s *discordgo.Session, m *discordgo.MessageCreate, userArg string) (*discordgo.User, error) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}
	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return nil, err
	}
	members := guild.Members

	for _, member := range members {
		if strings.ToLower(member.User.Username) == strings.ToLower(userArg) || member.User.ID == strings.TrimSuffix(strings.TrimPrefix(userArg, "<@"), ">") {
			return member.User, nil
		}
	}
	return nil, fmt.Errorf("failed to find user: %s", userArg)
}

func createEmbed(u models.User, avatarURL string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00,
		Title:  fmt.Sprintf("UserInfo for %s", u.Name),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Battleship",
				Value:  formatBSInfo(u),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Transporter",
				Value:  formatTSInfo(u),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Miner",
				Value:  formatMinerInfo(u),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Current time",
				Value:  u.TZ.CurrentTime.Format(time.Kitchen),
				Inline: true,
			},
		},
	}

	if avatarURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		}
	}
	return embed
}

func formatBSInfo(u models.User) string {
	info := []string{}
	info = append(info, fmt.Sprintf("**Role**: %s", u.BsRole))
	info = append(info, "**Modules**:")
	info = append(info, formatModulesInfo(u.BsModules))

	return strings.Join(info, "\n")
}

func formatTSInfo(u models.User) string {
	info := []string{}
	info = append(info, fmt.Sprintf("**Capacity**: %s", u.TsCapacity))
	info = append(info, "**Modules**:")
	info = append(info, formatModulesInfo(u.TsModules))

	return strings.Join(info, "\n")
}

func formatMinerInfo(u models.User) string {
	info := []string{}
	info = append(info, fmt.Sprintf("**Level**: %s", u.MinerLevel))
	info = append(info, "**Modules**:")
	info = append(info, formatModulesInfo(u.MinerModules))

	return strings.Join(info, "\n")
}

func formatModulesInfo(modules models.Modules) string {
	if modules == nil {
		return ""
	}

	result := []string{}
	for _, m := range modules {
		result = append(result, fmt.Sprintf("%s - %s", m.Name, m.Level))
	}
	return strings.Join(result, "\n")
}
