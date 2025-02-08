package media

import "ServerApp/UserData/Sessia"

type Media struct {
	MediaKey    string        `json:"-" sqlite:"-"`
	Extension   string        `json:"extension" sqlite:"-"`
	BytesBase64 string        `json:"bytes" sqlite:"-"`
	Bytes       []byte        `json:"-" sqlite:"-"`
	Sessia      Sessia.Sessia `json:"sessia" sqlite:"-"`
}

