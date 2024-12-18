package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
	follower, err := strconv.Atoi(ps.ByName("id"))
	followed, err := strconv.Atoi(ps.ByName("id2"))

	// Checks if user exists
	_, err = rt.db.GetUserById(followed)
	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Follows the user
	err = rt.db.FollowUser(follower, followed)
	if err == sql.ErrNoRows {
		http.Error(w, "User is already following the user", http.StatusConflict)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
