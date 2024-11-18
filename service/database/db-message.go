package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) InsertMessage(message components.MessageToSend, isforwarded bool, chatid int, userperformingid int) (int, error) {

	// check if there is a text in message
	var text sql.NullString
	text.Valid=true
	if len(message.Text)==0{
		text.Valid=false
	}else{
		text.String = message.Text
	}

	// check if there is a photo in message
	var photo sql.NullString
	photo.Valid=true
	if len(message.Photo)==0{
		photo.Valid=false
	}else{
		photo.String = message.Photo
	}

	var messageid int
	err := db.c.QueryRow(`INSERT INTO Message(ChatId,UserId,Text,Photo, IsForwarded) VALUES(?,?,?,?,?) RETURNING MessageId`,chatid,userperformingid,text,photo,isforwarded).Scan(&messageid)
	if err != nil{
		return 0, err
	}

	return messageid, err
}

func (db *appdbimpl) GetMessage(messageid int) (components.Message, error) {

	var message components.Message
	var text sql.NullString
	var photo sql.NullString
	err := db.c.QueryRow(`SELECT * FROM Message WHERE MessageId=?`,messageid).Scan(&message.MessageId,&message.ChatId,&message.UserId,&text,&photo,&message.IsForwarded,&message.TimeStamp)
	if err != nil{
		if errors.Is(err,sql.ErrNoRows){
			return message, ErrMessNotFound
		}
		return message, err
	}

	// copies the values that can be NULL
	if text.Valid{
		message.Text=text.String
	}
	if photo.Valid{
		message.Photo=photo.String
	}


	return message, err
}