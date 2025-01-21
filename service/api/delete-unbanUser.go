package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// unbanUser
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the users
	banned, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	banner, err := strconv.Atoi(ps.ByName("id2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checks if the banner is the same as the authenticated user
	if ctx.UserID != banner {
		http.Error(w, "You can't unban this user", http.StatusForbidden)
		return
	}

	// Unbans the user
	err = rt.db.UnbanUser(banned, banner)
	// If the user wasn't banned
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User is not banned", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
