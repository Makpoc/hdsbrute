package userinfo

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hades-api/sheet/models"
	"github.com/makpoc/hdsbrute"
	"github.com/makpoc/hdsbrute/commands"
)

const cmd = "info"

var userAPI commands.UserAPI

// UserInfoCommand ...
var UserInfoCommand = hdsbrute.Command{
	Cmd:      []string{cmd},
	HelpFunc: helpFunc,
	Init: func(brute *hdsbrute.Brute) error {
		userAPI = commands.NewUserApi(brute.Config.BackendURL, brute.Config.Secret)

		log.Println("UserInfo initialized")
		return nil
	},
	Exec: handleFunc,
	Auth: getAllowedRoles(),
}

func getAllowedRoles() []string {
	var result []string
	envRoles, ok := os.LookupEnv("MEMBER_ROLES")
	if ok {
		for _, role := range strings.Split(envRoles, ",") {
			result = append(result, role)
		}
	}
	return result
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
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s%s` is currently disabled. Check `%s%s` for module info", b.Prefix, cmd, b.Prefix, "sheet"))
}

func createEmbed(u *models.User, avatarURL string) *discordgo.MessageEmbed {
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

func formatBSInfo(u *models.User) string {
	info := []string{}
	info = append(info, fmt.Sprintf("**Role**: %s", u.BsRole))
	info = append(info, fmt.Sprintf("**Weapon**: %s", u.BsWeapon))
	info = append(info, fmt.Sprintf("**Shield**: %s", u.BsShield))
	info = append(info, "**Modules**:")
	info = append(info, formatModulesInfo(u.BsModules))

	return strings.Join(info, "\n")
}

func formatTSInfo(u *models.User) string {
	info := []string{}
	info = append(info, fmt.Sprintf("**Capacity**: %s", u.TsCapacity))
	info = append(info, "**Modules**:")
	info = append(info, formatModulesInfo(u.TsModules))

	return strings.Join(info, "\n")
}

func formatMinerInfo(u *models.User) string {
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
