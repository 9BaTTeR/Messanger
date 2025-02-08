package user

import (
	Tocken "ServerApp/UserData/Sessia"
	"database/sql"
)

type Password struct {
	Pass       string `json:"password" sqlite:"-"`
	DataCreate string `json:"-" sqlite:"DataCreate"`
	Activated  bool   `json:"-" sqlite:"-"`
}

type Nickname struct {
	Nickname   string `json:"login" sqlite:"-"`
	DataCreate string `json:"-" sqlite:"-"`
	Activated  bool   `json:"-" sqlite:"-"`
}

type Email struct {
	Email      string `json:"email" sqlite:"-"`
	DataCreate string `json:"-" sqlite:"-"`
	Activated  bool   `json:"-" sqlite:"-"`
	Code       string `json:"-" sqlite:"-"`
}

type Photo struct {
	Path       string `json:"-" sqlite:"-"`
	DataCreate string `json:"-" sqlite:"-"`
	Activated  bool   `json:"-" sqlite:"-"`
}

type sqlInstance struct {
	instance *sql.DB
}

type User struct {
	PK         uint64        `json:"-" sqlite:"Id"`
	Password   Password      `json:"password" sqlite:"-"`
	Login      Nickname      `json:"login" sqlite:"-"`
	Email      Email         `json:"email" sqlite:"-"`
	Photo      Photo         `json:"-" sqlite:"-"`
	Sessia     Tocken.Sessia `json:"sessia" sqlite:"-"`
	code       string        `json:"-" sqlite:"-"`
	sql        sqlInstance   `json:"-" sqlite:"-"`
	restrained bool          `json:"-" sqlite:"-"`
}

type AllUserData struct {
	History []Tocken.History
}

type Request struct {
	Login    string        `json:"nickname"`
	OldPass  string        `json:"oldPass"`
	NewPass  string        `json:"newPass"`
	Image    string        `json:"imagekey"`
	Email    string        `json:"email"`
	Sessia   Tocken.Sessia `json:"sessia"`
	user     User          `json:"-"`
	oldLogin string        `json:"-"`
	oldEmail string        `json:"-"`
	oldImage string        `json:"-"`
}

// /disable
