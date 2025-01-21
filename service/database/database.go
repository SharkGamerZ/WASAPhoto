/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User
	ExistsName(username string) (bool, error)
	GetUserById(userID int) (_struct.User, error)
	GetUserByUsername(username string) (_struct.User, error)
	GetUsers() ([]_struct.User, error)
	GetMyStream(userID int) ([]_struct.Photo, error)

	CreateUser(_struct.User) (int, error)
	DeleteUser(UserID int) error

	UpdateUsername(UserID int, username string) error

	// Follows
	FollowUser(follower int, followed int) error
	UnfollowUser(follower int, followed int) error
	GetFollowings(userID int) ([]_struct.User, error)

	// Bans
	BanUser(banner int, banned int) error
	UnbanUser(banner int, banned int) error
	GetBanneds(userID int) ([]_struct.User, error)
	IsBanned(banner int, banned int) (bool, error)

	// Photos
	CreatePhoto(photo _struct.Photo) error
	GetUserPhotos(id int) ([]_struct.Photo, error)
	GetPhotoById(userid int, photoid int) (_struct.Photo, error)
	DeletePhoto(userid int, photoid int) error

	// Likes
	LikePhoto(like _struct.Like) error
	UnlikePhoto(like _struct.Like) error
	GetLikes(ownerID, photoID int) ([]_struct.User, error)

	// Comments
	CommentPhoto(ownerID int, photoID int, comment _struct.Comment, userID int) error
	GetCommentByID(commentID int) (_struct.Comment, error)
	DeleteComment(commentID int) error
	GetComments(ownerID, photoID int) ([]_struct.Comment, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// fmt.Println("Deleting database")
	// DeleteDatabase(db)

	// Drops table comments
	// _, _ = db.Exec("DROP TABLE COMMENTS")

	// Creates the database structure if it doesn't exist
	// USERS
	createUsersTableQuery := `CREATE TABLE IF NOT EXISTS USERS (id INTEGER NOT NULL PRIMARY KEY,
												username TEXT
												bio TEXT
												propic TEXT);`
	_, err := db.Exec(createUsersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating USERS table: %w", err)
	}

	// FOLLOWERS
	createFollowersTableQuery := `CREATE TABLE IF NOT EXISTS FOLLOWERS (
												followed_id INTEGER references USERS(id),
												follower_id INTEGER references USERS(id),
												PRIMARY KEY (followed_id, follower_id));`
	_, err = db.Exec(createFollowersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating FOLLOWERS table: %w", err)
	}

	// BANS
	createBansTableQuery := `CREATE TABLE IF NOT EXISTS BANS (
												banned_id INTEGER references USERS(id),
												banner_id INTEGER references USERS(id),
												PRIMARY KEY (banned_id, banner_id));`
	_, err = db.Exec(createBansTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating BANS table: %w", err)
	}

	// PHOTOS
	createPhotosTableQuery := `CREATE TABLE IF NOT EXISTS PHOTOS (
												id INTEGER NOT NULL,
												user_id INTEGER NOT NULL references USERS(id),
												photo TEXT,
												caption TEXT,
												timestamp DATETIME,
												PRIMARY KEY (id, user_id));`
	_, err = db.Exec(createPhotosTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating PHOTOS table: %w", err)
	}

	// LIKES
	createLikesTableQuery := `CREATE TABLE IF NOT EXISTS LIKES (
												owner_id INTEGER references USERS(id),
												photo_id INTEGER references PHOTOS(id),
												user_id INTEGER references USERS(id),
												time DATETIME,
												PRIMARY KEY (owner_id, photo_id, user_id));`
	_, err = db.Exec(createLikesTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating LIKES table: %w", err)
	}

	// COMMENTS
	createCommentsTableQuery := `CREATE TABLE IF NOT EXISTS COMMENTS (
												id INTEGER NOT NULL PRIMARY KEY,
												owner_id INTEGER references USERS(id),
												photo_id INTEGER references PHOTOS(id),
												user_id INTEGER references USERS(id),
												comment TEXT,
												time DATETIME);`

	_, err = db.Exec(createCommentsTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating COMMENTS table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// Deletes the database
func DeleteDatabase(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE USERS")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE FOLLOWERS")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE BANS")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE PHOTOS")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE LIKES")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE COMMENTS")
	if err != nil {
		return err
	}

	return nil
}
