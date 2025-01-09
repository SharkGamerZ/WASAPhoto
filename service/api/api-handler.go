package api

import (
	"net/http"
)

// TODO:
// USERS
// - [x] doLogin (see simplified login)
// - [x] getUsers
// - [x] getUserProfile
// - [x] deleteUser
// - [x] setMyUserName
// - [ ] getMyStream
// -----------------
// - [ ] uploadPhoto
// - [ ] getPhotoById
// - [ ] getUserPhotos
// - [ ] deletePhoto
// -----------------
// - [x] followUser
// - [ ] unfollowUser
// - [x] getFollowing
// -----------------
// - [ ] banUser
// -----------------
// - [ ] unbanUser
// - [ ] getBanned
// -----------------
// - [ ] likePhoto
// - [ ] unlikePhoto
// - [ ] getLikes
// -----------------
// - [ ] commentPhoto
// - [ ] uncommentPhoto
// - [ ] getComments
// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Session
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	// Get Routes
	rt.router.GET("/users", rt.wrap(rt.GetUsers, true))
	rt.router.GET("/users/:id", rt.wrap(rt.GetUserProfile, true))
	rt.router.GET("/users/:id/following", rt.wrap(rt.getFollowings, true))

	// Post Routes
	// rt.router.POST("/photos", rt.wrap(rt.createPhoto))

	// Put Routes
	rt.router.PUT("/users/:id/username", rt.wrap(rt.setMyUserName, true))
	rt.router.PUT("/users/:id/following/:id2", rt.wrap(rt.followUser, true))

	// Delete Routes
	rt.router.DELETE("/users/:id", rt.wrap(rt.deleteUser, true))
	rt.router.DELETE("/users/:id/following/:id2", rt.wrap(rt.unfollowUser, true))

	return rt.router
}
