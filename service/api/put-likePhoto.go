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
	ownerID := ps.ByName("id")
	photoID := ps.ByName("photoid")
	userID, err := strconv.Atoi(ps.ByName("user_like_id"))

	// Check if the user is authenticated
	if userID != ctx.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Checks if the photo exists
	_, err = rt.db.GetPhotoById(ownerID, photoID)
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
	like.OwnerID, err = strconv.Atoi(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	like.PhotoID, err = strconv.Atoi(photoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	like.UserID = userID
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
