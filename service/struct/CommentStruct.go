package _struct

// Comment is a struct that contains the comment data
type Comment struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	PhotoOwnerID int    `json:"owner_id"`
	PhotoID      int    `json:"post_id"`
	Content      string `json:"content"`
	Timestamp    string `json:"timestamp"`
}
