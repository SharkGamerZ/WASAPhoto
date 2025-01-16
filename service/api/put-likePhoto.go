package api

import (
	// "database/sql"
	// "encoding/json"
	// "errors"
	"net/http"

	"github.com/SharkGamerZ/WASAPhoto/service/api/reqcontext"
	// "github.com/SharkGamerZ/WASAPhoto/service/struct"
	"github.com/julienschmidt/httprouter"
)

// likePhoto likes the photo with the given id.
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// w.Header().Set("Content-type", "application/json")
	//
	// // Gets the id of the user from the URL
	// url_id := ps.ByName("id")
	// photo_id := ps.ByName("photoid")
	//
	// // Gets the user who liked the photo
	// user, err := rt.db.GetUserById(ctx.UserID)
	// if errors.Is(err, sql.ErrNoRows) {
	// 	http.Error(w, "User not found", http.StatusNotFound)
	// 	return
	// }
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Gets the photo
	// photo, err := rt.db.GetPhotoById(url_id, photo_id)
	// if errors.Is(err, sql.ErrNoRows) {
	// 	http.Error(w, "Photo not found", http.StatusNotFound)
	// 	return
	// }
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Checks if the user has already liked the photo
	// for _, like := range photo.Likes {
	// 	if like.UserID == user.ID {
	// 		http.Error(w, "You have already liked this photo", http.StatusForbidden)
	// 		return
	// 	}
	// }
	//
	// // Likes the photo
	// photo.Likes = append(photo.Likes, struct.Like{UserID: user.ID})
	// err = rt.db.UpdatePhoto(url_id, photo)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	//
}
