package Sessia

import (
	"ServerApp/Utility"
	"database/sql"
	"fmt"
)

const ()

func (h *History) Field(s *sql.DB) error {
	row, err := s.Query(getRowHistory, h.DateUse)
	if err != nil {
		return err
	}
	err = row.Scan(&h)
	return err
}

func (h *History) Flush(s *sql.DB) error {
	_, err := s.Exec(insertRowHistory, h.DateUse, h.Ip)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса FLUSH для HISTORY. Подробности: %w", err)
	}
	return nil
}

func (h *History) Update(s *sql.DB) error {
	_, err := s.Exec(updateHistory, Utility.GetTime(), h.Ip)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса UPDATE для HISTORY. Подробности: %w", err)
	}
	return nil
}
