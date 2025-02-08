package friends

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"
)

const (
	rowBlackList    = "Select * FROM BLACKLIST WHERE ID like ?;"
	insertBlackList = "INSERT INTO BLACKLIST(ID,DATE)\n" +
		"VALUES(?,?); "
	deleteBlackList = "DELETE FROM BLACKLIST WHERE ID like ?;"
	selectBlackList = "SELECT ID,DATE FROM BLACKLIST ORDER BY DATE DESC LIMIT ? OFFSET ?"
)

func (r Request) alreadyBlock() (bool, error) {
	if r.sql == nil {
		return false, fmt.Errorf("нет соединения с БД в методе alreadyBlock")
	}
	rows, err := r.sql.Query(rowBlackList, r.Id)
	if err != nil {
		return false, fmt.Errorf("сбой проверки наличия существующей блокировки. Подробности: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// Первый bool - пользователь из реквеста заблокировал пользователя из r.id. (x1)
// Второй bool - пользователь из r.id заблокировал пользователя из реквеста. (x2)
func (r Request) AnyInBlockList() (bool, bool, error) {
	id1, err := r.Sessia.ParseID()
	if err != nil {
		return false, false, fmt.Errorf("сбой чтения сведений о блокировках. Подробности: %w", err)
	}
	block1, err := checkBlock(id1.IDConversion(), r.Id)
	if err != nil {
		return false, false, fmt.Errorf("сбой чтения сведений о блокировках реквестера. Подробности: %w", err)
	}
	block2, err := checkBlock(r.Id, id1.IDConversion())
	if err != nil {
		return false, false, fmt.Errorf("сбой чтения сведений о блокировке не реквестера. Подробности: %w", err)
	}
	return block1, block2, nil
}

func checkBlock(from uint64, check uint64) (bool, error) {
	uid := userID.ConvertionID(from)
	ufolder := uid.PathDbUserId()
	fmt.Println(Utility.Combine([]string{ufolder, Configs.NameBLDB()}))
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{ufolder, Configs.NameBLDB()}))
	if err != nil {
		return false, fmt.Errorf("сбой соединения с БД блокировок. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(rowBlackList, check)
	if err != nil {
		return false, fmt.Errorf("сбой выполнения запроса БД. Подробности: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (r Request) listBlocks() (responces.ListBlocksResponce, error) {
	responce := responces.ListBlocksResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка при работче с чёрным списком."),
	}
	path, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пути пользователя из сессии. Подробности: %w", err)
	}
	pathDB := Utility.Combine([]string{path, Configs.NameBLDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД пользователя. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(selectBlackList, r.Take, r.Skip)
	if err != nil {
		return responce, fmt.Errorf("ошибка получения записей из БД пользователя. Подробности: %w", err)
	}
	defer rows.Close()
	users := []responces.Friend{}
	for rows.Next() {
		user := Friend{}
		err := rows.Scan(&user.ID, &user.DateAdd)
		if err != nil {
			return responce, fmt.Errorf("ошибка загрузки данных из БД. Подробности: %w", err)
		}
		err = user.loadData(user.ID)
		if err != nil {
			return responce, fmt.Errorf("ошибка загрузки дополнительных данных из БД. Подробности: %w", err)
		}
		users = append(users, user.Convert())
	}
	responce.Status = responce.Status.OK("Список сформирован")
	responce.Blocks = users
	return responce, nil
}

func (r Request) addBlock() (responces.BlackListResponce, error) {
	responce := responces.BlackListResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка при работче с чёрным списком."),
	}
	path, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пути пользователя из сессии. Подробности: %w", err)
	}
	pathDB := Utility.Combine([]string{path, Configs.NameBLDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД пользователя. Подробности: %w", err)
	}
	defer sql.Close()
	r.sql = sql
	exists, err := r.alreadyBlock()
	if err != nil {
		return responce, fmt.Errorf("ошибка проверки наличия добавляемой блокировки у пользователя. Подробности: %w", err)
	}
	if exists {
		responce.Status = responce.Status.BadRequest("пользователь уже заблокирован")
		return responce, nil
	}
	_, err = sql.Exec(insertBlackList, r.Id, Utility.GetTime())
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления пользователя в чёрный список. Подробности: %w", err)
	}
	if sql.Close() == nil {
		_, err = r.remove()
	} else {
		return responce, fmt.Errorf("сбой завершения соединения с БД при добавления пользователя в ЧС. Подробности: %w", err)
	}
	if err != nil {
		return responce, fmt.Errorf("ошибка удаление пользователей из друзей после добавления в ЧС. Подробности: %w", err)

	}
	responce.Status = responce.Status.OK("Пользователь заблокирован")
	return responce, nil
}

func (r Request) delBlock() (responces.BlackListResponce, error) {
	responce := responces.BlackListResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка при работче с чёрным списком."),
	}
	path, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пути пользователя из сессии. Подробности: %w", err)
	}
	pathDB := Utility.Combine([]string{path, Configs.NameBLDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД пользователя. Подробности: %w", err)
	}
	defer sql.Close()
	r.sql = sql
	exists, err := r.alreadyBlock()
	if err != nil {
		return responce, fmt.Errorf("ошибка проверки наличия добавляемой блокировки у пользователя. Подробности: %w", err)
	}
	if !exists {
		responce.Status = responce.Status.BadRequest("пользователь не заблокирован")
		return responce, nil
	}
	_, err = sql.Exec(deleteBlackList, r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка удаления пользователя из чёрного списка. Подробности: %w", err)
	}
	responce.Status = responce.Status.OK("Пользователь разблокирован")
	return responce, nil
}
