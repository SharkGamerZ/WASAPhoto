package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user
	banner, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	banned, err := strconv.Atoi(ps.ByName("id2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checks if the uer is authorized
	if banner != ctx.UserID {
		http.Error(w, "You are not authorized", http.StatusUnauthorized)
		return
	}

	// Checks if user exists
	_, err = rt.db.GetUserById(banned, ctx.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User does not exist", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Bans the user
	err = rt.db.BanUser(banner, banned)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusCreated)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
