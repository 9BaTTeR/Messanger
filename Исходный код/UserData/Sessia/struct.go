package Sessia

import (
	cl "ServerApp/UserData/TrustClient"
	"database/sql"
)

type Device struct {
	DeviceID string `sqlite:"DEVICEID" json:"deviceid"`
	OS       string `sqlite:"OS" json:"os"`
	MAC      string `sqlite:"MAC" json:"mac"`
	HostName string `sqlite:"HOSTNAME" json:"hostname"`
}
type Sessia struct {
	Tocken    string    `sqlite:"TOCKEN" json:"tocken"`
	Available bool      `sqlite:"AVAILABLE" json:"available"`
	DeviceID  Device    `sqlite:"DEVICEID" json:"device"`
	History   History   `sqlite:"HISTORY" json:"-"`
	Client    cl.Client `sqlite:"-" json:"client"`
	sql       *sql.DB   `sqlite:"-" json:"-"`
}
type History struct {
	DateUse string `sqlite:"DATAUSE" json:"dateuse"`
	Ip      string `sqlite:"IP" json:"ip"`
}
type EmptySessia struct {
	Sessia Sessia `json:"sessia"`
}
