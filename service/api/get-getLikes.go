package api

import (
	// "encoding/json"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getLikes returns a list of users who liked a photo
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the photoID and owner ID from the URL
	photoID := ps.ByName("photoid")
	ownerID := ps.ByName("id")

	// Get the likes of the photo
	likes, err := rt.db.GetLikes(ownerID, photoID)
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
