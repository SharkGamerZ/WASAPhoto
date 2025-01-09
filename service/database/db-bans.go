package database

import (
	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// BanUser bans a user from another user
func (db *appdbimpl) BanUser(banner int, banned int) error {
	// Checks if the user is already banned
	_, err := db.c.Exec("SELECT * FROM BANS WHERE banner_id = ? AND banned_id = ?", banner, banned)
	if err != nil {
		return err
	}

	// If not, ban the user
	query := "INSERT INTO BANS (banner_id, banned_id) VALUES (?, ?)"
	_, err = db.c.Exec(query, banner, banned)
	return err
}

// UnbanUser unbans a user from another user
func (db *appdbimpl) UnbanUser(banner int, banned int) error {
	// Checks if the user is already banned
	_, err := db.c.Exec("SELECT * FROM BANS WHERE banner_id = ? AND banned_id = ?", banner, banned)
	if err != nil {
		return err
	}

	query := "DELETE FROM BANS WHERE banner_id = ? AND banned_id = ?"
	_, err = db.c.Exec(query, banner, banned)
	return err
}

// GetBanneds returns all the users that are banned from the user
func (db *appdbimpl) GetBanneds(userID int) ([]_struct.User, error) {
	rows, err := db.c.Query("SELECT banned_id, username FROM BANS, USERS WHERE banner_id = ? AND banned_id = id", userID)
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
