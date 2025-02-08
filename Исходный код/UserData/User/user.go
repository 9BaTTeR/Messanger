package user

import (
	"ServerApp/UserData/Sessia"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	updateHistory = "UPDATE Email\n" +
		"SET Code = ?\n" +
		"WHERE Email like ?"
	getRowEmailByPK = "Select * From EMAIL"
	getRowUserByPK  = "Select Photo,Login,Password FROM USER ORDER BY rowid DESC LIMIT 1; "
)

func (u User) addPassword() error {
	querry := "INSERT INTO PASSWORD VALUES (?, ?, ?);"
	_, err := u.sql.instance.Exec(querry, Utility.GetTime(), u.Password.Pass, Utility.BoolToUInt(true))
	if err != nil {
		return err
	}
	return nil
}

func UserBySessia(t Sessia.Sessia) (User, error) {
	u := User{}
	uid, err := t.ParseID()
	if err != nil {
		return u, nil
	}
	u.PK = uid.IDConversion()
	return u, nil
}

func (u User) addLogin() error {
	querry := "INSERT INTO LOGIN VALUES (?, ?, ?);"
	_, err := u.sql.instance.Exec(querry, Utility.GetTime(), u.Login.Nickname, Utility.BoolToUInt(true))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) FieldByPk() error {
	if u.PK == 0 {
		return fmt.Errorf("PK не определен")
	}
	err := u.fieldUserByPK()
	if err != nil {
		return err
	}
	err = u.fieldEmaildByPk()
	if err != nil {
		return err
	}
	return nil
}

func FromSessia(s Sessia.Sessia) (u User, err error) {
	uid, err := s.ParseID()
	if err != nil {
		err = fmt.Errorf("ошибка чтения ID из сессии. Подробности: %w", err)
		return
	}
	u.PK = uid.IDConversion()
	err = u.FieldByPk()
	if err != nil {
		err = fmt.Errorf("ошибка заполнение user через ID. Подробности: %w", err)
		return
	}
	return
}

func (u *User) fieldEmaildByPk() error {
	instance, err := u.UserSqlInstance()
	if err != nil {
		return err
	}
	defer instance.Close()
	row, err := instance.Query(getRowEmailByPK)
	if err != nil {
		return err
	}
	defer row.Close()
	if !row.Next() {
		return fmt.Errorf("записей нет")
	}
	var activated sql.NullInt64
	err = row.Scan(&u.Email.DataCreate, &u.Email.Email, &activated, &u.Email.Code)
	if err != nil {
		return err
	}
	u.Email.Activated = Utility.UintToBool(uint64(activated.Int64))
	return err
}

func (u *User) fieldUserByPK() error {
	sql, err := u.UserSqlInstance()
	if err != nil {
		return err
	}
	defer sql.Close()
	row, err := sql.Query(getRowUserByPK)
	if err != nil {
		return err
	}
	defer row.Close()
	if !row.Next() {
		return fmt.Errorf("записей нет")
	}
	err = row.Scan(&u.Photo.Path, &u.Login.Nickname, &u.Password.Pass)
	if err != nil {
		return fmt.Errorf("ошибка сканирование сведений о пользователе")
	}
	return err
}

func (u User) IsVerified() bool {
	if u.restrained {
		return true
	}
	return false
}

func (u User) addEmail() error {
	querry := "INSERT INTO EMAIL VALUES (?, ?, ?, ?);"
	_, err := u.sql.instance.Exec(querry, Utility.GetTime(), u.Email.Email, Utility.BoolToUInt(false), u.code)
	if err != nil {
		return err
	}
	return nil
}

func (u User) addUserDB() error {
	querry := "INSERT INTO USER(UserId, Password, Login, Email) VALUES (?, ?, ?, ?);"
	_, err := u.sql.instance.Exec(querry, strconv.FormatUint(u.PK, 10), u.Password.Pass, u.Login.Nickname, u.Email.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u User) UserSql() (*sql.DB, error) {
	if u.sql.instance == nil {
		return nil, fmt.Errorf("соеденение с БД не открыто")
	}
	return u.sql.instance, nil
}
