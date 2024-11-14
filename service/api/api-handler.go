package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login route
	rt.router.POST("/session",rt.wrap(rt.doLogin))

	// Chats route
	rt.router.POST("/chats/newchat",rt.wrap(rt.createChat))
	rt.router.PUT("/chats/:chat_id/users",rt.wrap(rt.addToGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
