package hdsbrute

import (
	"os"

	"github.com/bwmarrin/discordgo"
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
func GetGuildMembers(s *discordgo.Session, m *discordgo.MessageCreate, roles []string) ([]*discordgo.Member, error) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}
	guild, err := s.Guild(channel.GuildID)
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
