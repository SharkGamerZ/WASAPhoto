package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Post Routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	// rt.router.POST("/photos", rt.wrap(rt.createPhoto))

	// Delete Routes
	rt.router.DELETE("/users/{user_id}", rt.wrap(rt.deleteUser))

	return rt.router
}
