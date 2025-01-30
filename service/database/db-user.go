package database

import (
	// "database/sql"

	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"

	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// Checks if user exists in the database
func (db *appdbimpl) ExistsName(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// Gets user from DB
func (db *appdbimpl) GetUserById(userID int, requesterID int) (_struct.User, error) {
	var user _struct.User
	err := db.c.QueryRow(`SELECT id, username, bio, propic FROM USERS WHERE id= ?;`, userID).Scan(&user.UserID, &user.Username, &user.Bio, &user.Propic)
	if err != nil {
		return user, err
	}

	// Check if the requester follows the user
	err = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?)`, requesterID, userID).Scan(&user.Followed)
	if err != nil {
		return user, err
	}

	// Gets the number of followers
	err = db.c.QueryRow(`SELECT COUNT(*) FROM FOLLOWERS WHERE followed_id = ?`, userID).Scan(&user.Followers)
	if err != nil {
		return user, err
	}

	// Gets the number of followings
	err = db.c.QueryRow(`SELECT COUNT(*) FROM FOLLOWERS WHERE follower_id = ?`, userID).Scan(&user.Followings)
	if err != nil {
		return user, err
	}

	// Gets the number of posts
	err = db.c.QueryRow(`SELECT COUNT(*) FROM PHOTOS WHERE user_id = ?`, userID).Scan(&user.PostNum)
	if err != nil {
		return user, err
	}
	return user, err
}

// Gets users from DB by username
func (db *appdbimpl) GetUsersByUsername(username string, requesterID int) ([]_struct.User, error) {
	rows, err := db.c.Query(`SELECT id, username, propic FROM USERS WHERE username LIKE ?;`, "%"+username+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []_struct.User
	for rows.Next() {
		var user _struct.User
		err = rows.Scan(&user.UserID, &user.Username, &user.Propic)
		if err != nil {
			return nil, err
		}

		// Check if the requester follows the user
		err = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?)`, requesterID, user.UserID).Scan(&user.Followed)
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

// Gets Users from the database
func (db *appdbimpl) GetUsers(requesterID int) ([]_struct.User, error) {
	rows, err := db.c.Query("SELECT id, username, propic FROM USERS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []_struct.User
	for rows.Next() {
		var user _struct.User
		err = rows.Scan(&user.UserID, &user.Username, &user.Propic)
		if err != nil {
			return nil, err
		}

		// Check if the requester follows the user
		err = db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM FOLLOWERS WHERE follower_id = ? AND followed_id = ?)`, requesterID, user.UserID).Scan(&user.Followed)
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
	// Open the default profile picture file
	file, err := os.Open("assets/default.png")
	if err != nil {
		return 0, err
	}
	defer file.Close() // Ensure the file is closed when done

	// Read the file contents
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	// Encode the image to base64
	defaultPropic := base64.StdEncoding.EncodeToString(fileBytes)

	_, err = db.c.Exec("INSERT INTO USERS (username, propic) VALUES (?, ?)", user.Username, defaultPropic)
	if err != nil {
		return 0, err
	}

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

// Changes the user's bio
func (db *appdbimpl) UpdateBio(userID int, newBio string) error {
	_, err := db.c.Exec("UPDATE USERS SET bio = ? WHERE id = ?", newBio, userID)
	return err
}

// Changes the user's profile picture
func (db *appdbimpl) UpdatePropic(userID int, newPropic string) error {
	_, err := db.c.Exec("UPDATE USERS SET propic = ? WHERE id = ?", newPropic, userID)
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

		// Pass the userID as viewerID to check for likes
		userPhotos, err := db.GetUserPhotos(currentUserID, userID)
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
