package userinfo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hades-api/sheet/models"
	"github.com/makpoc/hdsbrute"
)

var backendURL string
var backendSecret string

var memberRoles []string

const cmd = "info"

// UserInfoCommand ...
var UserInfoCommand = hdsbrute.Command{
	Cmd:      cmd,
	HelpFunc: helpFunc,
	Init: func(brute *hdsbrute.Brute) error {
		backendSecret = brute.Config.Secret
		backendURL = brute.Config.BackendURL
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

// handleFunc handles requests for the info command
func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	if len(query) == 0 {
		helpFunc(b, s, m)
		return
	}

	url := fmt.Sprintf("%s/api/v1/users", backendURL)
	if backendSecret != "" {
		url = fmt.Sprintf("%s?secret=%s", url, backendSecret)
	}

	backendUsers, err := getSheetUsers(url)
	if err != nil {
		log.Printf("Failed to get Sheet Users: %v", err)
		_, err = s.ChannelMessageSend(m.ChannelID, "Failed to get sheet users :poop:")
		return
	}

	userArg := strings.TrimSpace(strings.Join(query[0:], " "))
	// search the sheet for the provided arguments first
	if backendUser, err := getBackendUser(userArg, backendUsers); err == nil {
		avatarURL := ""
		if discordUser, err := findDiscordUser(s, m, backendUser.Name); err == nil {
			avatarURL = discordUser.AvatarURL("")
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(backendUser, avatarURL))
		if err != nil {
			log.Printf("Failed to send User Info message: %v\n", err)
		}
		return
	}

	discordUser, err := findDiscordUser(s, m, userArg)
	if err == nil {
		// matching discord user found
		if backendUser, err := getBackendUser(discordUser.Username, backendUsers); err == nil {
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(backendUser, discordUser.AvatarURL("")))
			if err != nil {
				log.Printf("Failed to send User Info message: %v\n", err)
			}
			return
		}
	} else {
		if strings.HasPrefix(err.Error(), "multiple users matched") {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ambiguous user - *%s*: %s", userArg, err.Error()))
			if err != nil {
				log.Printf("Failed to send message: %v\n", err)
			}
			return
		}
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No such user in sheet database: %s", userArg))
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
	}
}

func getBackendUser(userName string, backendUsers models.Users) (models.User, error) {
	for _, backendUser := range backendUsers {
		if strings.ToLower(userName) == strings.ToLower(backendUser.Name) {
			return backendUser, nil
		}
	}

	for _, backendUser := range backendUsers {
		if matchPartialUser(backendUser.Name, userName) {
			return backendUser, nil
		}
	}
	return models.User{}, fmt.Errorf("user not found in sheet database")
}

func getSheetUsers(url string) (models.Users, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info. Status was: %s", resp.Status)
	}

	var users models.Users
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return users, nil
}

func findDiscordUser(s *discordgo.Session, m *discordgo.MessageCreate, userArg string) (*discordgo.User, error) {
	members, err := hdsbrute.GetGuildMembers(s, m, memberRoles)
	if err != nil {
		return nil, err
	}

	var matched []*discordgo.User
	var matchedUsernames []string

	userArg = strings.TrimSpace(strings.ToLower(userArg))
	for _, member := range members {
		if matchExactUser(member.User, userArg) {
			matched = append(matched, member.User)
			matchedUsernames = append(matchedUsernames, member.User.Username)
		}
	}

	if len(matched) == 0 {
		for _, member := range members {
			if matchPartialUser(member.User.Username, userArg) {
				matched = append(matched, member.User)
				matchedUsernames = append(matchedUsernames, member.User.Username)
			}
		}
	}

	switch len(matched) {
	case 0:
		return nil, fmt.Errorf("failed to find user: %s", userArg)
	case 1:
		return matched[0], nil
	default:
		return nil, fmt.Errorf("multiple users matched %s: %s", userArg, strings.Join(matchedUsernames, ", "))
	}
}

func matchExactUser(given *discordgo.User, wanted string) bool {
	return strings.ToLower(given.Username) == wanted || given.ID == strings.TrimSuffix(strings.TrimPrefix(wanted, "<@"), ">")
}

func matchPartialUser(given string, wanted string) bool {
	return strings.Contains(strings.TrimSpace(strings.ToLower(given)), strings.ToLower(wanted))
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
				Name:   "Current time",
				Value:  u.TZ.CurrentTime.Format(time.Kitchen),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				// add a little spacing
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: false,
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
	info = append(info, fmt.Sprintf("**Weapon**: %s", u.BsWeapon))
	info = append(info, fmt.Sprintf("**Shield**: %s", u.BsShield))
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
