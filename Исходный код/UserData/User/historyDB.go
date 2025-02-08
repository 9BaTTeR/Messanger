package user

import (
	"ServerApp/Configs"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const (
	DbHistoryDialogsUser = "CREATE TABLE HistoryDialogsUser (" +
		"Hash TEXT PRIMARY KEY," +
		"UpdateAt TEXT NOT NULL);"
	fillDialogToHistory = "INSERT OR REPLACE INTO HistoryDialogsUser" +
		" VALUES (?,?);"
)

func CreateHistoryDialogDB() error {
	err := createBaseFilesHistoryDB()
	if err != nil {
		return fmt.Errorf("ошибка создания файла базового dialogs.db  Подробности: %w", err)
	}
	instance, err := sql.Open("sqlite3", Configs.HistoryDefDB())
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	defer instance.Close()
	err = createDbHistoryDialog(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания historyDialog.db Подробности: %w", err)
	}


	return nil
}
func createDbHistoryDialog(s sql.DB) error {
	_, err := s.Exec(DbHistoryDialogsUser)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы истории диалогов. Подробности: %w", err)
	}
	return nil
}

func createBaseFilesHistoryDB() error {
	path := Configs.HistoryDefDB()
	f, err := Utility.CreateFile(path)
	if err != nil {
		return fmt.Errorf("сбой генерации файла базы данных истории диалогов. Подробности: %w", err)
	}
	f.Close()
	return nil
}

func AddHistoryDialog(IdUsers []uint64, NameDialog string) error {
	//d.IdUsers = append(IdUsers, idUserCreator)
	for i := range IdUsers {
		uid := uid.ConvertionID(IdUsers[i])
		path := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameHDB()})
		err := checkHistoryDialog(path)
		if err != nil {
			return err
		}
		instance, err := sql.Open("sqlite3", path)
		if err != nil {
			return err
		}
		_, err = instance.Exec(fillDialogToHistory, NameDialog, Utility.GetTime())
		if err != nil {
			return fmt.Errorf("ошибка наполнения таблицы историй диалогов. Подробности: %w", err)
		}
		err = instance.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// here
func checkHistoryDialog(path string) error {
	pathdefault := Configs.HistoryDefDB()
	exist, err := Utility.Exists(pathdefault)
	if !exist {
		err = CreateHistoryDialogDB()
	}
	if err != nil {
		return err
	}
	exist, err = Utility.Exists(path)
	if !exist {
		err = Utility.CopyFile(pathdefault, path)
	}
	if err != nil {
		return err
	}
	return nil
}
