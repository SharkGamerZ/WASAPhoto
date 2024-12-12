package database

import _struct "github.com/SharkGamerZ/WASAPhoto/service/struct"

func (db *appdbimpl) FollowUser(userID int, username string) error {
	_, err := db.c.Exec("INSERT INTO FOLLOWS (follower_id, followed_id) VALUES (?, ?)", userID, username)
	return err
}

func (db *appdbimpl) UnfollowUser(userID int, username string) error {
	_, err := db.c.Exec("DELETE FROM FOLLOWS WHERE follower_id = ? AND followed_id = ?", userID, username)
	return err
}

func (db *appdbimpl) GetFollowers(userID int) ([]_struct.User, error) {
	rows, err := db.c.Query("SELECT followed_id FROM FOLLOWS WHERE follower_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []_struct.User
	for rows.Next() {
		var user _struct.User
		err = rows.Scan(&user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
