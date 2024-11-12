package database

import "errors"

// user
var ErrUserNotFound = errors.New("user requested not found")

// transaction
var ErrTransaction = errors.New("error with the transaction")

// chat
var ErrNotInChat = errors.New("you are not in the chat")

// message
var ErrMessageEmpty = errors.New("the message must have photo or text")