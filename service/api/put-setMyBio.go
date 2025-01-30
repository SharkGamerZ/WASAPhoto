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
func (rt *_router) setMyBio(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checks if the user is authenticated
	if userID != ctx.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Gets the user from the request body
	var bio string
	err = json.NewDecoder(r.Body).Decode(&bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Check if the bio is valid

	// Updates the user's bio
	err = rt.db.UpdateBio(userID, bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
