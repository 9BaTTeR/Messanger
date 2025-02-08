package trustclient

type Client struct {
	HashCert string `json:"hashcert" sqlite:"-"`
	Name     string `json:"nameclient" sqlite:"-"`
	Version  uint64 `json:"version" sqlite:"-"`
}

type Certificate struct {
	Hash       string `json:"hash" sqlite:"HashCert"`
	Cert       []byte `json:"-" sqlite:"-"`
	Name       string `json:"name" sqlite:"NameClient"`
	MinVersion uint64 `json:"version" sqlite:"MinVersion"`
	MaxVersion uint64 `json:"-" sqlite:"MaxVersion"`
	Email      string `json:"email" sqlite:"-"`
	verifed    bool   `json:"-" sqlite:"Active"`
	id         uint64 `json:"-" sqlite:"-"`
}

type UpdateCertificate struct {
	Hash       string `json:"hash"`
	NewVersion uint64 `json:"newversion"`
}

type EmptySessia struct {
	Client Client `json:"client"`
}

type EmptyJson struct {
	EmptySessia EmptySessia `json:"sessia"`
}
