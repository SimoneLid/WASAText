package database

import "errors"

// user
var ErrUserNotFound = errors.New("the user requested was not found")
var ErrUsernameLength = errors.New("the username must contain 3-16 characters")
var ErrAlreadyInChat = errors.New("the user requested is already in the chat")

// transaction
var ErrTransaction = errors.New("error with the transaction")

// chat
var ErrChatNotFound = errors.New("the chat requested was not found")
var ErrNotInChat = errors.New("the user performing an action on a chat must be in it")
var ErrLessTwoUserInChat = errors.New("the chat must contain at least two users")

// message
var ErrMessageEmpty = errors.New("the message must have photo or text")