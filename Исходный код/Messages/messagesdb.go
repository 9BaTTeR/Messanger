package message

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const (
	DbDMessages = "CREATE TABLE DMessages (" +
		"DMKey TEXT PRIMARY KEY," +
		"IdUser Integer NOT NULL," +
		"Date TEXT NOT NULL," +
		"Content TEXT," +
		"ForwardedKey TEXT," +
		"Read TEXT NOT NULL DEFAULT false," +
		"Important TEXT NOT NULL DEFAULT false," +
		"DeletedAt TEXT," +
		"UpdateAt TEXT)" +
		"WITHOUT ROWID;"
	DbOMessages = "CREATE TABLE OMessages (" +
		"OMKey TEXT PRIMARY KEY," +
		"Date TEXT NOT NULL," +
		"DMKey TEXT NOT NULL, " +
		"FOREIGN KEY(DMKey) REFERENCES DMessages(DMKey))" +
		"WITHOUT ROWID;"
	DbMedia = "CREATE TABLE Media (" +
		"IdMedia Integer PRIMARY KEY AUTOINCREMENT, " +
		"Hash TEXT NOT NULL, " +
		"OrderBy Integer NOT NULL, " +
		"DMKey TEXT, " +
		"FOREIGN KEY(DMKey) REFERENCES DMessages(DMKey));"
)

func CopyMessagesDB(old string, new string) error {
	exists, err := Utility.Exists(old)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("не удается прочитать файл при копировании")
	}
	err = Utility.CopyFile(old, new)
	if err != nil {
		return err
	}
	return nil
}

func CreateMessages() error {
	path := Configs.MessagesDefDB()

	instance, err := openDB(path)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	defer instance.Close()
	err = createDbDMessages(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового messages.db Подробности: %w", err)
	}
	err = createDbOMessages(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового messages.db Подробности: %w", err)
	}
	err = createDbMedia(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового messages.db Подробности: %w", err)
	}
	return nil
}

func createDbDMessages(s sql.DB) error {
	_, err := s.Exec(DbDMessages)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы свалки сообщений. Подробности: %w", err)
	}
	return nil
}
func createDbOMessages(s sql.DB) error {
	_, err := s.Exec(DbOMessages)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы оригинальных сообщений. Подробности: %w", err)
	}
	return nil
}
func createDbMedia(s sql.DB) error {
	_, err := s.Exec(DbMedia)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы медиа. Подробности: %w", err)
	}
	return nil
}

func openDB(path string) (*sql.DB, error) {
	sql, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return sql, nil
}

