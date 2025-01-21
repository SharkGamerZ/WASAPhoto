package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// likePhoto likes the photo with the given id.
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	photoOwnerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	likeOwnerID, err := strconv.Atoi(ps.ByName("user_like_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is authenticated
	if likeOwnerID != ctx.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Checks if the user is banned
	banned, err := rt.db.IsBanned(photoOwnerID, likeOwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "You are banned", http.StatusForbidden)
		return
	}

	// Checks if the photo exists
	_, err = rt.db.GetPhotoById(photoOwnerID, photoID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creates the like
	var like _struct.Like
	like.OwnerID = photoOwnerID
	like.PhotoID = photoID
	like.UserID = likeOwnerID
	like.Timestamp = time.Now().Format("2025-01-01 00:00:00")

	// Likes the photo
	err = rt.db.LikePhoto(like)
	if errors.Is(err, sql.ErrNoRows) {
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
