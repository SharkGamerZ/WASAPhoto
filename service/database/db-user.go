package database

import (
	// "database/sql"

	"github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// Checks if user exists in the database
func (db *appdbimpl) ExistsName(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// Create a new user in the database
func (db *appdbimpl) CreateUser(user _struct.User) (int, error) {
	_, err := db.c.Exec("INSERT INTO USERS (username) VALUES (?)", user.Username)

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
	_, err := db.c.Exec("DELETE FROM USERS WHERE userID = ?", userID)
	return err
}
