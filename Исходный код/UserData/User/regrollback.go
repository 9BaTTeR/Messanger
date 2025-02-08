package user

import (
	cf "ServerApp/Configs"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"
)

const (
	removeEmail = "Delete from MAP where EMAIL like ?"
)

func (reg User) Rollback() {
	fmt.Printf("\nОткат емейла, логи: %v\n", rollbackEmail(reg.Email.Email))
	fmt.Printf("\nОткат папки, логи: %v\n", rollbackFolder(reg.PK))
}

func rollbackEmail(email string) error {
	pathDB := cf.EmailDB()
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return fmt.Errorf("ошибка отката email. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(removeEmail, email)
	if err != nil {
		return fmt.Errorf("ошибка удаления email из БД. Подробности: %w", err)
	}
	return err
}

func rollbackFolder(id uint64) error {
	uid := userID.ConvertionID(id)
	return Utility.RemoveFolder(uid.PathDbUserId())
}
