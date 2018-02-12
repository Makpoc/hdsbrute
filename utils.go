package hdsbrute

import "os"

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
