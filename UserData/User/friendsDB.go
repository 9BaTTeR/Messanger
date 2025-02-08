package user

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	friendsDB = "CREATE TABLE FRIENDS(\n" +
		"	ID INTEGER PRIMARY KEY,\n" +
		"	DESCRIPTION TEXT,\n" +
		"	DATEADD TEXT) WITHOUT ROWID;"
	requestDB = "CREATE TABLE REQUEST(\n" +
		"	DATE TEXT NOT NULL, \n" +
		"  	ID INTEGER NOT NULL,\n" +
		"  	TYPEREQUEST TEXT,\n" +
		"  	DESCRIPTION TEXT,\n" +
		"	COMING INTEGER,\n" +
		"  	PRIMARY KEY(DATE,ID)\n" +
		")WITHOUT ROWID;"
	blacklistDB = "CREATE TABLE BLACKLIST(\n" +
		"ID INTEGER PRIMARY KEY,\n" +
		"DATE TEXT\n" +
		") WITHOUT ROWID;"
)

func CreateFriendsBaseDB() error {
	err := createBaseFilesFriendDB()
	if err != nil {
		return fmt.Errorf("ошибка инициализации базовой БД. Подробности: %s", err)
	}
	sql, err := sqldrv.Open("sqlite3", Configs.FriendsDefDB())
	if err != nil {
		return fmt.Errorf("ошибка подключения к базовой БД. Подробности: %s", err)
	}
	defer sql.Close()
	_, err = sql.Exec(friendsDB)
	if err != nil {
		return fmt.Errorf("ошибка генерации сущности friends в базовой БД. Подробности: %s", err)
	}
	_, err = sql.Exec(requestDB)
	if err != nil {
		return fmt.Errorf("ошибка генерации сущности request в базовой БД. Подробности: %s", err)
	}
	err = createBaseFilesBlackListDB()
	if err != nil {
		return fmt.Errorf("ошибка инициализации базовой БД. Подробности: %s", err)
	}
	sql, err = sqldrv.Open("sqlite3", Configs.BlackListDefDB())
	if err != nil {
		return fmt.Errorf("ошибка подключения к базовой БД. Подробности: %s", err)
	}
	defer sql.Close()
	_, err = sql.Exec(blacklistDB)
	if err != nil {
		return fmt.Errorf("ошибка генерации сущности blacklist в базовой БД. Подробности: %s", err)
	}

	return nil
}

func createBaseFilesFriendDB() error {

	f, err := Utility.CreateFile(Configs.FriendsDefDB())
	if err != nil {
		return fmt.Errorf("сбой генерации файла базы данных друзей. Подробности: %s", err)
	}
	f.Close()
	return nil
}

func createBaseFilesBlackListDB() error {

	f, err := Utility.CreateFile(Configs.BlackListDefDB())
	if err != nil {
		return fmt.Errorf("сбой генерации файла базы данных друзей. Подробности: %s", err)
	}
	f.Close()
	return nil
}
