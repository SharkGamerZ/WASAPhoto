package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the users
	follower, err := strconv.Atoi(ps.ByName("id"))
	followed, err := strconv.Atoi(ps.ByName("id2"))

	// Unfollows the user
	err = rt.db.UnfollowUser(follower, followed)
	if err == sql.ErrNoRows {

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sends the response
	w.Write([]byte(`{"success":true}`))
}
