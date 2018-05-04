package hdsbrute

import (
	"os"
	"strings"
)

func getEnvRoles(envKey string) []string {
	var result []string
	envRoles, ok := os.LookupEnv(envKey)
	if ok {
		for _, role := range strings.Split(envRoles, ",") {
			result = append(result, role)
		}
	}
	return result
}

func GetMemberRoles() []string {
	return getEnvRoles("MEMBER_ROLES")
}

func GetAdminRoles() []string {
	return getEnvRoles("ADMIN_ROLES")
}
