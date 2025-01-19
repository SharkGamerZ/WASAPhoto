package _struct

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Propic   string `json:"propic"`
	Bio      string `json:"bio"`
}
