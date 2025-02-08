package goEmail

import (
	responces "ServerApp/Responces"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	internal     = "Внутренняя ошибка"               //500
	unauthorized = "Пользователь не зарегестрирован" //401
	exists_email = "Почта уже зарегестрировнаа"      //400
)

func (ac *Request) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &ac)
	if err != nil {
		return err
	}
	return nil
}
func (ac Request) Compose() ([]byte, error) {
	data, err := json.Marshal(ac)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e Email) LinkEmail() (responces.Responce, error) {
	exists, err := e.EmailExists()
	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка проверки наличия email. Подробности: %w", err)
	}
	if exists {
		return responces.Responce{Code: 400, Description: exists_email}, nil
	}
	sql, err := openDB()
	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка соединения с БД. Подробности %w", err)
	}
	defer sql.Close()
	if e.ID == (Email{}).ID {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ID не определён")
	}
	_, err = sql.Exec(setEmail, e.Email, e.ID)

	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка записи [email-id] в БД. Подробности %w", err)
	}
	return responces.Responce{Code: 200, Description: "OK"}, sql.Close()
}

func (e *Email) FindID() (responces.Responce, error) {
	sql, err := openDB()
	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка соединения с БД. Подробности %w", err)
	}
	defer sql.Close()
	if e.ID != (Email{}).ID {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ID уже определён")
	}
	result, err := sql.Query(getID, e.Email)
	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка чтения [email-id] в БД. Подробности %w", err)
	}
	defer result.Close()
	if !(result.Next()) {
		return responces.Responce{Code: 401, Description: internal}, fmt.Errorf("пользователь не зарегистрирован")
	}
	err = result.Scan(&e.ID)
	if err != nil {
		return responces.Responce{Code: 500, Description: internal}, fmt.Errorf("ошибка парсинга значений ID из email.db. Подробности: %w", err)
	}
	return responces.Responce{}, sql.Close()
}

func (e Email) EmailExists() (bool, error) {
	sql, err := openDB()
	if err != nil {
		return false, fmt.Errorf("ошибка соединения с БД. Подробности %w", err)
	}
	defer sql.Close()
	result, err := sql.Query(getID, e.Email)
	if err != nil {
		return false, fmt.Errorf("ошибка чтения [email-id] в БД. Подробности %w", err)
	}
	defer result.Close()
	if !(result.Next()) {
		return false, sql.Close()
	}
	return true, sql.Close()
}

