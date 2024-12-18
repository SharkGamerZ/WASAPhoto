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
// - [ ] followUser
// - [ ] unfollowUser
// - [ ] getFollowing
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
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Get Routes
	rt.router.GET("/users", rt.wrap(rt.GetUsers))
	rt.router.GET("/users/:id", rt.wrap(rt.GetUserProfile))
	rt.router.GET("/users/:id/following", rt.wrap(rt.getFollowers))

	// Post Routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	// rt.router.POST("/photos", rt.wrap(rt.createPhoto))

	// Put Routes
	rt.router.PUT("/users/:id/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:id/following/:id2", rt.wrap(rt.followUser))

	// Delete Routes
	rt.router.DELETE("/users/:id", rt.wrap(rt.deleteUser))

	return rt.router
}
