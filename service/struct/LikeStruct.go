package _struct

// Like struct is a struct that represents a like
type Like struct {
	OwnerID   int    `json:"ownerID"`
	PhotoID   int    `json:"photoID"`
	UserID    int    `json:"userID"`
	Timestamp string `json:"timestamp"`
}
