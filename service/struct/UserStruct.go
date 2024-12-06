package _struct

import (
	"regexp"
)

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Propic   string `json:"propic"`
}

func (u *User) DoesHaveValidUsername() bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`).MatchString(u.Username)
}
