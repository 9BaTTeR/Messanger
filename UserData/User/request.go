package user

import (
	rs "ServerApp/Responces"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
	"time"
)

const (
	givePass      = "Select DataCreate From PASSWORD where Password like ? and Activated like 1 LIMIT 1;"
	disablePass   = "Update PASSWORD SET Activated=0 where Activated like 1;"
	enablePass    = "Update PASSWORD SET Activated=1, DataCreate=? where Password like ?;"
	insertNewPass = "INSERT INTO PASSWORD(DataCreate,Password,Activated) VALUES(?,?,1);"
	setKeyPass    = "UPDATE USER SET Password=?;"
)

const (
	day = time.Hour * 24
)

func (r *Request) changePass() (rs.NonAddResponce, error) {
	answer := rs.NonAddResponce{Status: rs.Responce{}.InternalError("внутренняя ошибка при смене пароля")}
	result, msg := LinterPass(r.NewPass)
	if !result {
		answer.Status = answer.Status.BadRequest(msg)
	}
	r.OldPass, r.NewPass = Utility.SaltPass(r.OldPass), Utility.SaltPass(r.NewPass)
	var err error
	err = r.user.openDB()
	if err != nil {
		return answer, fmt.Errorf("сбой соединения с БД пользователя. Подробности: %w", err)
	}
	defer r.user.closeDB()
	lastchange := time.Time{}
	oldexists, newexists, err := r.passExists(&lastchange)
	if err != nil {
		return answer, fmt.Errorf("ошибка проверка паролей на наличие. Подробности: %w", err)
	}
	if !oldexists {
		answer.Status = answer.Status.BadRequest("указанный пароль неверный")
		return answer, nil
	}
	if timePass(lastchange) {
		answer.Status = answer.Status.BadRequest("нельзя слишком часто менять пароль")
		return answer, nil
	}
	_, err = r.user.sql.instance.Exec(disablePass)
	if err != nil {
		errs := r.rollbackPass()
		return answer, fmt.Errorf("ошибка отключения старого пароля. Подробности: %s. Ошибки при откате: %w", err, errs)
	}
	if newexists {
		err = r.enableNewPass()
		if err != nil {
			errs := r.rollbackPass()
			return answer, fmt.Errorf("ошибка активации нового пароля. Подробности: %s. Ошибки при откате: %w", err, errs)
		}
		answer.Status = answer.Status.OK("пароль сменён")
		return answer, nil
	}
	err = r.createNewPass()
	if err != nil {
		errs := r.rollbackPass()
		return answer, fmt.Errorf("ошибка создания нового пароля. Подробности: %s. Ошибки при откате: %w", err, errs)
	}
	_, err = r.user.sql.instance.Exec(setKeyPass, r.NewPass)
	if err != nil {
		errs := r.rollbackPass()
		return answer, fmt.Errorf("ошибка привязки нового пароля. Подробности: %s. Ошибки при откате: %w", err, errs)

	}
	answer.Status = answer.Status.OK("пароль изменён")
	return answer, nil
}

func (r Request) rollbackPass() error {
	_, err := r.user.sql.instance.Exec(enablePass, Utility.GetTime(), r.OldPass)
	return err
}

func (r Request) enableNewPass() error {
	_, err := r.user.sql.instance.Exec(enablePass, Utility.GetTime(), r.NewPass)
	return err
}

func (r Request) createNewPass() error {
	_, err := r.user.sql.instance.Exec(insertNewPass, Utility.GetTime(), r.NewPass)
	return err
}

func (r Request) passExists(t *time.Time) (oldexists bool, newexists bool, err error) {
	rows := r.user.sql.instance.QueryRow(givePass, r.OldPass)
	var tempdate string
	err = rows.Scan(&tempdate)
	if err == sql.ErrNoRows {
		err = nil
		oldexists = false
		return
	} else if err == nil {
		oldexists = true
	} else {
		return false, false, fmt.Errorf("ошибка сканирования записи. Подробности: %w", err)
	}
	*t, err = time.Parse("2006-01-02 15:04:05.99", tempdate)

	if err != nil {
		return false, false, fmt.Errorf("невозможно проверить дату создания пароля. Подробности: %w", err)
	}
	rows = r.user.sql.instance.QueryRow(givePass, r.NewPass)
	err = rows.Scan()
	if err == sql.ErrNoRows {
		err = nil
		newexists = false
	} else if err == nil {
		newexists = true
	} else {
		return false, false, fmt.Errorf("ошибка сканирования записи. Подробности: %w", err)
	}
	return
}

