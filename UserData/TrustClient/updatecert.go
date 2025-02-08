package trustclient

import (
	"ServerApp/Configs"
	rr "ServerApp/Responces"
	"ServerApp/Utility"
	"fmt"
)

const (
	updateRow = "UPDATE CERTS\n" +
		"SET MaxVersion = ?\n" +
		"WHERE HashCert = ?"
)

func (uc UpdateCertificate) UpdateVersion() (rr.VerifedResponce, error) {
	pathDB := Configs.BlackListDefDB()
	exists, err := Utility.Exists(pathDB)
	vr := rr.VerifedResponce{}
	vr.Responce = vr.Responce.InternalError("сбой регистрации сертификата")
	if err != nil {

		return vr, fmt.Errorf("невозможно проверить наличие БД. Подробности: %w", err)
	}
	if !exists {
		vr.Responce = vr.Responce.BadRequest("сертификат необнаружен.")
		return vr, nil
	}
	sql, err := openDB()
	if err != nil {
		return vr, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	_, err = sql.Exec(updateRow, uc.NewVersion, uc.Hash)
	if err != nil {
		return vr, fmt.Errorf("сбой обновления записи в БД. Подробности: %w", err)
	}
	return rr.VerifedResponce{}, nil
}

func (uc UpdateCertificate) Disable() (rr.VerifedResponce, error) {
	// pathDB := Configs.Path{}.FolderDB().Path() + "cert.db"
	// exists, err := Utility.Exists(pathDB)
	vr := rr.VerifedResponce{}
	vr.Responce = vr.Responce.InternalError("сбой регистрации сертификата")
	return rr.VerifedResponce{}, nil
}

func (uc UpdateCertificate) Enable() (rr.VerifedResponce, error) {
	// pathDB := Configs.Path{}.FolderDB().Path() + "cert.db"
	// exists, err := Utility.Exists(pathDB)
	vr := rr.VerifedResponce{}
	vr.Responce = vr.Responce.InternalError("сбой регистрации сертификата")
	return rr.VerifedResponce{}, nil
}
