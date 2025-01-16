package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// DeletePhoto deletes the photo with the given id.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	url_id := ps.ByName("id")
	photo_id := ps.ByName("photoid")

	// Checks if the user is the owner of the photo
	photo, err := rt.db.GetPhotoById(url_id, photo_id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if photo.UserID != ctx.UserID {
		http.Error(w, "You are not the owner of this photo", http.StatusForbidden)
		return
	}

	// Delete the photo
	err = rt.db.DeletePhoto(url_id, photo_id)

	// Check if the photo was not found, double check
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Photo deleted"})
}
