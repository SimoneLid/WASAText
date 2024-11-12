package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/mattn/go-sqlite3"
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
		return 0, 0, err
	}

	// insert the chat in db
	var chatid int
	err = tx.QueryRow(`INSERT INTO Chat(ChatName,ChatPhoto) VALUES(?,?) RETURNING ChatId`,groupname,groupphoto).Scan(&chatid)
	if err != nil{
		tx.Rollback()
		return 0, 0, err
	}

	// insert all the user in ChatUser
	for i:=0; i<len(chat.UsernameList);i++{

		// takes the id of the user
		var userid int
		userid, err = db.GetIdFromUsername(chat.UsernameList[i])
		if err != nil{
			tx.Rollback()
			return 0, 0, ErrUserNotFound
		}

		// create a row in ChatUser
		_, err = tx.Exec("INSERT INTO ChatUser(UserId,ChatId) VALUES(?,?)",userid,chatid)
		if err != nil{
			tx.Rollback()
			return 0, 0, err
		}
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

	// if the error is the check failed means that the message doesn't have neither text nor photo
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintCheck{
        return 0, 0, ErrMessageEmpty
    } else if err != nil {
        return 0, 0, err
    }
	


	// commit the transaction
	err = tx.Commit()
	if err != nil{
		return 0, 0, err
	}

	return chatid, messageid, err
}


func (db *appdbimpl) AddUsersToGroup(usernamelist []string, chatid int) error {

	var user_exists bool
	var err error
	for i:=0; i<len(usernamelist);i++{

		// check if the user exists
		err = db.c.QueryRow(`SELECT EXISTS(SELECT * FROM User WHERE Username=?)`,usernamelist[i]).Scan(&user_exists)
		if err != nil{
			return err
		}

		// if the user doesn't exist raise an error
		if !user_exists{
			err=errors.New("user not found")
			return err
		}
		
		// if the user exists takes his id
		var userid int
		userid, err = db.GetIdFromUsername(usernamelist[i])
		if err != nil{
			return err
		}

		// create a row in ChatUser
		_, err = db.c.Exec("INSERT INTO ChatUser(UserId,ChatId) VALUES(?,?)",userid,chatid)
		if err != nil{
			return err
		}
	} 

	return err
}