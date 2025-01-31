package database

import (
	"database/sql"
	"errors"

	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// LikePhoto likes a photo
func (db *appdbimpl) LikePhoto(like _struct.Like) error {
	// Checks if the user has already liked the photo
	var tmp int

	err := db.c.QueryRow("SELECT owner_id, photo_id, user_id  FROM LIKES WHERE owner_id = ? AND photo_id = ? AND user_id = ?", like.OwnerID, like.PhotoID, like.UserID).Scan(&tmp, &tmp, &tmp)
	if errors.Is(err, sql.ErrNoRows) {
		// If not, like the photo
		query := "INSERT INTO LIKES (owner_id, photo_id, user_id, time) VALUES (?, ?, ?, ?)"
		_, err2 := db.c.Exec(query, like.OwnerID, like.PhotoID, like.UserID, like.Timestamp)
		if err2 != nil {
			return err2
		}
	}

	return err
}

// UnlikePhoto unlikes a photo
func (db *appdbimpl) UnlikePhoto(like _struct.Like) error {
	// Checks if the user has already liked the photo
	var tmp int

	err := db.c.QueryRow("SELECT owner_id, photo_id, user_id  FROM LIKES WHERE owner_id = ? AND photo_id = ? AND user_id = ?", like.OwnerID, like.PhotoID, like.UserID).Scan(&tmp, &tmp, &tmp)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("User has not liked the photo")
	}

	// If yes, unlike the photo
	query := "DELETE FROM LIKES WHERE owner_id = ? AND photo_id = ? AND user_id = ?"
	_, err = db.c.Exec(query, like.OwnerID, like.PhotoID, like.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetLikes returns a list of users who liked a photo
func (db *appdbimpl) GetLikes(ownerID, photoID int) ([]_struct.User, error) {
	rows, err := db.c.Query(`SELECT user_id, propic, username
							FROM LIKES, USERS
							WHERE owner_id = ? AND photo_id = ?
							AND LIKES.user_id = USERS.id`, ownerID, photoID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []_struct.User
	for rows.Next() {
		var user _struct.User

		err = rows.Scan(&user.UserID, &user.Propic, &user.Username)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}
