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
// - [x] unfollowUser
// - [x] getFollowing
// -----------------
// - [x] banUser
// - [x] unbanUser
// - [x] getBanneds
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
	rt.router.GET("/users/:id/bans", rt.wrap(rt.getBanneds, true))

	// Post Routes
	// rt.router.POST("/photos", rt.wrap(rt.createPhoto, true))

	// Put Routes
	rt.router.PUT("/users/:id/username", rt.wrap(rt.setMyUserName, true))
	rt.router.PUT("/users/:id/following/:id2", rt.wrap(rt.followUser, true))
	rt.router.PUT("/users/:id/bans/:id2", rt.wrap(rt.banUser, true))

	// Delete Routes
	rt.router.DELETE("/users/:id", rt.wrap(rt.deleteUser, true))
	rt.router.DELETE("/users/:id/following/:id2", rt.wrap(rt.unfollowUser, true))
	rt.router.DELETE("/users/:id/bans/:id2", rt.wrap(rt.unbanUser, true))

	return rt.router
}