const (
	lastNick    = "select Login, DataCreate from LOGIN Where Activated like 1 Order by DataCreate DESC Limit 1"
	disableNick = "UPDATE LOGIN SET Activated = 0 WHERE Activated like 1"
	enableNick  = "INSERT INTO LOGIN (DataCreate, Login, Activated)\n" +
		"VALUES (?, ?, 1)\n" +
		"ON CONFLICT(Login) DO UPDATE SET\n" +
		"DataCreate = excluded.DataCreate,\n" +
		"Activated = excluded.Activated;\n"
	setKeyNick = "UPDATE USER SET Login=?"
)

func (r Request) changeNickname() (rs.NonAddResponce, error) {
	answer := rs.NonAddResponce{Status: rs.Responce{}.InternalError("внутренняя ошибка при смене никнейма")}
	result, msg := LinterLogin(r.Login)
	if !result {
		answer.Status = answer.Status.BadRequest(msg)
		return answer, nil
	}
	err := r.user.openDB()
	if err != nil {
		return answer, fmt.Errorf("сбой соединения с базой данных пользователя. Подробности: %w", err)
	}
	defer r.user.closeDB()
	lastchange := time.Time{}
	err = r.latestNickname(&lastchange)
	if err != nil {
		return answer, fmt.Errorf("сбой получения последнего, активированного никнейма. Подробности: %w", err)
	}
	if timePass(lastchange) {
		answer.Status = answer.Status.BadRequest("нельзя слишком часто менять никнейм")
		return answer, nil
	}
	_, err = r.user.sql.instance.Exec(disableNick)
	if err != nil {
		errs := r.rollbackNickname()
		return answer, fmt.Errorf("ошибка отключения логина. Подробности: %s. Ошибка отката: %w", err, errs)
	}
	_, err = r.user.sql.instance.Exec(enableNick, Utility.GetTime(), r.Login)
	if err != nil {
		errs := r.rollbackNickname()
		return answer, fmt.Errorf("ошибка активации логина. Подробности: %s. Ошибка отката: %w", err, errs)
	}
	_, err = r.user.sql.instance.Exec(setKeyNick, r.Login)
	if err != nil {
		errs := r.rollbackNickname()
		return answer, fmt.Errorf("ошибка привязки нового никнейма. Подробности: %s. Ошибки при откате: %w", err, errs)

	}
	answer.Status = answer.Status.OK("никнейм изменён")
	return answer, nil
}

func (r Request) rollbackNickname() error {
	r.user.sql.instance.Exec(enableNick, Utility.GetTime(), r.oldLogin)
	return nil
}

func (r *Request) latestNickname(t *time.Time) error {
	rows := r.user.sql.instance.QueryRow(lastNick)
	var tempdate string
	err := rows.Scan(&r.oldLogin, &tempdate)
	if err != nil {
		return fmt.Errorf("ошибка получения последнего никнейма(логина). Подробности: %w", err)
	}
	*t, err = time.Parse("2006-01-02 15:04:05.99", tempdate)
	if err != nil {
		return fmt.Errorf("невозможно проверить дату создания логина. Подробности: %w", err)
	}
	return nil
}

const (
	month = time.August
)

const (
	lastEmail        = "select Email, DataCreate from EMAIL Where Activated like 1 Order by DataCreate DESC Limit 1"
	disableAllSessia = "Update "
)

// BROKEN
func (r Request) ChangeEmail() (rs.NonAddResponce, error) {
	answer := rs.NonAddResponce{Status: rs.Responce{}.InternalError("внутренняя ошибка при смене почты")}
	lastchanges := time.Time{}
	result, msg := LinterEmail(r.Email)
	if !result {
		answer.Status = answer.Status.BadRequest(msg)
		return answer, nil
	}
	err := r.latestEmail(&lastchanges)
	if err != nil {
		return answer, fmt.Errorf("ошибка полчения последней почты. Подробности: %w", err)
	}
	monthPass(lastchanges)
	return answer, nil
}

func (r *Request) latestEmail(t *time.Time) error {

	rows := r.user.sql.instance.QueryRow(lastEmail)
	var tempdate string
	err := rows.Scan(&r.oldEmail, &tempdate)
	if err != nil {
		return fmt.Errorf("ошибка получения последней почты. Подробности: %w", err)
	}
	*t, err = time.Parse("2006-01-02 15:04:05.99", tempdate)
	if err != nil {
		return fmt.Errorf("невозможно проверить дату создания логина. Подробности: %w", err)
	}
	_ = r.rollbackEmail()
	return nil
}

func (r Request) rollbackEmail() error {
	return nil
}

