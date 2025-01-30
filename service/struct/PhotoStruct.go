package _struct

type Photo struct {
	PhotoID   int    `json:"photoID"`
	UserID    int    `json:"userID"`
	Photo     string `json:"photo"`
	Caption   string `json:"caption"`
	Timestamp string `json:"timestamp"`
	Liked     bool   `json:"liked"`
	Likes     int    `json:"likes"`
}
