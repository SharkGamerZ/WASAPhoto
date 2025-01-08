package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	// "fmt"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin logs in a user, or creates a new user if the user does not exist.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gets the user from the request body
	var username string
	err = json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if the username is valid
	if !isValidUsername(username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Checks if user exists
	userExists, err := rt.db.ExistsName(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userExists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Updates the user's username
	err = rt.db.UpdateUsername(userID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
