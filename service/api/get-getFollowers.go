package api

import (
	"encoding/json"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the user from the request body
	var user _struct.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Gets the followers of the user
	followers, err := rt.db.GetFollowers(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sends the followers
	err = json.NewEncoder(w).Encode(followers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
