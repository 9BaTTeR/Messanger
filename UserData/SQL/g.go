package SQL

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLInterface interface {
	CreateDB() error
	CopyDB() error
	ExecQuerry(querry string) error
}

type SqlInstance struct {
	Instance sql.DB
}

// Получаем экземпляр sqlite куда будем производить запись.
func OpenDB(path string) (SqlInstance, error) {
	sql, err := sql.Open("sqlite3", path)
	if err != nil {
		return SqlInstance{}, err
	}
	sql.Begin()
	return SqlInstance{Instance: *sql}, nil
}

// Закрываем соединение.
func (sq *SqlInstance) Close() error {
	err := sq.Instance.Close()
	if err != nil {
		return err
	}
	return nil
}
func (sq *SqlInstance) ExecQuerry(querry string, values[]string) error {
	_, err := sq.Instance.Exec(querry)
	if err != nil {
		return err
	}
	return nil
}

func (sq *SqlInstance) ResultQuerry(querry string, values[]string) (*sql.Rows, error) {
	result, err := sq.Instance.Query(querry)
	if err != nil {
		return nil, err
	}
	return result, nil
}
