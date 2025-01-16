package database

import (
	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// CreatePhoto creates a new photo.
func (db *appdbimpl) CreatePhoto(photo _struct.Photo) error {
	_, err := db.c.Exec("INSERT INTO photos (user_id, photo, caption, timestamp) VALUES (?, ?, ?, ?)", photo.UserID, photo.Photo, photo.Caption, photo.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetUserPhotos returns the photos of a user.
func (db *appdbimpl) GetUserPhotos(id string) ([]_struct.Photo, error) {
	rows, err := db.c.Query("SELECT * FROM photos WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var photos []_struct.Photo
	for rows.Next() {
		var photo _struct.Photo

		err = rows.Scan(&photo.PhotoID, &photo.UserID, &photo.Photo, &photo.Caption, &photo.Timestamp)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return photos, nil
}

// GetPhotoById returns the photo with the given id.
func (db *appdbimpl) GetPhotoById(userid string, photoid string) (_struct.Photo, error) {
	var photo _struct.Photo

	err := db.c.QueryRow("SELECT * FROM photos WHERE user_id = ? AND id = ?", userid, photoid).Scan(&photo.PhotoID, &photo.UserID, &photo.Photo, &photo.Caption, &photo.Timestamp)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

// DeletePhoto deletes the photo with the given id.
func (db *appdbimpl) DeletePhoto(userid string, photoid string) error {
	_, err := db.c.Exec("DELETE FROM photos WHERE user_id = ? AND id = ?", userid, photoid)
	if err != nil {
		return err
	}

	return nil
}
