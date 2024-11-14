package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) InsertChat(chat components.ChatCreation, userperformingid int) (int, int, error) {
	
	// check if is a group
	var groupname sql.NullString
	var groupphoto sql.NullString
	groupname.Valid=true
	groupphoto.Valid=true
	if len(chat.GroupName)==0{
		groupname.Valid=false
		groupphoto.Valid=false
	}else{
		groupname.String=chat.GroupName
		// set the default image if not specified
		if len(chat.GroupPhoto)==0{
			groupphoto.String="prova.png"
		}else{
			groupphoto.String=chat.GroupPhoto
		}
	}

	// start a transaction
	tx, err :=db.c.Begin()
	if err != nil{
		return 0, 0, ErrTransaction
	}

	// insert the chat in db
	var chatid int
	err = tx.QueryRow(`INSERT INTO Chat(ChatName,ChatPhoto) VALUES(?,?) RETURNING ChatId`,groupname,groupphoto).Scan(&chatid)
	if err != nil{
		errtx := tx.Rollback()
		if errtx != nil{
			return 0, 0, ErrTransaction
		}
		return 0, 0, err
	}


	// insert all the user in ChatUser and also check if the user performing the action is in the list
	userinchat := false
	for i:=0; i<len(chat.UsernameList);i++{

		// takes the id of the user
		var userid int
		userid, err = db.GetIdFromUsername(chat.UsernameList[i])
		if err != nil{
			return 0, 0, ErrUserNotFound
		}

		// create a row in ChatUser
		if userid==userperformingid{
			userinchat = true
			// sets the LastRead for the user creating the chat
			_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId,LastRead) VALUES(?,?,CURRENT_TIMESTAMP)",userid,chatid)
			if err != nil{
				return 0, 0, err
			}
		}else{
			// for the other user LastRead is not set
			_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId) VALUES(?,?)",userid,chatid)
			if err != nil{
				return 0, 0, err
			}
		}
	}
	
	// checks userinchat
	if !userinchat{
		errtx := tx.Rollback()
		if errtx != nil{
			return 0, 0, ErrTransaction
		}
		return 0, 0, ErrNotInChat
	}


	// check if there is a text in message
	var text sql.NullString
	text.Valid=true
	if len(chat.FirstMessage.Text)==0{
		text.Valid=false
	}else{
		text.String = chat.FirstMessage.Text
	}

	// check if there is a photo in message
	var photo sql.NullString
	photo.Valid=true
	if len(chat.FirstMessage.Photo)==0{
		photo.Valid=false
	}else{
		photo.String = chat.FirstMessage.Photo
	}
	

	// insert the first message
	var messageid int
	err = tx.QueryRow("INSERT INTO Message(ChatId,UserId,Text,Photo,IsForwarded) VALUES(?,?,?,?,?) RETURNING MessageId",chatid,userperformingid,text,photo,chat.FirstMessage.IsForwarded).Scan(&messageid)
	if err != nil {
		errtx := tx.Rollback()
		if errtx != nil{
			return 0, 0, ErrTransaction
		}
        return 0, 0, err
    }

	// commit the transaction
	err = tx.Commit()
	if err != nil{
		return 0, 0, ErrTransaction
	}

	return chatid, messageid, err
}


func (db *appdbimpl) AddUsersToGroup(usernamelist []string, chatid int) error {
	// start a transaction
	tx, err :=db.c.Begin()
	if err != nil{
		return ErrTransaction
	}

	
	// insert all the user in ChatUser
	for i:=0; i<len(usernamelist);i++{

		// takes the id of the user
		var userid int
		userid, err = db.GetIdFromUsername(usernamelist[i])
		if err != nil{
			errtx := tx.Rollback()
			if errtx != nil{
				return ErrTransaction
			}
			return ErrUserNotFound
		}

		// create a row in ChatUser
		_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId) VALUES(?,?)",userid,chatid)
		if err != nil{
			errtx := tx.Rollback()
			if errtx != nil{
				return ErrTransaction
			}
			return ErrAlreadyInChat
		}
	}
	
	// commit the transaction
	err = tx.Commit()
	if err != nil{
		return ErrTransaction
	}


	return err
}


func (db *appdbimpl) GetChatComponents(chatid int) ([]int, error) {

	var idlist []int
	// takes the users in the chat
	rows, err := db.c.Query(`SELECT UserId FROM ChatUser WHERE ChatId=?`,chatid)
	if errors.Is(err,sql.ErrNoRows){
		return idlist, ErrChatNotFound
	}
	if err != nil{
		return idlist, err
	}

	defer rows.Close()

	// add all the ids to the list
	for rows.Next(){
		var tempid int
		err = rows.Scan(&tempid)
		if err != nil{
			return idlist, err
		}
		idlist = append(idlist, tempid)
	}
	return idlist, err
}