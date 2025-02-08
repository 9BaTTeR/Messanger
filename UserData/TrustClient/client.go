package trustclient

import (
	rr "ServerApp/Responces"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	keypath    = "security/rsa.key"
	hostemail  = "rolladie@distet.tech"
	commonname = "RollaDie"
)

const (
	selectRowHash = "select NameClient, MinVersion, MaxVersion, Active from CERTS where HashCert like ?;"
)

const (
	excludedRequest = "Certificate UpdateCertificate"
)

func (c Client) verifed() (rr.VerifedResponce, error) {
	resp := rr.VerifedResponce{}
	resp.Responce = rr.Responce{}.InternalError("Внутренняя ошибка")
	sql, err := openDB()
	if err != nil {
		return resp, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(selectRowHash, c.HashCert)
	if err != nil {
		resp.Responce = resp.Responce.Forbidden("Ваш сертификат не находится в системе доверенных.")
		return resp, fmt.Errorf("сбой чтения БД. Подробности: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		resp.Responce = resp.Responce.BadRequest("Неверный хэш-сертификата")
		return resp, nil
	}
	cert := Certificate{}
	err = rows.Scan(&cert.Name, &cert.MaxVersion, &cert.MinVersion, &cert.verifed)
	if err != nil {
		return resp, fmt.Errorf("ошибка чтения записи из БД. Подробности: %w", err)
	}
	if cert.MinVersion < c.Version || cert.MaxVersion > c.Version {
		resp.Responce = resp.Responce.BadRequest("версия не принадлежит сертификату")
		return resp, nil
	}
	if cert.Name != c.Name {
		resp.Responce = resp.Responce.BadRequest("сертификат не принадлежит вашему приложению")
		return resp, nil
	}
	active := cert.verifed
	if err != nil {
		return resp, fmt.Errorf("ошибка чтения флага активности сертификата. Подробности: %w", err)
	}
	if !active {
		resp.Responce = resp.Responce.BadRequest("сертификат неактивирован или был заблокирован")
		return resp, nil
	}
	resp.Responce = resp.Responce.OK("Сертификат валиден.")
	return resp, nil
}

func ConnVerify(content []byte, pathUrl string) (bool, rr.VerifedResponce, error) {
	excludes := strings.Split(excludedRequest, " ")

	if slices.Contains(excludes, pathUrl) {
		return true, rr.VerifedResponce{}, nil
	}
	c := EmptyJson{}
	err := json.Unmarshal(content, &c)
	if err != nil {
		return false, rr.VerifedResponce{}, err
	}
	resp, err := c.EmptySessia.Client.verifed()
	if resp.Responce.Code < 200 || resp.Responce.Code > 300 {
		return false, resp, err
	}
	if err != nil {
		return false, resp, err
	}
	return true, resp, err
}
