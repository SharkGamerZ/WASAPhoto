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

// getLikes returns a list of users who liked a photo
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the photoID and owner ID from the URL
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	photoOwnerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if the user is banned from the owner
	banned, err := rt.db.IsBanned(ctx.UserID, photoOwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "You are banned from this user", http.StatusForbidden)
		return
	}

	// Get the likes of the photo
	likes, err := rt.db.GetLikes(photoOwnerID, photoID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "No likes found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the likes
	err = json.NewEncoder(w).Encode(likes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
