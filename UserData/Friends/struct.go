package friends

import (
	"ServerApp/UserData/Sessia"
	"database/sql"
	"time"
)

type Friend struct {
	ID          uint64 `json:"id" sqlite:"ID"`
	Description string `json:"description" sqlite:"DESCRIPTION"`
	DateAdd     string `json:"dateadd" sqlite:"DATEADD"`
	Photo       string `json:"photo" sqlite:"-"`
	Nickname    string `json:"nickname" sqlite:"-"`
}

type Request struct {
	Sessia      Sessia.Sessia `json:"sessia" sqlite:"-"`
	Id          uint64        `json:"id"`
	Description string        `json:"description"`
	Operation   TypeOperation `json:"operation"`
	Coming      uint8         `json:"-"`
	Take        uint64        `json:"take"`
	Skip        uint64        `json:"skip"`
	Date        time.Time     `json:"-"`
	sql         *sql.DB       `json:"-" sqlite:"-"`
}

type TypeOperation string

const (
	accept           TypeOperation = "accept"       //Done
	appends          TypeOperation = "append"       //Done
	remove           TypeOperation = "remove"       //Done
	cancel           TypeOperation = "cancel"       //DOIT
	update           TypeOperation = "update"       //Done
	getFriend        TypeOperation = "getFriend"    //Done
	getIncoming      TypeOperation = "getIncoming"  //Done
	getOutcoming     TypeOperation = "getOutcoming" //Done
	addBlackList     TypeOperation = "addBlock"     //Done
	removeBlackList  TypeOperation = "delBlock"     //done
	getListBlackList TypeOperation = "listBlocks"   //done
)
