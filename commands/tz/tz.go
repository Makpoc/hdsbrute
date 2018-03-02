package tz

import (
	"encoding/json"
	"fmt"
	"io"
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

const cmd = "tz"

// TimeZoneCommand ...
var TimeZoneCommand = hdsbrute.Command{
	Cmd:      cmd,
	HelpFunc: helpFunc,
	Init: func(brute *hdsbrute.Brute) error {
		backendSecret = brute.Config.Secret
		backendURL = brute.Config.BackendURL

		log.Println("TimeZones initialized")
		return nil
	},
	Exec: handleFunc,
}

// helpFunc is the function called to display help/usage info
func helpFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := []string{
		"**Description**:",
		fmt.Sprintf("Shows the current time and time offset taken from the `%ssheet`", b.Prefix),
		"",
		"**Usage**:",
		"",
		fmt.Sprintf("`%s%s` - lists tz for all users", b.Prefix, cmd),
		fmt.Sprintf("`%s%s [username|mention]` - shows tz for the provided user", b.Prefix, cmd),
	}
	s.ChannelMessageSend(m.ChannelID, strings.Join(helpMessage, "\n"))
}

// handleFunc handles requests for the tz command
func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	isSingleUserRequest := len(query) > 0

	url := fmt.Sprintf("%s/api/v1/timezones", backendURL)
	if isSingleUserRequest {
		url = fmt.Sprintf("%s/%s", url, strings.Join(query, " "))
	}

	if backendSecret != "" {
		url = fmt.Sprintf("%s?secret=%s", url, backendSecret)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get time zones. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get timezones"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get time zones. Status was: %v\n", resp.Status)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get timezones"))
		return
	}

	if isSingleUserRequest {
		handleUserTzFunc(s, m, query, resp.Body)
	} else {
		handleAllTzFunc(s, m, query, resp.Body)
	}
}

func handleAllTzFunc(s *discordgo.Session, m *discordgo.MessageCreate, query []string, body io.Reader) {
	var timeZones []models.UserTime
	err := json.NewDecoder(body).Decode(&timeZones)
	if err != nil {
		log.Printf("Failed to decode TZ response. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get time zones - %s", err.Error()))
		return
	}

	messages := formatAllTzMessage(timeZones, 2000)
	for _, message := range messages {
		if len(message) == 0 {
			continue
		}
		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			log.Printf("Failed to send TimeZones message: %v\n", err)
		}
	}
}

func handleUserTzFunc(s *discordgo.Session, m *discordgo.MessageCreate, query []string, body io.Reader) {
	var timeZone = &models.UserTime{}
	err := json.NewDecoder(body).Decode(timeZone)
	if err != nil {
		if err == io.EOF {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":thinking: No such user: %s. Are you sure s/he is in the sheet?", strings.Join(query, " ")))
			return
		}
		log.Printf("Failed to decode TZ response. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get time zone for user - %s", err.Error()))
		return
	}

	if timeZone == nil || timeZone.UserName == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user in TZ database: %s", strings.Join(query, " ")))
		if err != nil {
			log.Printf("Failed to send message: %v\n", err)
		}
	}

	var userAvatarURL string
	discordUser, err := getUserFromArg(s, m, timeZone.UserName)
	if err != nil {
	} else {
		userAvatarURL = discordUser.AvatarURL("")
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, createTzEmbed(timeZone, userAvatarURL))
	if err != nil {
		log.Printf("Failed to send TimeZones message: %v\n", err)
		return
	}
}

func createTzEmbed(u *models.UserTime, avatarURL string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  getEmbedColor(u),
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

// getEmbedColor calculates the color based on the time of the day for the provided user
func getEmbedColor(u *models.UserTime) int {
	uHour := u.CurrentTime.Hour()

	// evening (22:00-23:59] or morning (7:00-9:59]
	if (uHour >= 22) || (uHour > 7 && uHour <= 9) {
		return 0xff9900 // orange
	}
	// night - (0:00-6:59]
	if uHour >= 0 && uHour <= 6 {
		return 0xff0000 // red
	}

	// day 10:00-21:59
	return 0x00ff00 // green
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
