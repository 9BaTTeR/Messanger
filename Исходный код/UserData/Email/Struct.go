package goEmail

import (
	"ServerApp/UserData/Sessia"
	"database/sql"
)

type Request struct {
	Code   string        `json:"code" sqlite:"-"`
	Email  string        `json:"email" sqlite:"-" `
	Sessia Sessia.Sessia `json:"sessia"`
	sql    *sql.DB       `json:"-"`
}

type Email struct {
	Email string `sqlite:"Email"`
	ID    uint64 `sqlite:"ID"`
}
