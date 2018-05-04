package hdsbrute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// hasPermissions checks if the author of the message has the necessary role to use the command
func hasPermission(b *Brute, m *discordgo.MessageCreate, cmd *Command) bool {
	if cmd.Auth == nil || len(cmd.Auth) == 0 {
		return true
	}

	fmt.Printf("Checking permissions for %s. Required Auth: %s\n", m.Author.Username, strings.Join(cmd.Auth, ","))
	g, err := GetGuild(b.Session, m.ChannelID)
	if err != nil {
		fmt.Printf("Failed to get guild: %v\n", err)
		return false
	}
	member, err := b.Session.GuildMember(g.ID, m.Author.ID)
	if err != nil {
		fmt.Printf("Failed to get guild member: %v\n", err)
		return false
	}

	for _, allowedRole := range cmd.Auth {
		// try to convert name to ID, but if that fails - use the provided string directly
		var convertedRoleId = getRoleId(g, allowedRole)
		if convertedRoleId != "" {
			allowedRole = convertedRoleId
		}
		if hasRole(member, allowedRole) {
			return true
		}
	}
	return false
}

func hasRole(member *discordgo.Member, roleId string) bool {
	for _, mRole := range member.Roles {
		if mRole == roleId {
			return true
		}
	}
	return false
}
