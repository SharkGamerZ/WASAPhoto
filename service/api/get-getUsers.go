package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Get all users
func (rt *_router) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Checks if a usernames is passed
	username := r.URL.Query().Get("username")

	// If a username is passed, gets the specifc user
	fmt.Println("User:" + username + ".")
	if username != "" {
		// Get the user by username
		user, err := rt.db.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	// If no username is passed, get all users
	users, err := rt.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
