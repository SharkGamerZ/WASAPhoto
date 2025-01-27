package database

import (
	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

func (db *appdbimpl) FollowUser(follower int, followed int) error {
	// Checks if the user is already following the user
	_, err := db.c.Exec("SELECT * FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?", follower, followed)
	if err != nil {
		return err
	}

	// If not, follow the user
	query := "INSERT INTO FOLLOWERS (follower_id, followed_id) VALUES (?, ?)"
	_, err = db.c.Exec(query, follower, followed)
	return err
}

func (db *appdbimpl) UnfollowUser(follower int, followed int) error {
	// Checks if the user is following the user
	_, err := db.c.Exec("SELECT * FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?", follower, followed)
	if err != nil {
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

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

// CountFollowers returns the number of followers of a user
func (db *appdbimpl) CountFollowers(userID int) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM FOLLOWERS WHERE followed_id = ?", userID).Scan(&count)
	return count, err
}

// CountFollowings returns the number of followings of a user
func (db *appdbimpl) CountFollowings(userID int) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM FOLLOWERS WHERE follower_id = ?", userID).Scan(&count)
	return count, err
}
