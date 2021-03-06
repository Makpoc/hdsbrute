package hdsbrute

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hades-api/sheet/models"
)

// contains checks if a set of strings contains given value
func Contains(set []string, val string) bool {
	for _, c := range set {
		if val == c {
			return true
		}
	}
	return false
}

// GetEnvPropOrDefault looks for an environment variable for the given key. If such is found - it returns it, otherwise it returns the provided default value
func GetEnvPropOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// GetGuildMembers returns all members, that have the given role(s). If roles list is empty - it returns all members
func GetMembersByRole(s *discordgo.Session, channelId string, roles []string) ([]*discordgo.Member, error) {
	guild, err := GetGuild(s, channelId)
	if err != nil {
		return nil, err
	}

	var corpMembers []*discordgo.Member

	for _, member := range guild.Members {
		for _, role := range member.Roles {
			if isAllowedRole(getRoleName(guild, role), roles) {
				corpMembers = append(corpMembers, member)
				break // role loop
			}
		}
	}

	return corpMembers, nil
}

// isAllowedRole checks if the provided role is part of the list of allowed roles
func isAllowedRole(roleName string, allowedRoles []string) bool {
	if len(allowedRoles) == 0 {
		// no specific role requirements
		return true
	}
	for _, r := range allowedRoles {
		if r == roleName {
			return true
		}
	}

	return false
}

// getRoleName gets the role name from the role ID
func getRoleName(guild *discordgo.Guild, roleId string) string {
	for _, role := range guild.Roles {
		if role.ID == roleId {
			return role.Name
		}
	}
	return ""
}

// getRoleId gets the role ID from the role name
func getRoleId(guild *discordgo.Guild, roleName string) string {
	for _, role := range guild.Roles {
		if role.Name == roleName {
			return role.ID
		}
	}
	return ""
}

func GetGuild(s *discordgo.Session, channelId string) (*discordgo.Guild, error) {
	channel, err := s.Channel(channelId)
	if err != nil {
		return nil, err
	}
	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return nil, err
	}

	return guild, nil
}

func GetDiscordUser(s *discordgo.Session, m *discordgo.MessageCreate, userName string) (*discordgo.User, error) {
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
		if strings.ToLower(member.User.Username) == strings.ToLower(userName) || member.User.Mention() == userName {
			return member.User, nil
		}
	}
	return nil, fmt.Errorf("failed to find user: %s", userName)
}

func TrimMentionPrefix(mention string) string {
	return strings.TrimSpace(strings.TrimRight(strings.TrimLeft(strings.TrimLeft(mention, "<@!"), "<@"), ">"))
}

func IsAvailable(u *models.UserTime) bool {
	avails := u.Availability
	now := u.CurrentTime
	nowH := now.Hour()
	nowM := now.Minute()
	for _, avail := range avails {
		fromH := int(avail.From.Hours())
		fromM := int(avail.From.Minutes()) - fromH*60
		toH := int(avail.To.Hours())
		toM := int(avail.To.Minutes()) - toH*60
		if fromH > toH {
			// crosses midnight
			// e.g. 23-2
			if nowH > fromH || nowH < toH {
				return true
			}
			if (nowH == fromH && nowM >= fromM) || (nowH == toH && nowM > toM) {
				return true
			}
		} else {
			if nowH > fromH && nowH < toH {
				return true
			}
			if nowH == fromH && nowM >= fromM && (nowH < toH || nowH == toH && nowM <= toM) {
				return true
			}
		}
	}
	return false
}
