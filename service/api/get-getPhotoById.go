package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getPhotoById returns the photo with the given id.
func (rt *_router) getPhotoById(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	url_id := ps.ByName("id")
	photo_id := ps.ByName("photoid")

	// Get the photo
	photo, err := rt.db.GetPhotoById(url_id, photo_id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the photo
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
