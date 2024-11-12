package database


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
		return 0, err
	}
	return userid, err
}


func (db *appdbimpl) GetUsernameFromId(userid int) (string, error) {

	var username string
	err := db.c.QueryRow(`SELECT Username FROM User WHERE UserId=?`,userid).Scan(&username)
	if err != nil{
		return "", err
	}
	return username, err
}