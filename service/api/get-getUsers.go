package api

import (
	"encoding/json"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Get all users
func (rt *_router) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Checks if a usernames is passed
	username := r.URL.Query().Get("username")

	// If a username is passed, gets the specific user
	if username != "" {
		// Get the user by username
		user, err := rt.db.GetUsersByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the user
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// If no username is passed, get all users
	users, err := rt.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter users, checking if they banned the ser
	for i, user := range users {
		// Check if the user is banned
		banned, err := rt.db.IsBanned(user.UserID, ctx.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if banned {
			users = append(users[:i], users[i+1:]...)
		}
	}

	// Send the users
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
