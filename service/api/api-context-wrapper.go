package api

import (
	// "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"net/http"
	"strconv"
	"strings"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler, auth bool) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Configura gli header CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Specifica l'origine in produzione
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Gestisci le richieste preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Authorization check
		userID := 0
		if auth {
			// Check if the user is authorized
			userID = getAuthToken(r.Header)
			if userID == 0 {
				http.Error(w, "User not authorized", http.StatusUnauthorized)
				return
			}
		}

		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
			UserID:  userID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

func getAuthToken(header http.Header) int {
	// The string is formatted as "Bearer <token>"
	token := header.Get("Authorization")

	// Rimuove eventuali prefissi "Bearer " se presenti
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// Verifica che il token sia un numero
	userID, err := strconv.Atoi(token)
	if err != nil {
		return 0
	}
	return userID
}
