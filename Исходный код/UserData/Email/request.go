package goEmail

import (
	"ServerApp/Configs"
	rs "ServerApp/Responces"
	gencode "ServerApp/UserData/GenCode"
	"ServerApp/UserData/MailSender"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"
	"path"
	"strings"
)

func (r Request) Do(ip4 string) (rs.NonAddResponce, error) {
	if strings.ReplaceAll(r.Code, " ", "") != "" {
		return r.acceptCode()
	}
	if strings.ReplaceAll(r.Email, " ", "") != "" {
		return r.newCode(ip4)
	}
	return rs.NonAddResponce{Status: rs.Responce{}.BadRequest("некорректный JSON")}, nil
}

const (
	changeEmail = "UPDATE EMAIL SET Code = ? Where Email like ?"
)

func (r Request) newCode(ip4 string) (rs.NonAddResponce, error) {
	err := r.openDB()
	answer := rs.NonAddResponce{
		Status: rs.Responce{}.InternalError("внутренний сбой при смене кода"),
	}
	defer r.closeDB()
	if err != nil {
		return answer, fmt.Errorf("сбой смены кода. Подробности: %w", err)
	}
	active, err := r.alreadyActive()
	if err != nil {
		return answer, fmt.Errorf("сбой смены кода. Подробности: %w", err)
	}
	if active {
		answer.Status = answer.Status.BadRequest("email уже активирован")
		return answer, nil
	}
	active, err = r.mailExists()
	if err != nil {
		return answer, fmt.Errorf("сбой проверки привязки почты. Подробности: %w", err)
	}
	if !active {
		answer.Status = answer.Status.BadRequest("Указанная почта не привязана к вашему аккаунту.")
		return answer, nil
	}
	code := gencode.GenCode()
	_, err = r.sql.Exec(changeEmail, code, r.Email)
	if err != nil {
		return answer, fmt.Errorf("сбой обновления записи кода в БД. Подробности: %w", err)
	}
	nick, err := r.nick()
	if err != nil {
		return answer, fmt.Errorf("сбой смены кода. Подробности: %w", err)
	}
	err = MailSender.MailAuthCode(r.Email, nick, string(code), MailSender.AboutUser(ip4, Utility.GetTime(), r.Sessia))
	if err != nil {
		return answer, fmt.Errorf("сбой отсылки нового кода на почту пользователя. Подробности: %w", err)
	}
	answer.Status = answer.Status.OK("Код сменён и отправлен пользователю на почту")
	return answer, nil
}

func (r *Request) openDB() error {
	a, err := r.Sessia.UserFolder()
	if err != nil {
		return fmt.Errorf("сбой получения пользователя из сессии. Подробности: %w", err)
	}
	paths := path.Join(a, Configs.NameUDB())
	sql, err := sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой содинения с БД. Подробности: %w", err)
	}
	r.sql = sql
	return nil
}

func (r *Request) closeDB() error {
	defer func() { r.sql = nil }()
	return r.sql.Close()
}

const (
	existsMail = "SELECT Email FROM USER WHERE Email like ?"
)

func (r Request) mailExists() (bool, error) {
	if r.sql == nil {
		return false, fmt.Errorf("sql соединение закрыто")
	}
	var email string
	rows := r.sql.QueryRow(existsMail, r.Email)
	err := rows.Scan(&email)
	if err == sqldrv.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("сбой получения активации. Подрнобности: %w", err)
	}
	return true, nil
}

const (
	nicknamesql = "select Login From USER LIMIT 1;"
)

func (r Request) nick() (string, error) {
	if r.sql == nil {
		return "", fmt.Errorf("sql соединение закрыто")
	}
	var nickname string
	rows := r.sql.QueryRow(nicknamesql)
	err := rows.Scan(&nickname)
	if err != nil {
		return "", fmt.Errorf("сбой получения имени пользователя. Подрнобности: %w", err)
	}
	return nickname, nil
}

const (
	active = "select Activated from EMAIL where Email like ? Limit 1"
)

func (r Request) alreadyActive() (bool, error) {
	if r.sql == nil {
		return false, fmt.Errorf("sql соединение закрыто")
	}
	var temp uint8
	rows := r.sql.QueryRow(active, r.Email)
	err := rows.Scan(&temp)
	if err == sqldrv.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("сбой получения активации. Подрнобности: %w", err)
	}
	result := Utility.UintToBool(uint64(temp))
	return result, nil
}

const (
	activateMail = "UPDATE EMAIL SET Activated = 1 WHERE Email IN (Select Email from USER LIMIT 1) AND CODE LIKE ?"
)

func (r Request) acceptCode() (rs.NonAddResponce, error) {
	if len(strings.ReplaceAll(r.Code, " ", "")) != 6 {
		return rs.NonAddResponce{
			Status: rs.Responce{}.BadRequest("Код некорректный"),
		}, nil
	}
	answer := rs.NonAddResponce{
		Status: rs.Responce{}.BadRequest("внутренний сбой при подтверждении почты"),
	}
	err := r.openDB()
	if err != nil {
		return answer, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer r.closeDB()
	opeate, err := r.sql.Exec(activateMail, r.Code)
	if err != nil {
		return answer, fmt.Errorf("сбой активации почты в БД. Подробности: %w", err)
	}
	count, err := opeate.RowsAffected()
	if err != nil {
		return answer, fmt.Errorf("сбой получения кол-ва обновленных записей в БД. Подробности: %w", err)
	}
	if count == 0 {
		answer.Status = answer.Status.BadRequest("неверный код")
		return answer, nil
	}
	answer.Status = answer.Status.OK("почта подтверждена")
	return answer, nil

}
