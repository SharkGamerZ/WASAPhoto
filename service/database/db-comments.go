package database

import (
	"database/sql"
	"errors"

	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
)

// CommentPhoto comments a photo
func (db *appdbimpl) CommentPhoto(ownerID int, photoID int, comment _struct.Comment, userID int) error {
	id, err := db.GetLastCommentID()
	if err != nil {
		return err
	}

	comment.ID = id + 1

	_, err = db.c.Exec("INSERT INTO COMMENTS (id, owner_id, photo_id, user_id, comment, time) VALUES (?, ?, ?, ?, ?, ?)", comment.ID, ownerID, photoID, userID, comment.Content, comment.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetLastCommentID returns the id of the last comment
func (db *appdbimpl) GetLastCommentID() (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM COMMENTS ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// GetCommentByID returns a comment
func (db *appdbimpl) GetCommentByID(commentID int) (_struct.Comment, error) {
	var comment _struct.Comment
	err := db.c.QueryRow("SELECT id, owner_id, photo_id, user_id, comment, time FROM COMMENTS WHERE id = ?", commentID).Scan(&comment.ID, &comment.PhotoOwnerID, &comment.PhotoID, &comment.UserID, &comment.Content, &comment.Timestamp)
	if err != nil {
		return _struct.Comment{}, err
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (db *appdbimpl) DeleteComment(commentID int) error {
	query := "DELETE FROM COMMENTS WHERE id = ?"
	_, err := db.c.Exec(query, commentID)
	return err
}

// GetComments returns comments of a photo
func (db *appdbimpl) GetComments(ownerID, photoID int) ([]_struct.Comment, error) {
	rows, err := db.c.Query("SELECT id, user_id, comment, time FROM COMMENTS WHERE owner_id = ? AND photo_id = ?", ownerID, photoID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []_struct.Comment
	for rows.Next() {
		var comment _struct.Comment
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.Timestamp)
		if err != nil {
			return nil, err
		}
		comment.PhotoOwnerID = ownerID
		comment.PhotoID = photoID
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return comments, nil
}
