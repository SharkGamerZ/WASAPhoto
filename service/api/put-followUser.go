package api

import (
	"encoding/json"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the user from the request body
	var user _struct.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if user exists
	userExists, err := rt.db.ExistsName(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !userExists {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Follows the user
	err = rt.db.FollowUser(user.UserID, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
