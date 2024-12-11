package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Get all users
func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Gets the user id from the request
	userID := ps.ByName("id")

	// Convert the user id to an int
	userIDInt, err := strconv.Atoi(userID)
	user, err := rt.db.GetUserById(userIDInt)

	// If user doesn't exist
	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
