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
		rt.baseLogger.Errorf("Invalid username: %v", user.Username)
		return
	}

	// Checks if user exists
	userExists, err := rt.db.ExistsName(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rt.baseLogger.Errorf("Error checking if user exists: %v", err)
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
			rt.baseLogger.Errorf("Error creating user: %v", err)
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
			rt.baseLogger.Errorf("Error encoding response: %v", err)
			return
		}

	} else {
		// If user exists, return OK
		w.WriteHeader(http.StatusOK)

		// Return the ID of the user
		var id int
		users, err := rt.db.GetUsersByUsername(user.Username, 0)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			rt.baseLogger.Errorf("Error getting user by username: %v", err)
			return
		}

		id = users[0].UserID

		err = json.NewEncoder(w).Encode(struct {
			ID int `json:"id"`
		}{
			ID: id,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			rt.baseLogger.Errorf("Error encoding response: %v", err)
			return
		}
	}

}
