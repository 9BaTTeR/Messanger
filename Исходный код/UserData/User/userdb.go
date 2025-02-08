package user

import (
	// dialog "ServerApp/Dialogs"
	"ServerApp/Configs"
	"ServerApp/UserData/Sessia"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

const (
	PasswordDB = "CREATE TABLE PASSWORD (\n" +
		"  DataCreate TEXT NOT NULL,\n" +
		"  Password TEXT PRIMARY KEY NOT NULL,\n" +
		"  Activated INTEGER NOT NULL);\n"
	LoginDB = "CREATE TABLE LOGIN (\n" +
		"  DataCreate TEXT NOT NULL,\n" +
		"  Login TEXT PRIMARY KEY NOT NULL,\n" +
		"  Activated INTEGER NOT NULL);"
	EmailDB = "CREATE TABLE EMAIL (\n" +
		"  DataCreate TEXT NOT NULL,\n" +
		"  Email TEXT PRIMARY KEY NOT NULL,\n" +
		"  Activated INTEGER NOT NULL,\n" +
		"  Code TEXT);"
	PhotoDB = "CREATE TABLE PHOTO (\n" +
		"  DataCreate TEXT ,\n" +
		"  Path TEXT PRIMARY KEY,\n" +
		"  Activated INTEGER);"
	UserDB = "CREATE TABLE USER (\n" +
		"  UserId INTEGER PRIMARY KEY NOT NULL,\n" +
		"  Password TEXT NOT NULL,\n" +
		"  Login TEXT NOT NULL,\n" +
		"  Email TEXT NOT NULL,\n" +
		"  Photo TEXT NOT NULL DEFAULT '',\n" +
		"  FOREIGN KEY(Password) REFERENCES PASSWORD(Password),\n" +
		"  FOREIGN KEY(Login) REFERENCES LOGIN(Login),\n" +
		"  FOREIGN KEY(Email) REFERENCES EMAIL(Email),\n" +
		"  FOREIGN KEY(Photo) REFERENCES PHOTO(Path));"
)

func (u *User) CopyDB() error {
	err := checkDefaultDbs()
	if err != nil {
		return err
	}
	pathBase := Configs.DefaultsDB()
	uId := uid.ConvertionID(uint64(u.PK))
	pathUserFolder := uId.PathDbUserId()
	err = Utility.CopyFile(path.Join(pathBase, Configs.NameUDB()), path.Join(pathUserFolder, Configs.NameUDB()))
	if err != nil {
		return fmt.Errorf("ошибка копирования User.db. Подробности: %w", err)
	}
	err = Utility.CopyFile(path.Join(pathBase, Configs.NameSDB()), path.Join(pathUserFolder, Configs.NameSDB()))
	if err != nil {
		return fmt.Errorf("ошибка копирования sessia.db. Подробности: %w", err)
	}
	err = Utility.CopyFile(path.Join(pathBase, Configs.NameFDB()), path.Join(pathUserFolder, Configs.NameFDB()))
	if err != nil {
		return fmt.Errorf("ошибка копирования friends.db. Подробности: %w", err)
	}
	err = Utility.CopyFile(path.Join(pathBase, Configs.NameHDB()), path.Join(pathUserFolder, Configs.NameHDB()))
	if err != nil {
		return fmt.Errorf("ошибка копирования dialogs.db. Подробности: %w", err)
	}
	err = Utility.CopyFile(path.Join(pathBase, Configs.NameBLDB()), path.Join(pathUserFolder, Configs.NameBLDB()))
	if err != nil {
		return fmt.Errorf("ошибка копирования blacklist.db. Подробности: %w", err)
	}
	u.sql.instance, err = sql.Open("sqlite3", pathUserFolder+Configs.NameUDB())
	if err != nil {
		return fmt.Errorf("ошибка открытия sql.instanse. Подробности: %w", err)
	}
	return nil
}

func checkDefaultDbs() error {
	exists, err := Utility.Exists(Configs.UserDefDB())
	if err != nil {
		return fmt.Errorf("ошибка чтения/создания базовой user.db. Подробности: %w", err)
	}
	if !exists {
		err = createBaseDB()
		if err != nil {
			return err
		}
	}
	exists, err = Utility.Exists(Configs.SessiaDefDB())
	if err != nil {
		return fmt.Errorf("ошибка чтения/создания базовой sessia.db. Подробности: %w", err)
	}
	if !exists {
		err = Sessia.CreateBaseDB()
		if err != nil {
			return err
		}
	}
	exists, err = Utility.Exists(Configs.FriendsDefDB())
	if err != nil {
		return fmt.Errorf("ошибка чтения/создания базовой friends.db. Подробности: %w", err)
	}
	if !exists {
		err = CreateFriendsBaseDB()
		if err != nil {
			return err
		}
	}
	exists, err = Utility.Exists(Configs.HistoryDefDB())
	if err != nil {
		return fmt.Errorf("ошибка чтения/создания базовой historyDialog.db. Подробности: %w", err)
	}
	if !exists {
		err = CreateHistoryDialogDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func createBaseDB() error {
	instance, err := openDB()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	defer instance.Close()
	err = createFileDB()
	if err != nil {
		return err
	}
	err = createLoginDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создание базовой user.db Подробности: %w", err)
	}
	err = createPasswordDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создание базовой user.db Подробности: %w", err)
	}
	err = createEmailDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создание базовой user.db Подробности: %w", err)
	}
	err = createPhotoDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создание базовой user.db Подробности: %w", err)
	}
	err = createUserDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создание базовой user.db Подробности: %w", err)
	}
	return nil
}

func createFileDB() error {

	file, err := Utility.CreateFile(Configs.UserDefDB())
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", Configs.UserDefDB())
}

func (u User) UserSqlInstance() (*sql.DB, error) {
	err := u.openDB()
	if err != nil {
		return nil, err
	}
	sql := u.sql.instance
	if err != nil {
		return nil, err
	}
	return sql, nil
}

func (u *User) openDB() error {
	if u.sql.instance != nil {
		return nil
	}
	uid := uid.ConvertionID(u.PK)
	path := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameUDB()})
	var err error
	u.sql.instance, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) closeDB() error {
	if u.sql.instance == nil {
		return nil
	}
	err := u.sql.instance.Close()
	u.sql.instance = nil
	return err
}

func createLoginDB(s sql.DB) error {
	_, err := s.Exec(LoginDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы логина. Подробности: %w", err)
	}
	return nil
}
func createPasswordDB(s sql.DB) error {
	_, err := s.Exec(PasswordDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы паролей. Подробности: %w", err)
	}
	return nil
}
func createEmailDB(s sql.DB) error {
	_, err := s.Exec(EmailDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы электронной почты. Подробности: %w", err)
	}
	return nil
}
func createPhotoDB(s sql.DB) error {
	_, err := s.Exec(PhotoDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы фотографий. Подробности: %w", err)
	}
	return nil
}
func createUserDB(s sql.DB) error {
	_, err := s.Exec(UserDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы пользователя. Подробности: %w", err)
	}
	return nil
}
