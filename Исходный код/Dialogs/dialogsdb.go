package dialogs

import (
	"ServerApp/Configs"
	"database/sql"
	"fmt"
)

const (
	DbDialog = "CREATE TABLE Dialog (" +
		"Name TEXT PRIMARY KEY," +
		"Privacy TEXT NOT NULL," +
		"CreatedAt TEXT NOT NULL," +
		"Photo TEXT);"
	DbMembers = "CREATE TABLE Members (" +
		"IdUser Integer PRIMARY KEY," +
		"Role TEXT NOT NULL," +
		"DateJoin TEXT NOT NULL," +
		"Notice TEXT NOT NULL," +
		"FOREIGN KEY(Role) REFERENCES Roles(NameRole));"
	DbRoles = "CREATE TABLE Roles (" +
		"NameRole TEXT PRIMARY KEY," +
		"Role TEXT);"
	AddDbRoles = "INSERT INTO Roles " +
		"VALUES" +
		"('Пользователь', 'Стандартные возможности: создавать ссылку-приглашение, добавлять новых участников, писать в чат ')," +
		"('Администратор', 'Стандартные возможности: имеет все права пользователя, а также может выгонять участников, удалять сообщения')," +
		"('Создатель', 'Имеет все права');"
	DbBans = "CREATE TABLE BANS (" +
		"IdUser Integer PRIMARY KEY," +
		"Reason TEXT," +
		"Duration TEXT DEFAULT 0 NOT NULL," +
		"BanBy Integer NOT NULL);"
	DbLinks = "CREATE TABLE LINKS (" +
		"Link TEXT PRIMARY KEY," +
		"IdUser Integer NOT NULL," +
		"Duration TEXT DEFAULT 0," +
		"Count Integer DEFAULT 0 NOT NULL);"
	DbNameSettings = "CREATE TABLE NameSettings (" +
		"Name TEXT PRIMARY KEY);"
	AddNameSettings = "INSERT INTO NameSettings " +
		"VALUES" +
		"('Пользователи могут создавать ссылку приглашение')," +
		"('Пользователи могут добавлять новых пользователей')," +
		"('Пользователи могут выгонять участников');"
	DbValueSettings = "CREATE TABLE ValueSettings (" +
		"Value TEXT PRIMARY KEY);"
	AddValueSettings = "INSERT INTO ValueSettings " +
		"VALUES" +
		"('Да')," +
		"('Нет');"
	DbSettings = "CREATE TABLE Settings (" +
		"IdSetting Integer PRIMARY KEY AUTOINCREMENT," +
		"Name TEXT NOT NULL," +
		"Value TEXT NOT NULL," +
		"Role TEXT NOT NULL," +
		"FOREIGN KEY(Name) REFERENCES NameSettings(Name)," +
		"FOREIGN KEY(Value) REFERENCES ValueSettings(Value)," +
		"FOREIGN KEY(Role) REFERENCES Roles(NameRole));"
	DbPinned = "CREATE TABLE Pinned (" +
		"Date TEXT PRIMARY KEY," +
		"IdUser Integer NOT NULL," +
		"OMKey TEXT NOT NULL);"
)
const (
	internalerror = "Внутренняя ошибка"
)

func (d Dialog) createBaseDialog() error {
	//path := Configs.Path{}.FolderDefDB().Path()
	path := Configs.SettingsDefDDB()
	instance, err := openDB(path)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД. Подробности: %w", err)
	}
	defer instance.Close()
	err = createDialogDB(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbMembers(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbRoles(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbBans(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbLinks(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbNameSettings(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbValueSettings(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbSettings(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}
	err = createDbPinned(*instance)
	if err != nil {
		return fmt.Errorf("ошибка создания базового dialog.db Подробности: %w", err)
	}

	return nil
}

func openDB(path string) (*sql.DB, error) {
	sql, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия бд. Подробности %w", err)
	}
	return sql, nil
}

func createDialogDB(s sql.DB) error {
	_, err := s.Exec(DbDialog)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы диалога. Подробности: %w", err)
	}
	return nil
}
func createDbMembers(s sql.DB) error {
	_, err := s.Exec(DbMembers)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы участников. Подробности: %w", err)
	}
	return nil
}
func createDbRoles(s sql.DB) error {
	_, err := s.Exec(DbRoles)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы роли. Подробности: %w", err)
	}
	_, err = s.Exec(AddDbRoles)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы роли. Подробности: %w", err)
	}
	return nil
}
func createDbBans(s sql.DB) error {
	_, err := s.Exec(DbBans)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы банов. Подробности: %w", err)
	}
	return nil
}
func createDbLinks(s sql.DB) error {
	_, err := s.Exec(DbLinks)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы ссылок приглашений. Подробности: %w", err)
	}
	return nil
}
func createDbNameSettings(s sql.DB) error {
	_, err := s.Exec(DbNameSettings)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы названий настроек. Подробности: %w", err)
	}
	_, err = s.Exec(AddNameSettings)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы названий настроек. Подробности: %w", err)
	}
	return nil
}
func createDbValueSettings(s sql.DB) error {
	_, err := s.Exec(DbValueSettings)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы значений настроек. Подробности: %w", err)
	}
	_, err = s.Exec(AddValueSettings)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы значений настроек. Подробности: %w", err)
	}
	return nil
}
func createDbSettings(s sql.DB) error {
	_, err := s.Exec(DbSettings)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы настроек. Подробности: %w", err)
	}
	return nil
}
func createDbPinned(s sql.DB) error {
	_, err := s.Exec(DbPinned)
	if err != nil {
		return fmt.Errorf("ошибка создание таблицы закрепленного сообщения. Подробности: %w", err)
	}
	return nil
}
