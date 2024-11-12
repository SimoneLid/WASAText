/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User
	InsertUser(username string ) (int, error)
	GetIdFromUsername(username string) (int, error)
	GetUsernameFromId(userid int) (string, error)

	// Chat
	InsertChat(chat components.ChatCreation, userperformingid int) (int, int, error)
	AddUsersToGroup(usernamelist []string, chatid int) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	User := `CREATE TABLE IF NOT EXISTS User(
				UserId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				Username TEXT NOT NULL UNIQUE,
				Photo TEXT NOT NULL,
				LastAccess DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CHECK(LENGTH(Username)>=3 AND LENGTH(Username)<=16)
				);`


	Chat := `CREATE TABLE IF NOT EXISTS Chat(
				ChatId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				ChatName TEXT,
				ChatPhoto TEXT,
				CHECK ((ChatName IS NULL AND ChatPhoto IS NULL) OR (ChatName IS NOT NULL AND ChatPhoto IS NOT NULL))
				);`


	Message := `CREATE TABLE IF NOT EXISTS Message(
					MessageId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					ChatId INTEGER NOT NULL,
					UserId INTEGER NOT NULL,
					Text TEXT,
					Photo TEXT,
					IsForwarded BOOLEAN NOT NULL,
					TimeStamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CHECK (Text IS NOT NULL OR Photo IS NOT NULL),
					FOREIGN KEY(ChatId) REFERENCES Chat(ChatId),
					FOREIGN KEY(UserId) REFERENCES User(UserId)
					);`


	Comment := `CREATE TABLE IF NOT EXISTS Comment(
					MessageId INTEGER NOT NULL,
					UserId INTEGER NOT NULL,
					Emoji TEXT,
					PRIMARY KEY(MessageId,UserId),
					FOREIGN KEY(MessageId) REFERENCES Message(MessageId),
					FOREIGN KEY(UserId) REFERENCES User(UserId)
					);`


	ChatUser := `CREATE TABLE IF NOT EXISTS ChatUser(
					UserId INTEGER NOT NULL,
					ChatId INTEGER NOT NULL,
					TimeAdded DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					LastRead DATETIME,
					PRIMARY KEY(UserId,ChatId),
					CHECK((LastRead>TimeAdded) OR (LastRead IS NULL)),
					FOREIGN KEY(ChatId) REFERENCES Chat(ChatId),
					FOREIGN KEY(UserId) REFERENCES User(UserId)
					);`


	_, err := db.Exec(User)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Chat)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Message)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Message)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(Comment)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(ChatUser)
	if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
	}





	// Inizio test

	// Decommenta per testare
	_, err = db.Exec(`
		INSERT OR IGNORE INTO User (Username, Photo, LastAccess) VALUES ('Sim', 'sim_photo.jpg', '2024-11-01 08:30:00');
		INSERT OR IGNORE INTO User (Username, Photo) VALUES ('Ivan','photo.png');
	`)
	if err != nil {
		return nil, fmt.Errorf("error inserting into User: %w", err)
	}

	/*
	// Inserimenti per la tabella Chat
	_, err = db.Exec(`
		INSERT INTO Chat (ChatName, ChatPhoto, IsGroup) VALUES (NULL, NULL, 0);
		INSERT INTO Chat (ChatName, ChatPhoto, IsGroup) VALUES ('Group Chat 1', 'group1_photo.jpg', 1);
		INSERT INTO Chat (ChatName, ChatPhoto, IsGroup) VALUES ('Group Chat 2', 'group2_photo.jpg', 1);
	`)
	if err != nil {
		return nil, fmt.Errorf("error inserting into Chat: %w", err)
	}

	// Inserimenti per la tabella Message
	_, err = db.Exec(`
		INSERT INTO Message (ChatId, UserId, Text, Photo, IsPhoto, IsForwarded, TimeStamp) 
		VALUES (1, 1, 'Hello, this is a private message.', NULL, 0, 0, '2024-11-01 08:45:00');

		INSERT INTO Message (ChatId, UserId, Text, Photo, IsPhoto, IsForwarded, TimeStamp) 
		VALUES (2, 2, NULL, 'photo1.jpg', 1, 1, '2024-11-02 09:20:00');

		INSERT INTO Message (ChatId, UserId, Text, Photo, IsPhoto, IsForwarded, TimeStamp) 
		VALUES (2, 3, 'Welcome to the group chat!', NULL, 0, 0, '2024-11-03 10:05:00');
	`)
	if err != nil {
		return nil, fmt.Errorf("error inserting into Message: %w", err)
	}

	// Inserimenti per la tabella Comment
	_, err = db.Exec(`
		INSERT INTO Comment (MessageId, UserId, Emoji) VALUES (1, 2, '👍');
		INSERT INTO Comment (MessageId, UserId, Emoji) VALUES (1, 3, '😂');
		INSERT INTO Comment (MessageId, UserId, Emoji) VALUES (2, 1, '😍');
	`)
	if err != nil {
		return nil, fmt.Errorf("error inserting into Comment: %w", err)
	}

	// Inserimenti per la tabella ChatUser
	_, err = db.Exec(`
		INSERT INTO ChatUser (UserId, ChatId, TimeAdded, LastRead) VALUES (1, 1, '2024-11-01 08:30:00', '2024-11-01 08:50:00');
		INSERT INTO ChatUser (UserId, ChatId, TimeAdded, LastRead) VALUES (2, 1, '2024-11-01 08:40:00', '2024-11-01 09:00:00');
		INSERT INTO ChatUser (UserId, ChatId, TimeAdded, LastRead) VALUES (2, 2, '2024-11-02 09:00:00', NULL);
		INSERT INTO ChatUser (UserId, ChatId, TimeAdded, LastRead) VALUES (3, 2, '2024-11-02 09:05:00', '2024-11-03 10:10:00');
		INSERT INTO ChatUser (UserId, ChatId, TimeAdded, LastRead) VALUES (1, 2, '2024-11-02 09:10:00', '2024-11-02 09:20:00');
	`)
	if err != nil {
		return nil, fmt.Errorf("error inserting into ChatUser: %w", err)
	}

	fmt.Println("Data inserted successfully!")
	*/


	return &appdbimpl{
		c: db,
	}, nil
}



func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
