package dialogs

import (
	msg "ServerApp/Messages"
	"ServerApp/UserData/Sessia"
	"database/sql"
)

type sqlInstance struct {
	instance *sql.DB
}

//Диалоги

type Members struct {
	IdUsers  []uint64 `json:"idUsers" sqlite:"IdUsers"`
	Role     []Role   `json:"role" sqlite:"Role"`
	DateJoin []string `json:"dateJoin" sqlite:"DateJoin"`
	Notice   []bool   `json:"notice" sqlite:"Notice"`
}

type Role struct {
	Name        string `json:"name" sqlite:"Name"`
	Description string `json:"description" sqlite:"Description"`
}

type Dialog struct {
	Hash          string      `json:"pkDialog" sqlite:"Hash"`
	Name          string      `json:"name" sqlite:"Name"`
	Private       string      `json:"private" sqlite:"Private"`
	CreatedAt     string      `json:"-" sqlite:"CreatedAt"`
	Photo         string      `json:"photo" sqlite:"Photo"`
	idUserCreator uint64      `json:"-" sqlite:"-"`
	IdUsers       []uint64    `json:"idUsers" sqlite:"IdUsers"`
	sql           sqlInstance `sqlite:"-" json:"-"`
}

type Link struct {
	Link     string `json:"-" sqlite:"Link"`
	IdUser   uint64 `json:"idUser" sqlite:"IdUser"`
	Duration string `json:"duration" sqlite:"Duration"`
	Count    uint64 `json:"count" sqlite:"Count"`
}

type Bans struct {
	IdUser   uint64 `json:"idUser" sqlite:"IdUser"`
	Reason   string `json:"reason" sqlite:"Reason"`
	Duration string `json:"duration" sqlite:"Duration"`
	BanBy    uint64 `json:"banBy" sqlite:"BanBy"`
}

type Settings struct {
	NameSetting []string `json:"nameSetting" sqlite:"NameSetting"`
	Value       []string `json:"value" sqlite:"Value"`
	Role        []Role   `json:"role" sqlite:"Role"`
}

type Pinned struct {
	OMessages msg.OMessages `json:"oMessages" sqlite:"OMessages"`
	Date      string        `json:"date" sqlite:"Date"`
	IdUser    uint64        `json:"idUser" sqlite:"IdUser"`
}

type DialogRequest struct {
	Dialog Dialog        `json:"dialog" sqlite:"-"`
	Sessia Sessia.Sessia `json:"sessia" sqlite:"-"`
}

type KickResponce struct {
	Whom     uint64        `json:"user" sqlite:"-"`
	PkDialog string        `json:"pkDialog" sqlite:"-"`
	Sessia   Sessia.Sessia `json:"sessia" sqlite:"-"`
}

// type CDRequest struct {
// 	NameDialog string        `json:"nameDialog" sqlite:"-"`
// 	PkDialog   string        `json:"pkDialog" sqlite:"-"`
// 	Photo      string        `json:"imageKey" sqlite:"-"`
// 	NewUsers   []uint64      `json:"addUser" sqlite:"-"`
// 	Privacy    string        `json:"privacy" sqlite:"-"`
// 	Sessia     Sessia.Sessia `json:"sessia" sqlite:"-"`
// }
