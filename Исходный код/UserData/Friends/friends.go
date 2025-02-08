package friends

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	Us "ServerApp/UserData/User"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"
)

const (
	isFriendSql = "Select id,description,dateadd FROM FRIENDS Where ID like ?"
	loadMore    = "SELECT PHOTO,LOGIN FROM USER ORDER BY rowid DESC LIMIT 1;"
)

func (f Friend) Convert() responces.Friend {
	r := responces.Friend{
		ID:          f.ID,
		Description: f.Description,
		DateAdd:     f.DateAdd,
		Photo:       f.Photo,
		Nickname:    f.Nickname,
	}
	return r

}

func description(from uint64, of uint64) (string, error) {
	uid := userID.ConvertionID(from)
	pathDB := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameFDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return "", fmt.Errorf("сбой соединения с БД при получении описания для %v у пользователя %v. Подробности: %s", of, from, err)
	}
	defer sql.Close()
	rows, err := sql.Query(getDescription, of)
	if err != nil {
		return "", fmt.Errorf("ошибка получения описания для %v у пользователя %v. Подробности: %s", of, from, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return "", nil
	}
	var result string
	err = rows.Scan(&result)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга полученного значения описания для %v у пользователя %v. Подробности: %s", of, from, err)
	}
	return result, nil
}

// User - кого опрашиваем.
// ID - кого загружаем.
// Необходимость в тестировании: Максимальная.
func LoadFriend(user Us.User, ID uint64) (Friend, error) {
	f := Friend{}
	uid := userID.ConvertionID(user.PK)
	if user.PK == 0 {
		var err error
		uid, err = user.Sessia.ParseID()
		if err != nil {
			return f, err
		}
	}
	path := uid.PathDbUserId()
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{path, Configs.NameFDB()}))
	if err != nil {
		return f, fmt.Errorf("ошибка загрузки БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(isFriendSql, ID)
	if err != nil {
		return f, fmt.Errorf("ошибка чтения записей из БД. Подробности: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return f, nil
	}
	err = rows.Scan(&f.ID, &f.Description, &f.DateAdd)
	if err != nil {
		return f, fmt.Errorf("ошибка скана значений в структуру. Подробности: %w", err)
	}
	err = f.loadData(ID)
	if err != nil {
		return f, fmt.Errorf("ошибка загрузки доп сведений. Подробности: %w", err)
	}
	return f, nil
}

func LoadUserData(ID uint64) (Friend, error) {
	f := Friend{}
	err := f.loadData(ID)
	if err != nil {
		return f, nil
	}
	return f, nil
}

func (user *Friend) loadData(ID uint64) error {
	uid := userID.ConvertionID(ID)
	path := uid.PathDbUserId()
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{path, Configs.NameUDB()}))
	if err != nil {
		return nil
	}
	defer sql.Close()
	rows, err := sql.Query(loadMore)
	if err != nil {
		return fmt.Errorf("ошибка получения дополнительной информации. Подробности: %w", err)
	}
	if !rows.Next() {
		return nil
	}
	err = rows.Scan(&user.Photo, &user.Nickname)
	if err != nil {
		return fmt.Errorf("ошибка выгрузки дополнительной информации. Подробности: %w", err)
	}
	defer rows.Close()
	return nil
}

func (f Friend) IsEmpty() bool {
	if (f.ID == Friend{}.ID) {
		return true
	}
	if (f.DateAdd == Friend{}.DateAdd) {
		return true
	}
	return false
}
