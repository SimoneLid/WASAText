package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login route
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User routes
	rt.router.PUT("/users/:user_id/name", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:user_id/photo", rt.wrap(rt.setMyPhoto))

	// Chats routes
	rt.router.POST("/newchat", rt.wrap(rt.createChat))
	rt.router.PUT("/chats/:chat_id/users", rt.wrap(rt.addToGroup))
	rt.router.PUT("/chats/:chat_id/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/chats/:chat_id/photo", rt.wrap(rt.setGroupPhoto))

	// Message routes
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/chats/:chat_id/forwardedmessages", rt.wrap(rt.forwardMessage))


	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
