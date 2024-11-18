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
	err = tx.QueryRow("INSERT INTO Message(ChatId,UserId,Text,Photo,IsForwarded) VALUES(?,?,?,?,false) RETURNING MessageId",chatid,userperformingid,text,photo).Scan(&messageid)
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
		_, err = tx.Exec("INSERT OR IGNORE INTO ChatUser(UserId,ChatId) VALUES(?,?)",userid,chatid)
		if err != nil{
			errtx := tx.Rollback()
			if errtx != nil{
				return ErrTransaction
			}
			return err
		}
	}
	
	// commit the transaction
	err = tx.Commit()
	if err != nil{
		return ErrTransaction
	}


	return err
}


func (db *appdbimpl) IsUserInChat(chatid int, userid int) (bool, error) {

	var idlist []int
	// takes the users in the chat
	rows, err := db.c.Query(`SELECT UserId FROM ChatUser WHERE ChatId=?`,chatid)
	if err != nil{
		return false, err
	}
	
	defer rows.Close()

	if !rows.Next(){
		return false, ErrChatNotFound
	}

	// add all the ids to the list
	for{
		var tempid int
		err = rows.Scan(&tempid)
		if err != nil{
			return false, err
		}
		idlist = append(idlist, tempid)

		if !rows.Next(){
			break
		}
	}
	if rows.Err() != nil{
		return false, err
	}

	// check if the user is in the chat
	userinchat := false
	for i:=0;i<len(idlist);i++{
		if idlist[i]==userid{
			userinchat = true
		}
	}


	return userinchat, err
}


func (db *appdbimpl) IsGroup(chatid int) (bool, error) {
	var name sql.NullString
	err := db.c.QueryRow("SELECT ChatName FROM Chat WHERE ChatId=?",chatid).Scan(&name)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return false, ErrChatNotFound
		}
		return false, err
	}

	if !name.Valid{
		return false, err
	}

	return true, err
}

func (db *appdbimpl) ChangeGroupName(chatid int, groupname string) error {

	res, err := db.c.Exec("UPDATE Chat SET ChatName=? WHERE ChatId=?",groupname,chatid)
	if err != nil{
		return err
	}
	
	// check if the row effected are 0 which mean the chat don't exists
	eff, err := res.RowsAffected()
	if err != nil{
		return err
	}

	if eff==0{
		return ErrChatNotFound
	}

	return err
}


func (db *appdbimpl) ChangeGroupPhoto(chatid int, photo string) error {

	res, err := db.c.Exec("UPDATE Chat SET ChatPhoto=? WHERE ChatId=?",photo,chatid)
	if err != nil{
		return err
	}
	
	// check if the row effected are 0 which mean the chat don't exists
	eff, err := res.RowsAffected()
	if err != nil{
		return err
	}

	if eff==0{
		return ErrChatNotFound
	}

	return err
}