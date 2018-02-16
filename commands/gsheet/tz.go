package gsheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hades-sheet/models"
	"github.com/makpoc/hdsbrute"
)

var backendURL string
var backendSecret string

// TimeZoneCommand ...
var TimeZoneCommand = hdsbrute.Command{
	Cmd:     "tz",
	HelpStr: "TODO",
	Init: func() error {
		backendSecret = hdsbrute.GetEnvPropOrDefault("secret", "")
		backendURL = hdsbrute.GetEnvPropOrDefault("tzURL", "http://localhost:3000")

		fmt.Println("TimeZones initialized")
		return nil
	},
	Exec: handlerFunc,
}

func handlerFunc(s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	url := fmt.Sprintf("%s/api/v1/timezones", backendURL)
	if backendSecret != "" {
		url = fmt.Sprintf("%s&secret=%s", url, backendSecret)
	}

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get map - got %s. Error was: %v\n", resp.Status, err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get timezones - %s", resp.Status))
		return
	}
	defer resp.Body.Close()

	var timeZones []models.UserTime
	err = json.NewDecoder(resp.Body).Decode(&timeZones)
	if err != nil {
		fmt.Printf("Failed to decode TZ response. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get time zones - %s", err.Error()))
		return
	}

	if len(query) == 0 {
		messages := formatAllTzMessage(timeZones, 2000)
		for _, message := range messages {
			if len(message) == 0 {
				continue
			}
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Printf("Failed to send TimeZones message: %v\n", err)
			}
		}
		return
	}

	// TODO: - implement as query parameter to the backend!
	userArg := strings.Join(query[0:], " ")
	user, err := getUserFromArg(s, m, userArg)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user: %s", userArg))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
		}
	}

	for _, u := range timeZones {
		if strings.ToLower(user.Username) == strings.ToLower(u.UserName) {
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, createTzEmbed(u, user.AvatarURL("")))
			if err != nil {
				fmt.Printf("Failed to send TimeZones message: %v\n", err)
			}
			return
		}
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user in TZ database: %s", userArg))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}

func createTzEmbed(u models.UserTime, avatarURL string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Title:  fmt.Sprintf("TimeZone for %s", u.UserName),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Current time",
				Value:  u.CurrentTime.Format(time.Kitchen),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Offset",
				Value:  fmt.Sprintf("%s", u.Offset),
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

func formatAllTzMessage(tz []models.UserTime, maxChars int) []string {
	var mSize int
	var result []string
	var singleMessage []string
	for _, u := range tz {
		m := formatTz(u)
		mSize += len(m)
		if mSize > maxChars {
			if len(singleMessage) > 0 {
				result = append(result, strings.Join(singleMessage, ""))
				singleMessage = []string{}
			}
			mSize = 0
			continue
		}
		singleMessage = append(singleMessage, m)
	}
	result = append(result, strings.Join(singleMessage, ""))

	return result
}

func formatTz(u models.UserTime) string {
	return fmt.Sprintf(`__**%s**__:
**CurrentTime**: %s
**Offset**: %s

`, u.UserName, u.CurrentTime.Format(time.Kitchen), u.Offset)
}
