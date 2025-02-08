package Sessia

import (
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

func (d *Device) Field(s *sql.DB) error {
	row, err := s.Query(getRowDevice, d.DeviceID)
	if err != nil {
		return err
	}
	err = row.Scan(&d)
	return err
}

func (d *Device) Flush(s *sql.DB) error {
	_, err := s.Exec(insertRowDevice, d.DeviceID, d.OS, d.MAC, d.HostName)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &d)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sessia) GenTocken(uid userID.Storage) {
	part1 := Utility.SHA256(s.DeviceID.HostName)
	part2 := Utility.SHA384(s.History.DateUse)
	part3 := Utility.MD5Hash(s.DeviceID.MAC)
	part4 := Utility.SHA512(s.DeviceID.DeviceID)
	part5 := "?" + strconv.FormatUint(uint64(uid.IDConversion()), 10)
	s.Tocken = fmt.Sprintf("%v%v%v%v%v", part1, part2, part3, part4, part5)

}
