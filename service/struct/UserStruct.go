package _struct

type User struct {
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	Propic     string `json:"propic"`
	Bio        string `json:"bio"`
	Followed   bool   `json:"followed"`
	Followers  int    `json:"followers"`
	Followings int    `json:"followings"`
	PostNum    int    `json:"post_num"`
}
