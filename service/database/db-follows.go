package database

import (
	"database/sql"
	"errors"

	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

func (db *appdbimpl) FollowUser(follower int, followed int) error {
	// Checks if the user is already following the user
	_, err := db.c.Exec("SELECT * FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?", follower, followed)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// If not, follow the user
	query := "INSERT INTO FOLLOWERS (follower_id, followed_id) VALUES (?, ?)"
	_, err = db.c.Exec(query, follower, followed)
	return err
}

func (db *appdbimpl) UnfollowUser(follower int, followed int) error {
	// Checks if the user is already following the user
	_, err := db.c.Exec("SELECT * FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?", follower, followed)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.c.Exec("DELETE FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?", follower, followed)

	return err
}

func (db *appdbimpl) GetFollowings(userID int) ([]_struct.User, error) {
	rows, err := db.c.Query("SELECT followed_id, username FROM FOLLOWERS, USERS WHERE follower_id = ? AND followed_id = id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []_struct.User
	for rows.Next() {
		var user _struct.User
		err = rows.Scan(&user.UserID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
