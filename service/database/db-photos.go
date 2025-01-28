package database

import (
	"database/sql"
	"errors"

	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// CreatePhoto creates a new photo.
func (db *appdbimpl) CreatePhoto(photo _struct.Photo) error {
	id, err := db.GetLastPhotoIDOfUser(photo.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	photo.PhotoID = id + 1

	_, err = db.c.Exec("INSERT INTO photos (id, user_id, photo, caption, timestamp) VALUES (?, ?, ?, ?, ?)", photo.PhotoID, photo.UserID, photo.Photo, photo.Caption, photo.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetLastPhotoID returns the id of the last photo.
func (db *appdbimpl) GetLastPhotoIDOfUser(userid int) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM photos WHERE user_id = ? ORDER BY id DESC LIMIT 1", userid).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// GetUserPhotos returns the photos of a user.
func (db *appdbimpl) GetUserPhotos(id int, viewerID int) ([]_struct.Photo, error) {
	const query = `
		SELECT p.id, p.user_id, p.photo, p.caption, p.timestamp,
			   CASE WHEN l.user_id IS NOT NULL THEN 1 ELSE 0 END as liked
		FROM photos p
		LEFT JOIN likes l ON p.id = l.photo_id AND l.user_id = ?
		WHERE p.user_id = ?
		ORDER BY p.timestamp DESC
	`

	rows, err := db.c.Query(query, viewerID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []_struct.Photo
	for rows.Next() {
		var photo _struct.Photo
		var liked int
		err = rows.Scan(
			&photo.PhotoID,
			&photo.UserID,
			&photo.Photo,
			&photo.Caption,
			&photo.Timestamp,
			&liked,
		)
		if err != nil {
			return nil, err
		}
		photo.Liked = liked == 1
		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return photos, nil
}

// GetPhotoById returns the photo with the given id.
func (db *appdbimpl) GetPhotoById(userid int, photoid int, viewerID int) (_struct.Photo, error) {
	const query = `
		SELECT p.id, p.user_id, p.photo, p.caption, p.timestamp,
			   CASE WHEN l.user_id IS NOT NULL THEN 1 ELSE 0 END as liked
		FROM photos p
		LEFT JOIN likes l ON p.id = l.photo_id AND l.user_id = ?
		WHERE p.user_id = ? AND p.id = ?
	`

	var photo _struct.Photo
	var liked int
	err := db.c.QueryRow(query, viewerID, userid, photoid).Scan(
		&photo.PhotoID,
		&photo.UserID,
		&photo.Photo,
		&photo.Caption,
		&photo.Timestamp,
		&liked,
	)
	if err != nil {
		return photo, err
	}

	photo.Liked = liked == 1
	return photo, nil
}

// DeletePhoto deletes the photo with the given id.
func (db *appdbimpl) DeletePhoto(userid int, photoid int) error {
	_, err := db.c.Exec("DELETE FROM photos WHERE user_id = ? AND id = ?", userid, photoid)
	if err != nil {
		return err
	}

	return nil
}
