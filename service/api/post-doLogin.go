package api

import (
	"encoding/json"
	"net/http"

	// "fmt"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// doLogin logs in a user, or creates a new user if the user does not exist.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the user from the request body
	var user _struct.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if the username is valid
	if !isValidUsername(user.Username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Checks if user exists
	userExists, err := rt.db.ExistsName(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If user does not exist, create a new user
	if !userExists {
		w.WriteHeader(http.StatusCreated)

		// Create user
		var id int
		id, err = rt.db.CreateUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the ID of the user
		err = json.NewEncoder(w).Encode(struct {
			ID int `json:"id"`
		}{
			ID: id,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// If user exists, return OK
		w.WriteHeader(http.StatusOK)
	}

}
