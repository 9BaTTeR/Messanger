package goEmail

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	emailDB = "CREATE TABLE MAP (\n" +
		"	Email TEXT PRIMARY KEY,\n" +
		"	ID NUMBER);"
	setEmail = "INSERT INTO MAP\n" +
		"VALUES (?,?)"
	getID = "SELECT ID FROM MAP \n" +
		"Where Email like ?"
)

func CreateDB() error {
	err := createFileDB()
	if err != nil {
		return err
	}
	err = genereDB()
	if err != nil {
		return err
	}
	return nil
}

func createFileDB() error {
	f, err := Utility.CreateFile(Configs.EmailDB())
	if err != nil {
		return fmt.Errorf("ошибка создание UserDir/email.db map[email-id] базы данных. Подробности: %w", err)
	}
	defer f.Close()
	return nil
}

func genereDB() error {
	sql, err := openDB()
	if err != nil {
		return err
	}
	defer sql.Close()
	_, err = sql.Exec(emailDB)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса. Подробности: %w", err)
	}
	return sql.Close()
}

func openDB() (*sql.DB, error) {
	p := Configs.EmailDB()
	exists, err := Utility.Exists(p)
	if err != nil {
		return nil, fmt.Errorf("невозможно проверить наличие БД. Подробности: %w", err)
	}
	if !exists {
		err := CreateDB()
		if err != nil {
			return nil, err
		}
	}
	sql, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, fmt.Errorf("ошибка соединения с базой данных. Подробности: %w", err)
	}
	return sql, nil
}
