package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login route
	rt.router.POST("/session",rt.wrap(rt.doLogin))

	// Chats route
	rt.router.POST("/chats",rt.wrap(rt.createChat))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
