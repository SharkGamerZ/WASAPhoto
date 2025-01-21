package database

import (
	// "database/sql"

	"database/sql"
	"errors"

	"github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// Checks if user exists in the database
func (db *appdbimpl) ExistsName(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// Gets user from DB
func (db *appdbimpl) GetUserById(userID int) (_struct.User, error) {
	var user _struct.User
	err := db.c.QueryRow(`SELECT id, username FROM USERS WHERE id= ?;`, userID).Scan(&user.UserID, &user.Username)
	return user, err
}

// Gets user from DB
func (db *appdbimpl) GetUserByUsername(username string) (_struct.User, error) {
	var user _struct.User
	err := db.c.QueryRow(`SELECT id, username FROM USERS WHERE username= ?;`, username).Scan(&user.UserID, &user.Username)
	return user, err
}

// Gets Users from the database
func (db *appdbimpl) GetUsers() ([]_struct.User, error) {
	rows, err := db.c.Query("SELECT id, username FROM USERS")
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

// Create a new user in the database
func (db *appdbimpl) CreateUser(user _struct.User) (int, error) {
	_, err := db.c.Exec("INSERT INTO USERS (username) VALUES (?)", user.Username)
	if err != nil {
		return 0, err
	}

	// Gets the ID of the user
	var id int
	err = db.c.QueryRow("SELECT id FROM USERS WHERE username = ?", user.Username).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

// Delete a user from the database
func (db *appdbimpl) DeleteUser(userID int) error {
	_, err := db.c.Exec("DELETE FROM USERS WHERE id = ?", userID)
	return err
}

// Changes the user's username
func (db *appdbimpl) UpdateUsername(userID int, newUsername string) error {
	_, err := db.c.Exec("UPDATE USERS SET username = ? WHERE id = ?", newUsername, userID)
	return err
}

// GetMyStream gets the user's stream
func (db *appdbimpl) GetMyStream(userID int) ([]_struct.Photo, error) {
	// Gets the list of users followed by the user
	rows, err := db.c.Query("SELECT followed_id FROM FOLLOWERS WHERE follower_id = ?", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var photos []_struct.Photo

	// Gets the photos of the users followed by the user
	for rows.Next() {
		var currentUserID int
		err = rows.Scan(&currentUserID)
		if err != nil {
			return nil, err
		}

		// Checks if the user is banned
		isBanned, err := db.IsBanned(currentUserID, userID)
		if err != nil {
			return nil, err
		}
		if isBanned {
			continue
		}

		userPhotos, err := db.GetUserPhotos(currentUserID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, userPhotos...)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return photos, nil
}
