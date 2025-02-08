package message

import (
	"ServerApp/UserData/Sessia"
	"database/sql"
)

//Сообщения

type DMessages struct {
	DMKey        string `json:"msgKey" sqlite:"DMKey"`
	idUser       uint64 `json:"-" sqlite:"IdUser"`
	senderUser   uint64 `json:"-" sqlite:"-"`
	Date         string `json:"date" sqlite:"Date"`
	Content      string `json:"content" sqlite:"Content"`
	UpdateAt     string `json:"updateAt" sqlite:"UpdateAt"`
	DeletedAt    string `json:"deleteAt" sqlite:"DeletedAt"`
	ForwardedKey string `json:"forwardedKey" sqlite:"ForwardedKey"`
	PkDialog     string `json:"pkDialog" sqlite:"-"`
	Read         string `json:"read" sqlite:"-"`
	Important    string `json:"important" sqlite:"-"`
}

type OMessages struct {
	OMKey     string    `json:"-" sqlite:"OMKey"`
	Date      string    `json:"-" sqlite:"Date"`
	DMessages DMessages `json:"-" sqlite:"DMessages"`
}

type Media struct {
	DMessages DMessages `json:"infoMessage" sqlite:"DMessages"`
	Hash      []string  `json:"hash" sqlite:"Hash"`
	Order     []uint64  `json:"order" sqlite:"Order"`
}

type MessageRequest struct {
	Message     Media         `json:"message" sqlite:"-"`
	OMessage    OMessages     `json:"-" sqlite:"-"`
	Sessia      Sessia.Sessia `json:"sessia" sqlite:"-"`
	Operation   string        `json:"operation" sqlite:"-"`
	sqlInstance *sql.DB
}

// Для получения всех участников
type Members struct {
	IdUsers []uint64 `json:"idUsers" sqlite:"IdUsers"`
}
