package api

import (
	// "encoding/json"
	// "fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// doLogin logs in a user, or creates a new user if the user does not exist.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("doLogin, frocio\n"))

}
