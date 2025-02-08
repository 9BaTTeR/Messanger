package auth

import (
	responces "ServerApp/Responces"
	"database/sql"
	"fmt"

	email "ServerApp/UserData/Email"
	us "ServerApp/UserData/User"
	userid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"encoding/json"
)

const (
	brokentocken  = "Некорректный токен"
	brokendata    = "Пароль и email не верны"
	internalerror = "Внутренняя ошибка"
)

func TryAuth(jsonsource string, ip4 string) (Auth, error) {
	auth := Auth{}
	err := json.Unmarshal([]byte(jsonsource), &auth)
	if err == nil {
		auth.authtype = tocken
		if auth.AT.Tocken != "" {
			return auth, nil
		}

	}
	err = json.Unmarshal([]byte(jsonsource), &auth.ap)
	auth.authtype = logpass
	if err != nil {
		return Auth{}, err
	}
	auth.ap.Sessia.History.Ip = ip4
	auth.ap.Sessia.History.DateUse = Utility.GetTime()
	return auth, nil
}

func (a Auth) IsTockenAuth() bool {
	return a.authtype == tocken
}

func (a Auth) Auth() (responces.Resp, error) {
	fmt.Printf("%+v", a)
	user := us.User{}
	var err error
	var responce responces.AuthResponce
	if a.IsTockenAuth() {
		responce, err = a.authTocken(&user)
	} else {
		responce, err = a.authPassword(&user)
	}

	if err != nil || responce.Status.Code != 200 {
		return responce, err
	}
	return responce, nil
}

func (a Auth) authTocken(u *us.User) (responces.AuthResponce, error) {
	responce := responces.AuthResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка при авторизации"),
	}
	responce.Tocken = a.AT.Tocken
	u.Sessia.Tocken = a.AT.Tocken
	uid, err := u.Sessia.ParseID()
	if err != nil {
		responce.Status = responce.Status.BadRequest(brokentocken)
		return responce, err
	}

	u.PK = uint64(uid.IDConversion())
	err = u.Sessia.FieldSessiaFromDB()
	if err != nil {
		responce.Status = responce.Status.InternalError("Внутренний сбой при заполнении пользователей через PK")
		return responce, err
	}

	if u.Sessia.DeviceID.DeviceID != a.AT.Device.DeviceID {
		responce.Status = responce.Status.BadRequest("Токен не относится к этому deviceid")
		return responce, err
	}

	u.Sessia.Available = true

	sql, err := u.Sessia.ConnectDB()
	if err != nil {
		return responce, err
	}
	defer sql.Close()
	err = u.Sessia.History.Update(sql)
	if err != nil {
		return responce, err
	}

	responce.Status = responce.Status.OK("Токен валиден")
	return responce, nil
}

func (a Auth) authPassword(u *us.User) (responces.AuthResponce, error) {
	a.ap.Password = Utility.SaltPass(a.ap.Password)
	email := email.Email{}
	email.Email = a.ap.Email
	emErr, err := email.FindID()
	responce := responces.AuthResponce{}
	responce.Status = responce.Status.InternalError("внутренняя ошибка при авторизации по логику и паролю.")
	if err != nil {
		rsErr := responces.AuthResponce{Status: emErr}
		return rsErr, fmt.Errorf("ошибка получения ID по email. Подробности: %w", err)
	}
	uId := userid.ConvertionID(uint64(email.ID))
	u.PK = email.ID
	instance, err := u.UserSqlInstance()
	if err != nil {
		return responce, fmt.Errorf("ошибка получения соединения с БД. Подробности: %w", err)
	}
	defer instance.Close()
	result := instance.QueryRow("SELECT Login FROM USER WHERE Email = ? AND Password = ?", a.ap.Email, a.ap.Password)
	var user string
	err = result.Scan(&user)
	if err == sql.ErrNoRows  {
		responce.Status = responce.Status.BadRequest(brokendata)
		return responce, fmt.Errorf("данные не верны. Подробности: %w", err)
	}else if err != nil{
		responce.Status = responce.Status.InternalError(internalerror)
		return responce, fmt.Errorf("ошибка чтения записей из БД. Подробности: %w", err)
	}
	if user == "" {
		responce.Status = responce.Status.BadRequest("пользователь с таким логином не зарегистрирован.")
		return responce, nil
	}

	u.Sessia.History = a.ap.Sessia.History
	u.Sessia.DeviceID = a.ap.Sessia.DeviceID
	u.Sessia.Available = true
	u.Sessia.GenTocken(uId)
	responce.Tocken = u.Sessia.Tocken
	responce.Status = responce.Status.InternalError(internalerror)
	pathSessia, err := u.Sessia.PathDB()
	if err != nil {
		return responce, fmt.Errorf("ошибка получения пути. Подробности: %w", err)
	}
	sessiaDB, err := sql.Open("sqlite3", pathSessia)
	if err != nil {
		return responce, fmt.Errorf("ошибка открытия подключения сессии. Подробности: %w", err)
	}
	defer sessiaDB.Close()
	err = u.Sessia.Flush()
	if err != nil {
		return responce, fmt.Errorf("ошибка обновления записи сессии. Подробности: %w", err)
	}
	err = u.Sessia.History.Update(sessiaDB)
	if err != nil {
		return responce, fmt.Errorf("ошибка обновления записи сессии. Подробности: %w", err)
	}
	responce.Status = responce.Status.OK("Авторизация прошла успешно")
	return responce, nil
}