const (
	lastImage    = "SELECT Path FROM PHOTO Where Activated like 1"
	disableImage = "UPDATE PHOTO SET Activated=0 WHERE Activated like 1"
	updateImage  = "INSERT INTO PHOTO (DataCreate, Path, Activated)\n" +
		"VALUES (?, ?, 1)\n" +
		"ON CONFLICT(Path) DO UPDATE SET\n" +
		"DataCreate = excluded.DataCreate,\n" +
		"Activated = excluded.Activated;"
	rollbackImage = "UPDATE PHOTO SET Activated=1 WHERE Path like ?"
	setKeyImage   = "UPDATE USER SET Photo=?"
)

func (r Request) changePhoto() (rs.NonAddResponce, error) {
	answer := rs.NonAddResponce{Status: rs.Responce{}.InternalError("внутренняя ошибка при смене фото")}
	exists, err := Utility.MediaExists(r.Image)
	if err != nil {
		return answer, fmt.Errorf("ошибка поиска файла фотографии. Подробности: %w", err)
	}
	if !exists {
		answer.Status = answer.Status.BadRequest("указанного файла не существует")
		return answer, nil
	}
	err = r.user.openDB()
	if err != nil {
		return answer, fmt.Errorf("сбой соединения с БД пользователя. Подробности: %w", err)
	}
	defer r.user.closeDB()
	err = r.latestPhoto()
	if err != nil {
		return answer, fmt.Errorf("сбой получения текущего фото с БД пользователя. Подробности: %w", err)
	}
	_, err = r.user.sql.instance.Exec(disableImage)
	if err != nil {
		errs := r.rollbackPhoto()
		return answer, fmt.Errorf("сбой активации нового фото пользователя. Подробности: %s. Ошибка отката: %w", err, errs)
	}
	_, err = r.user.sql.instance.Exec(updateImage, Utility.GetTime(), r.Image)
	if err != nil {
		errs := r.rollbackPhoto()
		return answer, fmt.Errorf("сбой активации нового фото пользователя. Подробности: %s. Ошибка отката: %w", err, errs)
	}
	_, err = r.user.sql.instance.Exec(setKeyImage, r.Image)
	if err != nil {
		errs := r.rollbackPhoto()
		return answer, fmt.Errorf("ошибка привязки нового фото. Подробности: %s. Ошибки при откате: %w", err, errs)

	}
	answer.Status = answer.Status.OK("Фото профиля изменено")
	return answer, nil
}

func (r Request) rollbackPhoto() error {
	if r.oldImage == "" {
		return nil
	}
	_, err := r.user.sql.instance.Exec(rollbackImage, r.oldImage)
	return err
}

func (r *Request) latestPhoto() error {
	rows := r.user.sql.instance.QueryRow(lastImage)
	err := rows.Scan(&r.oldImage)
	if err != sql.ErrNoRows && err != nil {
		return fmt.Errorf("ошибка получения последнего фото. Подробности: %w", err)
	}
	return nil
}

func (r *Request) Do() (rs.Resp, error) {
	result := rs.MultiResponce{Status: rs.Responce{}.BadRequest("запрос не распознан. Уточните запрос, и повторите")}
	var err error
	r.user, err = r.getUser()
	if err != nil {
		result.Status = result.Status.BadRequest("Некорректное значение токена")
		return result, fmt.Errorf("сбой получения пользователя из сессии. Подробности: %w", err)
	}
	if r.isChangePass() {
		answer, err := r.changePass()
		if err != nil {
			return answer, err
		}
		result.Status = result.Status.OK("Запрос распознан и выполнен")
		result.Messages = append(result.Messages, answer.Status.Description)
	}
	if r.isChangeEmail() {

		result.Status.Code = 501
		result.Status.Description = "в текущей версии не реализована эта функция"
	}
	if r.isChangeLogin() {
		answer, err := r.changeNickname()
		if err != nil {
			return answer, err
		}
		result.Status = result.Status.OK("Запрос распознан и выполнен")
		result.Messages = append(result.Messages, answer.Status.Description)
	}
	if r.isChangePhoto() {
		answer, err := r.changePhoto()
		if err != nil {
			return answer, err
		}
		result.Status = result.Status.OK("Запрос распознан и выполнен")
		result.Messages = append(result.Messages, answer.Status.Description)
	}
	return result, nil
}

func (r Request) isChangePass() bool {
	return (r.NewPass != "" && r.OldPass != "")
}

func (r Request) isChangeEmail() bool {
	return r.Email != ""
}

func (r Request) isChangeLogin() bool {
	return r.Login != ""
}

func (r Request) isChangePhoto() bool {
	return r.Image != ""
}

func (r Request) getUser() (User, error) {
	return UserBySessia(r.Sessia)
}

func timePass(lastchange time.Time) bool {
	return lastchange.Compare(time.Now().Add(time.Duration(-1*day))) != -1
}

func monthPass(lastchange time.Time) bool {
	return lastchange.Compare(time.Now().Add(time.Duration(-1*month))) != -1
}
