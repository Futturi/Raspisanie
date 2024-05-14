package entities

type Rm struct {
	Name  string `json:"name" db:"name"`
	Id    string `json:"id" db:"id"`
	Count string `json:"count"`
}

type CreateRoomReq struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type Message struct {
	ClientId string `json:"clientid"`
	Content  string `json:"content" db:"text"`
	RoomID   string `json:"roomId"`
	Username string `json:"username" db:"creator"`
	Data     int    `json:"time" db:"date"`
}
