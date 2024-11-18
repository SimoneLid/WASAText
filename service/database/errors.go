package database

import "errors"

// user
var ErrUserNotFound = errors.New("the user requested was not found")
var ErrUsernameLength = errors.New("the username must contain 3-16 characters")
var ErrDifferentUser = errors.New("information of a user can be changed only by the user himself")
var ErrUsernameAlreadyExists = errors.New("the username is already in use by another user")

// transaction
var ErrTransaction = errors.New("error with the transaction")

// chat
var ErrChatNotFound = errors.New("the chat requested was not found")
var ErrNotInChat = errors.New("the user performing an action on a chat must be in the chat")
var ErrLessTwoUserInChat = errors.New("the chat must contain at least two users")
var ErrOperationGroupOnly = errors.New("this operation can be done only on groups")
var ErrGroupNameLength = errors.New("the group name can't be empty")

// message
var ErrMessageEmpty = errors.New("the message must have photo or text")
var ErrMessNotFound = errors.New("the message requested was not found")
