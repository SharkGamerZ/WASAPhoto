package utils

import "regexp"

// Checks if a username is valid
func isValidUsername(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`).MatchString(username)
}
