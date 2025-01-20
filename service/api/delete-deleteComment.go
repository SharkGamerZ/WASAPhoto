package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// DeleteComment deletes a comment with the given id.

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	// userID, err := strconv.Atoi(ps.ByName("id"))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// photoID, err := strconv.Atoi(ps.ByName("photoid"))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	commentID, err := strconv.Atoi(ps.ByName("comment_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if the user is the owner of the comment
	comment, err := rt.db.GetCommentByID(commentID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if comment.UserID != ctx.UserID {
		http.Error(w, "You are not the owner of this comment", http.StatusForbidden)
		return
	}

	// Delete the comment
	err = rt.db.DeleteComment(commentID)

	// Check if the comment was not found, double check
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the success message
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
