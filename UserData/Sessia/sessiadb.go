package Sessia

import (
	"ServerApp/Configs"
	_ "ServerApp/Responces"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const (
	deviceDB = "CREATE TABLE DEVICE(\n" +
		"  DEVICEID TEXT PRIMARY KEY NOT NULL,\n" +
		"  OS TEXT NOT NULL,\n" +
		"  MAC TEXT NOT NULL,\n" +
		"  HOSTNAME TEXT NOT NULL);"
	insertRowDevice = "INSERT INTO DEVICE\n" +
		" VALUES(?,?,?,?);"
	getRowDevice = "Select * From DEVICE Where DEVICEID like ?"
)
const (
	historyDB = "CREATE TABLE HISTORY (\n" +
		"  DATEUSE TEXT PRIMARY KEY NOT NULL,\n" +
		"  IP TEXT NOT NULL);"
	getRowHistory    = "Select * From HISTORY Where DATEUSE like ?"
	insertRowHistory = "INSERT INTO HISTORY\n" +
		" VALUES(?,?);"
)
const (
	sessiaDB = "CREATE TABLE SESSIA (\n" +
		"  TOCKEN TEXT PRIMARY KEY NOT NULL,\n" +
		"  AVAILABLE NUMBER NOT NULL,\n" +
		"  DEVICEID TEXT NOT NULL,\n" +
		"  HISTORY TEXT NOT NULL,\n" +
		"  FOREIGN KEY(DEVICEID) REFERENCES DEVICE(DEVICEID)\n" +
		"  FOREIGN KEY(HISTORY) REFERENCES HISTORY(HISTORY));"
	insertRowSessia = "INSERT INTO SESSIA\n" +
		" VALUES(?,?,?,?);"
	updateHistory = "UPDATE SESSIA\n" +
		"SET HISTORY = ?\n" +
		"WHERE TOCKEN like ?;"
	updateDevice = "UPDATE SESSIA\n" +
		"SET DEVICEID = ?\n" +
		"WHERE TOCKEN like ?;"
)

const (
	certificateDB = "\n" +
		"\n" +
		"\n" +
		"\n"
)

const (
	internalerror = "Внутренняя ошибка"
)

func CreateBaseDB() error {
	instance, err := openDB()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	err = createDeviceDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базовой user.db Подробности: %w", err)
	}
	err = createSessiaDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базовой user.db Подробности: %w", err)
	}
	err = createHistoryDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базовой user.db Подробности: %w", err)
	}
	err = instance.Close()
	if err != nil {
		return fmt.Errorf("ошибка закрытия соединения с БД. Подробности: %w", err)
	}
	return nil
}

func createFileDB() error {
	file, err := Utility.CreateFile(Configs.SessiaDefDB())
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (t *Sessia) openDB() error {
	if t.sql != nil {
		return nil
	}
	uid, err := t.ParseID()
	if err != nil {
		return fmt.Errorf("неверный токен. Подробности: %w", err)
	}
	path := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameSDB()})
	t.sql, err = sql.Open("sqlite3", path)
	return err
}
func (t *Sessia) closeDB() error {
	if t.sql == nil {
		return nil
	}
	err := t.sql.Close()
	t.sql = nil
	return err

}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", Configs.SessiaDefDB())
}

func createDeviceDB(s sql.DB) error {
	_, err := s.Exec(deviceDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы устройств. Подробности: %w", err)
	}
	return nil
}
func createSessiaDB(s sql.DB) error {
	_, err := s.Exec(sessiaDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы сессий. Подробности: %w", err)
	}
	return nil
}
func createHistoryDB(s sql.DB) error {
	_, err := s.Exec(historyDB)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы электронной почты. Подробности: %w", err)
	}
	return nil
}
