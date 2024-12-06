package database

import (
	// "database/sql"

	"github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// Checks if user exists in the database
func (db *appdbimpl) CheckUser(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// Create a new user in the database
func (db *appdbimpl) CreateUser(user _struct.User) error {
	_, err := db.c.Exec("INSERT INTO USERS (username) VALUES (?)", user.Username)
	return err
}
