package database

import (
	"database/sql"
	"errors"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/mattn/go-sqlite3"
)


func (db *appdbimpl) InsertUser(username string) (int, error) {

	// check if the user already exists
	var user_exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM User WHERE Username=?)`,username).Scan(&user_exists)
	if err != nil{
		return 0, err
	}

	
	var id int
	if !user_exists{
		// insert the user in db if not exists, returning the id
		err = db.c.QueryRow(`INSERT INTO User(Username,Photo) VALUES(?,"prova.png") RETURNING UserId`, username).Scan(&id)
	}else{
		// take the id of the already existing user
		err = db.c.QueryRow(`SELECT UserId FROM User WHERE Username=?`,username).Scan(&id)
	}
	if err != nil{
		return 0, err
	}

	return id, err
}


func (db *appdbimpl) GetIdFromUsername(username string) (int, error) {

	var userid int
	err := db.c.QueryRow(`SELECT UserId FROM User WHERE Username=?`,username).Scan(&userid)
	if err != nil{
		if errors.Is(err,sql.ErrNoRows){
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return userid, err
}


func (db *appdbimpl) GetUsernameFromId(userid int) (string, error) {

	var username string
	err := db.c.QueryRow(`SELECT Username FROM User WHERE UserId=?`,userid).Scan(&username)
	if err != nil{
		if errors.Is(err,sql.ErrNoRows){
			return "", ErrUserNotFound
		}
		return "", err
	}
	return username, err
}


func (db *appdbimpl) ChangeUsername(userid int, username string) error {

	res, err := db.c.Exec("UPDATE User SET Username=? WHERE UserId=?",username,userid)
	if err != nil{
		if strings.Contains(err.Error(),sqlite3.ErrConstraintUnique.Error()){
			return ErrUsernameAlreadyExists
		}
		return err
	}
	
	// check if the row effected are 0 which mean the user don't exists
	eff, err := res.RowsAffected()
	if err != nil{
		return err
	}

	if eff==0{
		return ErrUserNotFound
	}

	return err
}


func (db *appdbimpl) ChangeUserPhoto(userid int, photo string) error {

	res, err := db.c.Exec("UPDATE User SET Photo=? WHERE UserId=?",photo,userid)
	if err != nil{
		return err
	}
	
	// check if the row effected are 0 which mean the user don't exists
	eff, err := res.RowsAffected()
	if err != nil{
		return err
	}

	if eff==0{
		return ErrUserNotFound
	}

	return err
}


func (db *appdbimpl) IsUserInChat(chatid int, userid int) (bool, error) {

	userinchat:=false
	err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM ChatUser WHERE ChatId=? AND UserId=?)`,chatid,userid).Scan(&userinchat)
	if err != nil{
		return false, err
	}

	return userinchat, err
}

func (db *appdbimpl) SetLastAccess(userid int) error {

	_, err := db.c.Exec(`UPDATE User SET LastAccess=CURRENT_TIMESTAMP WHERE UserId=?`,userid)
	if err != nil{
		return err
	}

	return err
}


func (db *appdbimpl) SetLastRead(userid int, chatid int) error {

	_, err := db.c.Exec(`UPDATE ChatUser SET LastRead=CURRENT_TIMESTAMP WHERE UserId=? AND ChatId=?`,userid,chatid)
	if err != nil{
		return err
	}

	return err
}

func (db *appdbimpl) SearchUsers(username string) ([]components.User, error) {
	
	userlist := []components.User{}
	userrows, err := db.c.Query(`SELECT * FROM User WHERE Username LIKE ?||'%'`,username)
	if err != nil{
		return userlist, err
	}
	
	defer userrows.Close()

	// cicle for all the users
	for userrows.Next(){
		var user components.User
		err = userrows.Scan(&user.UserId,&user.Username,&user.Photo,&user.LastAccess)
		if err != nil{
			return userlist, err
		}

		// append the user to userlist
		userlist = append(userlist, user)
	}

	if userrows.Err() != nil{
		return userlist, err
	}
	
	return userlist, err
}