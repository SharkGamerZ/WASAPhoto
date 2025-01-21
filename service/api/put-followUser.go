package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
	follower, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followed, err := strconv.Atoi(ps.ByName("id2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checks if the uer is authorized
	if follower != ctx.UserID {
		http.Error(w, "You are not authorized", http.StatusUnauthorized)
		return
	}

	// Checks if the user is banned
	banned, err := rt.db.IsBanned(followed, follower)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if banned {
		http.Error(w, "You are banned", http.StatusForbidden)
		return
	}

	// Checks if user exists
	_, err = rt.db.GetUserById(followed)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User does not exist", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Follows the user
	err = rt.db.FollowUser(follower, followed)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User is already following the user", http.StatusConflict)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
