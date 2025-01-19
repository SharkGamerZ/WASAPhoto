package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// UnlikePhoto unlikes a photo
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	url_id := ps.ByName("id")
	photo_id := ps.ByName("photoid")
	user_like_id := ps.ByName("user_like_id")

	// Checks if the user is authenticated
	if user_like_id != strconv.Itoa(ctx.UserID) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Checks if the user has already liked the photo
	var like _struct.Like
	like.OwnerID, _ = strconv.Atoi(url_id)
	like.PhotoID, _ = strconv.Atoi(photo_id)
	like.UserID, _ = strconv.Atoi(user_like_id)

	err := rt.db.UnlikePhoto(like)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the success message
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Photo unliked"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
