package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// post /chats
func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	// take the user performing the action from the bearer
	userperformingid, err := strconv.Atoi(r.Header.Get("Authorization")[7:])
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}


	// take the chat to create from the request body
	var chat components.ChatCreation
	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}

	// check if the user performing is in the chat
	userperformingname, err := rt.db.GetUsernameFromId(userperformingid)
	if err != nil{
		http.Error(w,err.Error(),http.StatusUnauthorized)
		return
	}
	if userperformingname != chat.UsernameList[0]{
		http.Error(w,database.ErrNotInChat.Error(),http.StatusUnauthorized) // 401
		return
	}

	// Inserts the chat in the database
	var chatid, messageid int
	chatid, messageid, err = rt.db.InsertChat(chat, userperformingid)
	if errors.Is(err, database.ErrTransaction) {
        http.Error(w,err.Error(),http.StatusInternalServerError) // 500
	}
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}


	var ids components.ChatMessId
	ids.ChatId=chatid
	ids.MessageId=messageid

	// set the header of the response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(ids)
}