package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUserPhotos returns the photos of a user.
func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	url_id := ps.ByName("id")

	// Get the photos of the user
	photos, err := rt.db.GetUserPhotos(url_id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User has no posts", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the photos
	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
