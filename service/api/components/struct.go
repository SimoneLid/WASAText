package components

type Username struct{
	Username string `json:"username"`
}

type GroupName struct{
	GroupName string `json:"groupname"`
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

type MessageId struct{
	MessageId int `json:"messageid"`
}

type MessageToSend struct{
	Text string `json:"text"`
	Photo string `json:"photo"`
}

type Message struct{
	MessageId int `json:"messageid"`
	ChatId int `json:"chatid"`
	UserId int `json:"userid"`
	Text string `json:"text"`
	Photo string `json:"photo"`
	IsForwarded bool `json:"isforwarded"`
	TimeStamp string `json:"timestamp"`
}


type Photo struct{
	Photo string `json:"photo"`
}