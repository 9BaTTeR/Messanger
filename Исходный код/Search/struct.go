package search

import (
	"ServerApp/UserData/Sessia"
	"encoding/json"
)

type Request struct {
	Dialogs     string        `json:"dialogs"`
	Search      string        `json:"search"`
	TakeDialogs uint16        `json:"takedialogs"`
	SkipDialogs uint          `json:"skipdialogs"`
	Take        uint16        `json:"take"`
	Skip        uint          `json:"skip"`
	Operation   string        `json:"operation"`
	Sessia      Sessia.Sessia `json:"sessia"`
}

func (m *Request) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &m)
	if err != nil {
		return err
	}
	return nil
}
func (m Request) Compose() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

const (
	searchindialogs = "inDialogs"
	searchAny       = "everyWhere"
)
