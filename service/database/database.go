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

	CreateUser(_struct.User) (int, error)
	DeleteUser(UserID int) error

	UpdateUsername(UserID int, username string) error

	// Follows
	FollowUser(follower int, followed int) error
	UnfollowUser(follower int, followed int) error
	GetFollowings(userID int) ([]_struct.User, error)
	Ping() error
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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// fmt.Println("Creating database structure")
		createUsersTableQuery := `CREATE TABLE USERS (id INTEGER NOT NULL PRIMARY KEY,
													username TEXT
													bio TEXT
													propic TEXT);`
		_, err = db.Exec(createUsersTableQuery)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createFollowersTableQuery := `CREATE TABLE FOLLOWERS (
													followed_id INTEGER references USERS(id),
													follower_id INTEGER references USERS(id),
													PRIMARY KEY (followed_id, follower_id));`
		_, err = db.Exec(createFollowersTableQuery)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
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
	return nil
}
