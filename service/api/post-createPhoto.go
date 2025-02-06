package api

import (
	// "encoding/json"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	_struct "github.com/SharkGamerZ/WASAPhoto/service/struct"

	// "github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// createPhoto creates a new photo.
func (rt *_router) createPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")

	// Gets the id of the user from the URL
	url_id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Checks if the user is authorized
	if url_id != ctx.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode the request body as a multipart form
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Close the file after the function ends
	defer file.Close()

	// Read the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the caption from the form
	caption := r.FormValue("caption")

	// Create the photo
	var photo _struct.Photo
	photo.UserID = url_id
	photo.Photo = base64.StdEncoding.EncodeToString(fileBytes)
	photo.Caption = caption
	photo.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	// Saves the photo in the database
	photoid, err := rt.db.CreatePhoto(photo)

	// Checks if the photo was saved
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Returns the photo id
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(strconv.Itoa(photoid)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
