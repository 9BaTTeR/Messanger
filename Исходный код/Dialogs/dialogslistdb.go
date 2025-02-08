package dialogs

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const (
	DbAllDialogs = "CREATE TABLE ALLDIALOGS (" +
		"Hash TEXT PRIMARY KEY," +
		"Privacy TEXT NOT NULL);"
)

func CreateAllDialogsDB() error {
	//path := Configs.Path{}.FolderDialogs().Path()
	path := Configs.ChatsDB()
	_, err := Utility.CreateFile(path)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	instance, err := openDB(path)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	defer instance.Close()
	err = createDbAllDialog(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового allDialog.db Подробности: %w", err)
	}
	return nil
}
func createDbAllDialog(s sql.DB) error {
	_, err := s.Exec(DbAllDialogs)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы allDialogs. Подробности: %w", err)
	}
	return nil
}
