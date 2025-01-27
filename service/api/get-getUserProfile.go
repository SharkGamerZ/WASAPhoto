package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Get all users
func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Gets the user id from the request
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checks if the user is banned
	banned, err := rt.db.IsBanned(userID, ctx.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "User is banned", http.StatusForbidden)
		return
	}

	// Gets the user profile
	user, err := rt.db.GetUserById(userID, ctx.UserID)

	// If user doesn't exist
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}
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
}
