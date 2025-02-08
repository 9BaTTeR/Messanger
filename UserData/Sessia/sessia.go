package Sessia

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	getRowSessiaByTocken = "Select AVAILABLE,DEVICEID,HISTORY From SESSIA Where Tocken like ?"
	rowSessiaByTocken    = "Select AVAILABLE from Sessia Where Tocken like ?"
	getLastRowSessia     = "Select * From SESSIA ORDER BY rowid DESC LIMIT 1"
	setRowSessia         = "UPDATE SESSIA\n" +
		"SET AVAILABLE = ?," +
		"DEVICEID = ?," +
		"HISTORY = ?;"
)

const (
	excludedRequest = "Certificate UpdateCertificate Auth Registration"
	disabletocken   = "UPDATE SESSIA SET AVAILABLE = 0 WHERE TOCKEN like ?"
)

func DisableTocken(content []byte) (responces.NonAddResponce, error) {
	resp := responces.NonAddResponce{}
	resp.Status = responces.Responce{}.InternalError("Внутренняя ошибка")
	s := EmptySessia{}
	err := json.Unmarshal(content, &s)
	if err != nil {
		return resp, fmt.Errorf("ошибка unmarshal json. Подробности: %v", err)
	}
	err = s.Sessia.disable()
	if err != nil {
		return resp, fmt.Errorf("ошибка отключения токена. Подробности: %v", err)
	}
	resp.Status = resp.Status.OK("Сессия отключена")
	return resp, nil
}

func (s Sessia) disable() error {
	err := s.openDB()
	if err != nil {
		return fmt.Errorf("ошибка соединения с БД. Подробности: %w", err)
	}
	defer s.closeDB()
	result, err := s.sql.Exec(disabletocken, s.Tocken)
	if err != nil {
		return fmt.Errorf("ошибка отключения токена в БД. Подробности: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества отключенных токенов в БД. Подробности: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("ни одна сессия не была отключена. Выполненный запрос: %v", strings.Replace(disabletocken, "?", s.Tocken, 1))
	}
	return nil
}

func ValidateTocken(content []byte, pathUrl string) (bool, responces.ValidateTockenResponce, error) {
	resp := responces.ValidateTockenResponce{}
	resp.Status = responces.Responce{}.InternalError("Внутренняя ошибка")
	excludes := strings.Split(excludedRequest, " ")
	if slices.Contains(excludes, pathUrl) {
		return true, responces.ValidateTockenResponce{}, nil
	}
	s := EmptySessia{}
	err := json.Unmarshal(content, &s)
	if err != nil {
		return false, resp, fmt.Errorf("ошибка unmarshal json. Подробности: %v", err)
	}
	valide, err := s.Sessia.validate()
	if err != nil {
		return false, resp, fmt.Errorf("ошибка чтения валидации. Подробности: %v", err)
	}
	if !valide {
		resp.Status = responces.Responce{}.BadRequest("Токен не валиден")
		return false, resp, fmt.Errorf("токен невалиден")
	}
	return true, resp, nil
}

// Пхах, sessiaByTocken сломан. Надо исправить
func (s *Sessia) FieldSessiaFromDB() error {
	if (s.Tocken == Sessia{}.Tocken) {
		return fmt.Errorf("токен пустой")
	}
	err := s.openDB()
	if err != nil {
		return fmt.Errorf("Сбой соединения с БД. Подробности: %w", err)
	}
	err = s.sessiaByTocken()
	if err != nil {
		return err
	}
	s.closeDB()
	return nil
}

func (s Sessia) PathDB() (string, error) {
	path, err := s.UserFolder()
	return Utility.Combine([]string{path, Configs.NameSDB()}), err
}

func (s *Sessia) Flush() error {
	err := s.openDB()
	if err != nil {
		return err
	}
	exists, err := s.sql.Query(getRowSessiaByTocken, s.Tocken)
	if err != nil {
		return err
	}
	defer exists.Close()
	if !exists.Next() {
		return s.flush()
	}
	_, err = s.sql.Exec(setRowSessia, Utility.BoolToUInt(s.Available), s.DeviceID.DeviceID, s.History.DateUse)
	if err != nil {
		return err
	}
	return s.closeDB()
}

func (s Sessia) ConnectDB() (*sqldrv.DB, error) {
	err := s.openDB()
	return s.sql, err
}

func (s *Sessia) flush() error {
	_, err := s.sql.Exec(insertRowSessia, s.Tocken, Utility.BoolToUInt(s.Available), s.DeviceID.DeviceID, s.History.DateUse)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы SESSIA. Подробности: %w", err)
	}
	return nil
}

