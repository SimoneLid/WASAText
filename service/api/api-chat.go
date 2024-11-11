package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// post /chats
func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	// take the username from the request body
	var chat components.ChatCreation
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}

	/* Capisci come fare le 3 funzioni sotto:
	   - InsertChat
	   - AddUsersToGroup
	   - InsertMessage
	   In contemporanea, cioè se da errore 1 non esegui nessuna delle 3
	*/

	// Inserts the chat in the database
	var chatid int
	chatid, err = rt.db.InsertChat(chat)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}


	err = rt.db.AddUsersToGroup(chat.UsernameList, chatid)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	}

	/* // Inserts the message in the database
	var messageid int
	messageid, err = InsertMessage(chat.FirstMessage)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest) // 400
		return
	} */

	var ids components.ChatMessId
	ids.ChatId=chatid
	//ids.MessageId=messageid

	// set the header of the response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(ids)
}