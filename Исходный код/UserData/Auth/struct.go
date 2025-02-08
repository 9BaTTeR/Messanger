package auth

import s "ServerApp/UserData/Sessia"

const (
	tocken  = "tocken"
	logpass = "logpass"
)

type Auth struct {
	AT       AuthTocken `json:"sessia"`
	ap       AuthPass
	authtype string
}

type AuthTocken struct {
	Tocken string `json:"tocken"`
	Device Device `json:"device"`
}

type Device struct {
	DeviceID string `json:"deviceid"`
}

type AuthPass struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Sessia   s.Sessia `json:"sessia"`
}
