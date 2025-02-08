package getter

import "ServerApp/UserData/Sessia"

// Входящая структура вытягивания диалогов

type GetDialogs struct {
	Date   string        `json:"date" sqlite:"-"`
	Sessia Sessia.Sessia `json:"sessia" sqlite:"-"`
	Take   uint64        `json:"take"`
	Skip   uint64        `json:"skip"`
}

// Входящая структура вытягивания сообщений

type GetMessages struct {
	OMKey    string        `json:"msgKey" sqlite:"-"`
	PkDialog string        `json:"pkDialog" sqlite:"-"`
	Date     string        `json:"date" sqlite:"-"`
	Take     uint64        `json:"take"`
	Skip     uint64        `json:"skip"`
	Sessia   Sessia.Sessia `json:"sessia" sqlite:"-"`
}

// Входящая структура вытягивания медиа-контента

type GetMedia struct {
	Hash   string        `json:"hash" sqlite:"-"`
	Sessia Sessia.Sessia `json:"sessia" sqlite:"-"`
}

// Входящая структура вытягивания пользователя

type GetUser struct {
	Id       uint64        `json:"id" sqlite:"-"`
	PkDialog string        `json:"pkDialog" sqlite:"-"`
	Sessia   Sessia.Sessia `json:"sessia" sqlite:"-"`
}

// Входящая структура вытягивания всех медиа диалога

type GetAllMediaDialog struct {
	PkDialog string        `json:"pkDialog"`
	Take     uint64        `json:"take"`
	Skip     uint64        `json:"skip"`
	Sessia   Sessia.Sessia `json:"sessia" sqlite:"-"`
}
