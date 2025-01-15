package api

import (
	// "encoding/json"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	// "github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// createPhoto creates a new photo.
func (rt *_router) createPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// w.Header().Set("Content-type", "application/json")
	//
	// // Gets the photo from the request body
	// var photo _struct.Photo
	// err := json.NewDecoder(r.Body).Decode(&photo)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	//
	// // Checks if the photo is valid
	// if !isValidPhoto(photo) {
	// 	http.Error(w, "Invalid photo", http.StatusBadRequest)
	// 	return
	// }
	//
	// // Create photo
	// var id int
	// id, err = rt.db.CreatePhoto(photo)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Return the ID of the photo
	// err = json.NewEncoder(w).Encode(struct {
	// 	ID int `json:"id"`
	// }{
	// 	ID: id,
	// })
	//
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}
