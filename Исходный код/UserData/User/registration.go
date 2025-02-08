package user

import (
	"ServerApp/Configs"
	rr "ServerApp/Responces"
	goEmail "ServerApp/UserData/Email"
	gencode "ServerApp/UserData/GenCode"
	"ServerApp/UserData/MailSender"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"encoding/json"
	"fmt"
)

func (rg *User) Parse(source string) error {
	json.Unmarshal([]byte(source), &rg)
	return nil
}

func (rg User) Compose() ([]byte, error) {
	data, err := json.Marshal(rg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (reg User) Linter() (bool, string) {
	result, msg := LinterPass(reg.Password.Pass)
	if !result {
		return result, msg
	}
	result, msg = LinterLogin(reg.Login.Nickname)
	if !result {
		return result, msg
	}
	result, msg = LinterEmail(reg.Email.Email)
	if !result {
		return result, msg
	}
	return true, ""

}

func (reg *User) Registration(ip4 string) (rr.RegResponce, error) {
	lint, msg := reg.Linter()
	if !lint {
		return rr.RegResponce{
				Status: rr.Responce{}.BadRequest(msg),
				Tocken: "-1",
			},
			nil
	}
	reg.Password.Pass = Utility.SaltPass(reg.Password.Pass)
	userdir := Configs.UserDir()
	id := userID.LastFreeID()
	email := goEmail.Email{}
	email.Email = reg.Email.Email
	email.ID = id.IDConversion()
	reg.Sessia.History.Ip = ip4
	reg.Sessia.History.DateUse = Utility.GetTime()
	reg.PK = uint64(id.IDConversion())
	reg.Sessia.Available = true
	reg.Sessia.GenTocken(id)
	
	responce, err := email.LinkEmail()
	if err != nil || responce.Code != 200 {
		id.DecrementID()
		return rr.RegResponce{Status: responce, Tocken: "-1"}, err
	}

	pathes := Utility.Combine([]string{id.MasterDir(), id.SlaveDir()})

	err = Utility.CreateFolders(Utility.Combine([]string{userdir, pathes, "/Dialogs"}))
	if err != nil {

		id.DecrementID()
		return rr.RegResponce{Status: (rr.Responce{}).InternalError("Ошибка создание пользователя"), Tocken: "-1"}, fmt.Errorf("ошибка при создании каталога сообщений. Подробности: %v", err)
	}

	reg.code = string(gencode.GenCode())
	err = dbFilling(*reg)
	if err != nil {
		id.DecrementID()

		return rr.RegResponce{Status: (rr.Responce{}).InternalError("Ошибка создание пользователя"), Tocken: "-1"}, fmt.Errorf("ошибка наполнения базовых файлов пользователя. Подробности: %v", err)
	}
	err = MailSender.MailAuthCode(reg.Email.Email, reg.Login.Nickname, reg.code, MailSender.AboutUser(ip4, Utility.GetTime(), reg.Sessia))
	if err != nil {

		id.DecrementID()
		return rr.RegResponce{Status: (rr.Responce{}).InternalError("Ошибка создание пользователя"), Tocken: "-1"}, fmt.Errorf("ошибка отправки кода на почту пользователя. Подробности: %v", err)
	}

	return rr.RegResponce{Status: (rr.Responce{}).OK("Регистрация завершена."), Tocken: reg.Sessia.Tocken}, nil
}

func dbFilling(reg User) error {
	err := reg.CopyDB()
	if err != nil {
		return err
	}
	err = reg.addEmail()
	if err != nil {
		return err
	}
	err = reg.addLogin()
	if err != nil {
		return err
	}
	err = reg.addPassword()
	if err != nil {
		return err
	}
	err = reg.addUserDB()
	if err != nil {
		return err
	}
	err = reg.Sessia.FlushAll()
	if err != nil {
		return err
	}
	return nil
}
