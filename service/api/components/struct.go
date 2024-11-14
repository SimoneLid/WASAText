package components

type Username struct{
	Username string `json:"username"`
}

type UsernameList struct{
	UsernameList []string `json:"usernamelist"`
}

type UserId struct{
	UserId int `json:"userid"`
}

type User struct{
	UserId int `json:"userid"`
	Username string `json:"username"`
	Photo string `json:"photo"`
	LastAccess string `json:"lastaccess"`
}

type ChatMessId struct{
	ChatId int `json:"chatid"`
	MessageId int `json:"messageid"`
}

type ChatCreation struct{
	UsernameList []string `json:"usernamelist"`
	GroupName string `json:"groupname"`
	GroupPhoto string `json:"groupphoto"`
	FirstMessage MessageToSend `json:"firstmessage"`
}

type MessageToSend struct{
	IsForwarded bool `json:"isforwarded"`
	Text string `json:"text"`
	Photo string `json:"photo"`
}