package search

import (
	"ServerApp/Configs"
	rs "ServerApp/Responces"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	nameMF  = "Dialogs"
)

const (
	searchMOnD    = "SELECT OMessages.OMKey FROM DMessages, OMessages WHERE DMessages.DMKey like OMessages.DMKey AND content like ? LIMIT ?  OFFSET ?;"
	searchDialogs = "select Hash from HistoryDialogsUser Order By UpdateAt Desc Limit ? Offset ?;"
	searchMsg     = "Select * FROM DMessages, OMessages WHERE DMessages.DMKey like OMessages.DMKey AND instr(content, ?) COLLATE NOCASE LIMIT 1;"
)

func (r Request) Do() (rs.Resp, error) {
	switch r.Operation {
	case searchindialogs:
		return r.doSearchOnDialog()
	case searchAny:
		return r.searchAny()
	default:
		return rs.NonAddResponce{
			Status: rs.Responce{}.BadRequest(""),
		}, nil
	}
}

func (r Request) searchAny() (rs.SearchDialogsContainsResponce, error) {
	answer := rs.SearchDialogsContainsResponce{
		Status: rs.Responce{}.InternalError("Внутренняя ошибка при поиске"),
	}
	path, err := r.Sessia.UserFolder()
	if err != nil {
		return answer, fmt.Errorf("ошибка чтения папки пользователя . Подробности: %w", err)
	}
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{path, Configs.NameHDB()}))
	if err != nil {
		return answer, fmt.Errorf("сбой соединения с БД историей диалогов пользователя . Подробности: %w", err)
	}

	defer sql.Close()
	rows, err := sql.Query(searchDialogs, r.TakeDialogs, r.SkipDialogs)
	if err != nil {
		return answer, fmt.Errorf("сбой получения записей последних актуальных диалогов. Подробности: %s ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var hash string
		rows.Scan(&hash)
		result, err := dialogContains(path, r.Search, hash)
		if err != nil {
			return answer, fmt.Errorf("ошибка при перекрёстном поиске. Подробности: %w", err)
		}
		if result {
			answer.Dialogs = append(answer.Dialogs, hash)
		}
	}
	answer.Status = answer.Status.OK("Успешно!")
	return answer, nil
}

func dialogContains(userpath string, desired string, hash string) (bool, error) {
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{userpath, "Dialogs", hash, Configs.NameMDB()}))
	if err != nil {
		return false, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(searchMsg, desired)
	if err != nil {
		return false, fmt.Errorf("ошибка получения записей в БД. Подробности: %s ", err)
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (r Request) doSearchOnDialog() (rs.SearchMessageResponce, error) {
	answer := rs.SearchMessageResponce{
		Status: rs.Responce{}.InternalError("Внутренняя ошибка при поиске"),
	}
	uf, err := r.Sessia.UserFolder()
	if err != nil {
		return answer, fmt.Errorf("ошибка получения диалога при поиске. Подробности: %w", err)
	}
	path := Utility.Combine([]string{uf, nameMF, r.Dialogs, Configs.NameMDB()})
	sql, err := sqldrv.Open("sqlite3", path)
	if err != nil {
		return answer, fmt.Errorf("ошибка подключения к БД диалога при поиске. Подробности: %w", err)
	}
	stmt, err := sql.Prepare(searchMOnD)
	if err != nil {
		return answer, fmt.Errorf("ошибка препарирования БД диалога при поиске. Подробности: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query("%"+r.Search+"%", r.Take, r.Skip)
	if err != nil {
		return answer, fmt.Errorf("ошибка получения списка сообщений. Подробности: %w", err)
	}

	for rows.Next() {
		message := ""
		err := rows.Scan(&message)
		if err != nil {
			return answer, fmt.Errorf("ошибка парсинга значения. Подробности: %w", err)
		}
		answer.Messages = append(answer.Messages, message)
	}
	answer.Status = answer.Status.OK("Список сообщений сформирован.")
	return answer, nil
}