package database

import (
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) InsertChat(chat components.ChatCreation) (int, error) {
	
	isgroup:=false
	// check if is a group
	if len(chat.GroupName)>0{
		isgroup=true
		// set the default image if not specified
		if len(chat.GroupPhoto)==0{
			chat.GroupPhoto="prova.png"
		}
	}


	var chatid int
	var err error
	if isgroup{
		// if is a group sets name and photo
		err = db.c.QueryRow(`INSERT INTO Chat(ChatName,ChatPhoto,IsGroup) VALUES(?,?,?) RETURNING ChatId`,chat.GroupName,chat.GroupPhoto,isgroup).Scan(&chatid)
	}else{
		// if isn't a group sets name and photo as NULL
		err = db.c.QueryRow(`INSERT INTO Chat(IsGroup) VALUES(?) RETURNING ChatId`,isgroup).Scan(&chatid)
	}
	if err != nil{
		return 0, err
	}

	return chatid, err
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