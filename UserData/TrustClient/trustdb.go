package trustclient

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const (
	trustDB = "CREATE TABLE CERTS (\n" +
		"  	HashCert TEXT NOT NULL PRIMARY KEY,\n" +
		"	Certificate BLOB ,	\n" +
		"  	NameClient TEXT NOT NULL,\n" +
		"  	MinVersion INTEGER NOT NULL,\n" +
		"  	MaxVersion INTEGER NOT NULL,\n" +
		"  	Active INTEGER NOT NULL);"
)

func createDB() error {
	err := createFileDB()
	if err != nil {
		return fmt.Errorf("сбой формирование файла базы данных сертификатов. Подробности: %w", err)
	}
	instance, err := openDB()
	instance.SetConnMaxLifetime(lifetime)
	if err != nil {
		return fmt.Errorf("подключение к БД провалилось. Подробности: %w", err)
	}
	err = createCertDB(instance)
	if err != nil {
		return fmt.Errorf("сбой инициализации sessia.db. Подробности: %w", err)
	}
	return instance.Close()
}

func createFileDB() error {
	file, err := Utility.CreateFile(Configs.CertsDB())
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", Configs.CertsDB())
}

func createCertDB(s *sql.DB) error {
	_, err := s.Exec(trustDB)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу сертификатов. Подробности: %w", err)
	}
	return nil
}
