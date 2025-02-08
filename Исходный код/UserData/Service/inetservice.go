package userservice

import (
	"fmt"
	"strconv"
	"strings"
)

type IP4 struct {
	ip4 string
}

func (ip4 *IP4) IP4(ips string) error {
	ipsplit := strings.Split(ips, ".")
	var err error
	if len(ipsplit) != 4 {
		err = fmt.Errorf("некорректный ip4 адрес")
	}
	for _, k := range ipsplit {
		_, errs := strconv.ParseUint(k, 10, 8)
		if errs != nil {
			err = fmt.Errorf("некорректный ip4 адрес. Подробности: %v", err)
			break
		}
	}
	if err != nil {
		ip4.ip4 = ips
	}
	return err
}

func (ip4 *IP4) GeoInfo() (string, error) {
	return "Vologda/Russia", nil
}
