package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
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
		http.Error(w, "You are banned from this user", http.StatusForbidden)
		return
	}

	// Gets the followers of the user
	followers, err := rt.db.GetFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sends the followers
	err = json.NewEncoder(w).Encode(followers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
