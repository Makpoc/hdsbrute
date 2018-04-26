package available

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

var cmd = []string{"available", "avail", "online", "here"}

// Command for .available...
var Command = hdsbrute.Command{
	Cmd:     cmd,
	HelpStr: fmt.Sprintf("Lists all people currently available (based on their sheet configuraiton). Aliases: `%s`", strings.Join(cmd, ", ")),
	Init: func(brute *hdsbrute.Brute) error {
		backendSecret = brute.Config.Secret
		backendURL = brute.Config.BackendURL

		log.Println("Available initialized")
		return nil
	},
	Exec: handleFunc,
}

// handleFunc handles requests for the .available command
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
		log.Printf("Failed to get list of available members. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get list of available members"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get list of available members. Status was: %v\n", resp.Status)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get list of available members"))
		return
	}

	respondWithAvailables(s, m, query, resp.Body)
}

func respondWithAvailables(s *discordgo.Session, m *discordgo.MessageCreate, query []string, body io.Reader) {
	var timeZones []models.UserTime
	err := json.NewDecoder(body).Decode(&timeZones)
	if err != nil {
		log.Printf("Failed to decode TZ response. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get list of available members - %s", err.Error()))
		return
	}

	embed := prepareEmbed(timeZones)
	embedMsg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Printf("Failed to send available users message: %v\n", err)
		return
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Printf("Failed to delete .avail trigger message: %v", err)
	}
	time.AfterFunc(time.Second*30, func() {
		err := s.ChannelMessageDelete(embedMsg.ChannelID, embedMsg.ID)
		if err != nil {
			log.Printf("Failed to delete .avail embed: %v", err)
		}
	})
}

func prepareEmbed(users models.UserTimes) *discordgo.MessageEmbed {
	var availableUsers []string
	for _, u := range users {
		if hdsbrute.IsAvailable(&u) {
			availableUsers = append(availableUsers, fmt.Sprintf("- %s (%s)", u.UserName, u.CurrentTime.Format(time.Kitchen)))
		}
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Available players",
				Value: strings.Join(availableUsers, "\n"),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Note - this list does not guarantee availability",
		},
	}

	return embed
}