func (s *Sessia) sessiaByTocken() error {
	row, err := s.sql.Query(getRowSessiaByTocken, s.Tocken)
	if err != nil {
		return err
	}
	defer row.Close()
	if !row.Next() {
		return fmt.Errorf("записей нет")
	}
	var tempavailable uint8
	err = row.Scan(&tempavailable, &s.DeviceID.DeviceID, &s.History.DateUse)
	if err != nil {
		return fmt.Errorf("cбой получения записей. Подробности: %w", err)
	}
	return err
}

func (t Sessia) ParseID() (userID.Storage, error) {
	tempstringIDs := strings.Split(t.Tocken, "?")
	if len(tempstringIDs) < 2 {
		return userID.Storage{}, fmt.Errorf("некорректная строка. Подробности: %+v", tempstringIDs)
	}
	stringIDs := tempstringIDs[1]
	stringID := strings.Replace(stringIDs, "?", "", 1)
	ID, err := strconv.ParseUint(stringID, 10, 32)
	fmt.Println(stringID)
	if err != nil {
		return userID.Storage{}, fmt.Errorf("некорректный токен, подробности: %w", err)
	}
	userID := userID.Storage{}
	userID.Set(ID)
	return userID, nil
}

func (t Sessia) UserFolder() (string, error) {
	paths := Configs.UserDir()
	userID := userID.Storage{}
	userID, err := t.ParseID()
	if err != nil {
		return "", err
	}
	pathUserFolder := path.Join(paths, strconv.FormatUint(uint64(userID.MDir), 10), strconv.FormatUint(uint64(userID.SDir), 10))
	return pathUserFolder, nil
}

func (t Sessia) validate() (bool, error) {
	err := t.openDB()
	if err != nil {
		return false, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer t.closeDB()
	rows, err := t.sql.Query(rowSessiaByTocken, t.Tocken)
	if err != nil {
		return false, fmt.Errorf("ошибка чтения сессии по токену. Подробности: %w", err)
	}
	if !rows.Next() {
		return false, nil
	}
	var availabletemp uint8
	err = rows.Scan(&availabletemp)
	if err != nil {
		return false, fmt.Errorf("ошибка парсинга bool значения. Подробности: %w", err)
	}
	defer rows.Close()
	t.Available = Utility.UintToBool(uint64(availabletemp))
	return t.Available, nil
}

func (t Sessia) FlushAll() error {
	err := t.openDB()
	if err != nil {
		return fmt.Errorf("ошибка соединения с БД сессии при FlushAll. Подробности: %w", err)
	}
	defer t.closeDB()
	err = t.History.Flush(t.sql)
	if err != nil {
		return err
	}
	err = t.openDB()
	if err != nil {
		return fmt.Errorf("ошибка соединения с БД сессии при FlushAll(2). Подробности: %w", err)
	}
	err = t.DeviceID.Flush(t.sql)
	if err != nil {
		return err
	}
	defer t.closeDB()
	err = t.Flush()
	if err != nil {
		return err
	}
	return nil
}

const (
	isActive = "SELECT Activated FROM Email WHERE Activated LIKE 1 LIMIT 1"
)

func (s Sessia) IsWithActivateEmail() (bool, error) {
	pathusr, err := s.UserFolder()
	if err != nil {
		return false, fmt.Errorf("сбой получения папки пользователя. Подробности: %w", err)
	}
	path := path.Join(pathusr, Configs.NameUDB())
	sql, err := sqldrv.Open("sqlite3", path)
	if err != nil {
		return false, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	rows := sql.QueryRow(isActive)
	var temp uint8
	err = rows.Scan(&temp)
	if err == sqldrv.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("сбой получения записей. Подробности: %w", err)
	}
	result := Utility.UintToBool(uint64(temp))
	return result, nil
}
