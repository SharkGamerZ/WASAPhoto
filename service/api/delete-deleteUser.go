package api

import (
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	// "fmt"
)

// Deltes a user
func (rt *_router) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the user from the request body
	userID, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.Logger.WithField("id", userID).Debug("deleting user")

	// Checks if user exists
	// userExists, err := rt.db.ExistsName(user.Username)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// If user does not exist, return 404
	// if !userExists {
	// 	http.Error(w, "User not found", http.StatusNotFound)
	// 	return
	// }

	// If user exists, delete user
	err = rt.db.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
