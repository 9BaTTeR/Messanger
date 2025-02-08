package friends

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	user "ServerApp/UserData/User"
	userID "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	sqldrv "database/sql"
	"fmt"
	"path"
	"time"
)

const (
	incoming  uint8 = 0
	outcoming uint8 = 1
)

// FRIENDS TABLE
const (
	countFriends    = "select COUNT(*) From FRIENDS"
	updateFriends   = "UPDATE FRIENDS SET DESCRIPTION = ? Where ID like ?"
	insertFriends   = "INSERT INTO FRIENDS(ID,DESCRIPTION,DATEADD) VALUES(?,?,?)"
	removeFreinds   = "DELETE FROM FRIENDS WHERE ID = ?"
	isAlreadyFriend = "Select * FROM FRIENDS WHERE ID like ?"
	getDescFriend   = "SELECT DESCRIPTION FROM FRIENDS WHERE ID = ? AND DESCRIPTION not like '' AND DESCRIPTION IS NOT NULL LIMIT 1 "
	takeFriends     = "SELECT ID FROM FRIENDS LIMIT ? OFFSET ?"
)

// REQUEST TABLE
const (
	insertRequest      = "INSERT into REQUEST(date,id,typerequest,description,coming) VALUES(?,?,?,?,?)"
	isAlreadyRequested = "Select * FROM REQUEST WHERE ID like ? AND coming like ?"
	deleteRequest      = "delete from REQUEST WHERE ID = ?"
	getDescription     = "SELECT DESCRIPTION FROM REQUEST WHERE ID = ?  AND DESCRIPTION not like '' AND DESCRIPTION IS NOT NULL LIMIT 1 "
	getReqByComing     = "SELECT DATE,ID,DESCRIPTION FROM REQUEST WHERE COMING like ? LIMIT ? OFFSET ?"
	removeRequest      = "DELETE FROM REQUEST WHERE ID = ?"
)

// Метод добавления заявок в БД.
func (r Request) append() (responces.Resp, error) {
	responce := responces.AppendResponce{}
	responce.Status = responce.Status.InternalError("внутренняя ошибка")
	//blocks - проверка блокировок с обеих сторон
	blocks1, blocks2, err := r.AnyInBlockList()
	if err != nil {
		return responce, fmt.Errorf("сбой проверки блокировок пользователей. Подробности: %w", err)
	}
	if blocks1 {
		responce.Status = responce.Status.BadRequest("Пользователь находится в чёрном списке")
		return responce, nil
	}
	if blocks2 {
		responce.Status = responce.Status.BadRequest("Вы находитесь в чёрном списке")
		return responce, nil
	}
	//already - уже друзья.
	uid, err := r.Sessia.ParseID()
	if err != nil {
		return responce, fmt.Errorf("сбой получения ID из сессии. Подробности: %w", err)
	}
	already, err := r.isAlreadyFriend(uid.IDConversion(), r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления пользователя в друзья. Подробности: %w", err)
	}
	if already {
		responce.Status = responce.Status.BadRequest("пользователь уже в друзьях")
		return responce, nil
	}
	//cross - пересекающиеся запросы в друзья. already - заявка уже создана.
	cross, already, err := r.isAlreadyRequested(uid.IDConversion(), r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления пользователя в друзья. Подробности: %w", err)
	}
	if cross {
		uid, err := r.Sessia.ParseID()
		if err != nil {
			return responce, fmt.Errorf("неверный токен. Подробности: %w", err)
		}
		inID := uid.IDConversion()
		return r.acceptFriend(inID, r.Id)
	}
	if already {
		responce.Status = responce.Status.BadRequest("Заявка уже была отправлена.")
		return responce, nil
	}
	userDB1, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пользователя из сессии. Подробности: %s", err)
	}
	userDB1 = Utility.Combine([]string{userDB1, Configs.NameFDB()})
	//uid = userID.ConvertionID(r.Id)
	uid2 := userID.ConvertionID(r.Id)

	userDB2 := uid2.PathDbUserId()
	userDB2 = Utility.Combine([]string{userDB2, Configs.NameFDB()})

	sql, err := sqldrv.Open("sqlite3", userDB1)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД исходящего пользователя. Подробности: %s", err)
	}
	defer sql.Close()

	sql2, err := sqldrv.Open("sqlite3", userDB2)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД запрашиваемого пользователя. Подробности: %s", err)
	}
	defer sql2.Close()
	_, err = sql.Exec(insertRequest, r.Date.String(), r.Id, appends, r.Description, outcoming)
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления записи запроса в БД исходящего пользователя. Подробности: %w", err)
	}
	_, err = sql2.Exec(insertRequest, r.Date.String(), uid.IDConversion(), appends, "", incoming)
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления записи запроса в БД запрашиваемого пользователя. Подробности: %w", err)
	}
	responce.Status = responce.Status.OK("Заявка в друзья отправлена.")
	return responce, nil
}

func (r Request) isAlreadyFriend(first uint64, second uint64) (bool, error) {
	fpath := path.Join(userID.ConvertionID(first).PathDbUserId(), Configs.NameFDB())
	fsql, err := sqldrv.Open("sqlite3", fpath)
	if err != nil {
		return false, fmt.Errorf("сбой подключения к БД друзей пользователя %v. Подробности: %w", first, err)
	}
	defer fsql.Close()
	rows1, err := fsql.Query(isAlreadyFriend, second)
	if err != nil {
		return false, fmt.Errorf("невозможно прочитать список друзей пользователя. Подробности: %s", err)
	}
	defer rows1.Close()
	spath := path.Join(userID.ConvertionID(second).PathDbUserId(), Configs.NameFDB())
	ssql, err := sqldrv.Open("sqlite3", spath)
	if err != nil {
		return false, fmt.Errorf("сбой подключения к БД друзей пользователя %v. Подробности: %w", second, err)
	}
	defer ssql.Close()
	rows2, err := ssql.Query(isAlreadyFriend, first)
	if err != nil {
		return false, fmt.Errorf("невозможно прочитать список друзей пользователя. Подробности: %w", err)
	}
	defer rows2.Close()
	if rows1.Next() || rows2.Next() {
		return true, nil
	}
	return false, nil
}

func (r Request) isAlreadyRequested(first uint64, second uint64) (bool, bool, error) {
	fpath := path.Join(userID.ConvertionID(first).PathDbUserId(), Configs.NameFDB())
	fsql, err := sqldrv.Open("sqlite3", fpath)
	if err != nil {
		return false, false, fmt.Errorf("сбой подключения к БД друзей пользователя %v. Подробности: %w", first, err)
	}
	defer fsql.Close()
	spath := path.Join(userID.ConvertionID(second).PathDbUserId(), Configs.NameFDB())
	ssql, err := sqldrv.Open("sqlite3", spath)
	if err != nil {
		return false, false, fmt.Errorf("сбой подключения к БД друзей пользователя %v. Подробности: %w", second, err)
	}
	defer ssql.Close()
	rows1, err := fsql.Query(isAlreadyRequested, second, outcoming)
	if err != nil {
		return false, false, fmt.Errorf("сбой чтения записей из БД при проверке запроса в друзья у первого пользователя. Подробности: %w", err)
	}
	defer rows1.Close()
	if err != nil {
		return false, false, fmt.Errorf("сбой получения ID пользователя из сессии. Подробности: %w", err)
	}
	rows2, err := ssql.Query(isAlreadyRequested, first, outcoming)
	if err != nil {
		return false, false, fmt.Errorf("сбой чтения записей из БД при проверке запроса в друзья у второго пользователя. Подробности: %w", err)
	}
	defer rows2.Close()

	b1, b2 := rows1.Next(), rows2.Next()
	//Наличие пересекающейся заявки
	if b2 {
		return true, true, nil
	}
	//Наличие заявки
	if b1 {
		return false, true, nil
	}

	//1. bool - Пересекающиеся заявки
	//2. bool - Наличие заявки
	return false, false, nil
}

// Подробное описание
// InID - тот, кому запрос пришёл. В случае с запросом accept, это будет Request.Sessia.ParseID()
// OutID - тот, от кого запрос пришёл. В случае с запросом accept, это будет Request.ID
func (r Request) acceptFriend(inID uint64, outID uint64) (responces.AddFriendResponce, error) {

	responce := responces.AddFriendResponce{}
	responce.Status = responce.Status.InternalError("внутренняя ошибка")
	if inID == outID{
		return responce, fmt.Errorf("попытка принять запрос от самого себя. Подробности: %w", nil)
	}
	inDS, err := description(outID, inID)
	outDS := r.Description
	if err != nil {
		return responce, fmt.Errorf("ошибка получения описания друзей. Подробности: %w", err)
	}
	err = r.deleteRequest(inID, outID)
	if err != nil {
		errs := r.restoreRequest(inID, inDS, outID, r.Description)
		return responce, fmt.Errorf("ошибка добавления в друзья у inID. Подробности: %w. Лог отката: %w", err, errs)
	}
	err = r.addFriend(inID, outDS, outID, inDS)
	if err != nil {
		errs := r.restoreRequest(inID, inDS, outID, r.Description)
		return responce, fmt.Errorf("ошибка добавления в друзья у outID. Подробности: %w. Лог отката: %w", err, errs)
	}
	responce.Status = responce.Status.OK("Заявка в друзья принята")
	return responce, nil
}

func (r Request) addFriend(inID uint64, inDescription string, outID uint64, outDescription string) error {
	uid := userID.ConvertionID(inID)
	paths := path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err := sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД inDS. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(insertFriends, outID, inDescription, Utility.GetTime())
	if err != nil {
		return fmt.Errorf("сбой записи в друзья в БД inID. Подробности: %w", err)
	}
	uid = userID.ConvertionID(outID)
	paths = path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err = sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД outID. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(insertFriends, inID, outDescription, Utility.GetTime())
	if err != nil {
		return fmt.Errorf("сбой записи в друзья в БД outID. Подробности: %w", err)
	}
	return nil

}

// Подробное документирование
// InID - тот, кому изначальный запрос пришёл
// InDescription - описание того, от кого пришёл запрос, чаще всего ""
// outDS - тот, от кого изначальный запрос пришёл
// outDescription - описание того, к кому запрос был направлен. Возможен ""
// Вернёт nil если всё ок.
func (r Request) deleteRequest(inID uint64, outID uint64) error {
	uid := userID.ConvertionID(inID)
	paths := path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err := sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД inDS. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(deleteRequest, outID)
	if err != nil {
		return fmt.Errorf("сбой удаления запроса в БД inID для outID. Подробности: %w", err)
	}

	uid = userID.ConvertionID(outID)
	paths = path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err = sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД outID. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(deleteRequest, inID)
	if err != nil {
		return fmt.Errorf("сбой удаления запроса в БД outID для inID. Подробности: %w", err)
	}

	return nil
}

// Подробное документирование
// InID - тот, кому изначальный запрос пришёл
// InDescription - описание того, от кого пришёл запрос, чаще всего ""
// outDS - тот, от кого изначальный запрос пришёл
// outDescription - описание того, к кому запрос был направлен. Возможен ""
// Вернёт nil если всё ок.
func (r Request) restoreRequest(inID uint64, inDescription string, outID uint64, outDescription string) error {
	uid := userID.ConvertionID(inID)
	paths := path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err := sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД inDS. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(insertRequest, Utility.GetTime(), outID, appends, inDescription, incoming)
	if err != nil {
		return fmt.Errorf("сбой восстановления запроса для inID. Подробности: %w", err)
	}
	uid = userID.ConvertionID(outID)
	paths = path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err = sqldrv.Open("sqlite3", paths)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД outDS. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(insertRequest, Utility.GetTime(), inID, appends, outDescription, outcoming)
	if err != nil {
		return fmt.Errorf("сбой восстановления запроса для outID. Подробности: %w", err)
	}
	return nil
}

func (r Request) cancel() (responces.AddFriendResponce, error) {
	responce := responces.AddFriendResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка"),
	}

	userDB1, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пользователя из сессии. Подробности: %w", err)
	}
	userDB1 = Utility.Combine([]string{userDB1, Configs.NameFDB()})
	uid := userID.ConvertionID(r.Id)
	userDB2 := uid.PathDbUserId()
	userDB2 = Utility.Combine([]string{userDB2, Configs.NameFDB()})
	sqlin, err := sqldrv.Open("sqlite3", userDB1)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД in при попытке добавления в друзья. Подробности: %w", err)
	}
	defer sqlin.Close()
	sqlout, err := sqldrv.Open("sqlite3", userDB2)
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД out при попытке добавления в друзья. Подробности: %w", err)
	}
	defer sqlout.Close()
	_, err = sqlin.Exec(removeRequest, uid.IDConversion())
	if err != nil {
		return responce, fmt.Errorf("ошибка удаление запроса в друзья. Подробности: %w", err)
	}
	_, err = sqlout.Exec(removeRequest, r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка удаление запроса в друзья. Подробности: %w", err)
	}
	responce.Status = responce.Status.OK("Заявка в друзья отменена.")
	return responce, nil
}

func (r Request) accept() (responces.AddFriendResponce, error) {
	responce := responces.AddFriendResponce{Status: responces.Responce{}.InternalError("внутренняя ошибка")}
	uid, err := r.Sessia.ParseID()
	if err != nil {
		return responces.AddFriendResponce{
			Status: responces.Responce{}.BadRequest("некорректный токен"),
		}, fmt.Errorf("сбой получения ID. Подробности: %w", err)
	}
	already, err := r.isAlreadyFriend(uid.IDConversion(), r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка принятия пользователя в друзья. Подробности: %w", err)
	}
	if already {
		responce.Status = responce.Status.BadRequest("пользователь уже в друзьях")
		return responce, nil
	}
	can, err := r.canAccept(uid.IDConversion(), r.Id)
	if err != nil {
		return responce, fmt.Errorf("у пользователя нет заявки в друзья")
	}
	if !can {
		responce.Status = responce.Status.BadRequest("нет заявки в друзья от этого пользователя")
		return responce, nil
	}
	responce, err = r.acceptFriend(uid.IDConversion(), r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка добавления в друзья. Подробности: %w", err)
	}

	responce.Status = responce.Status.OK("Заявка в друзья принята")
	return responce, nil
}

const (
	isComing = "SELECT * FROM REQUEST WHERE ID LIKE ? AND  COMING LIKE ?"
)

func (r Request) canAccept(firstID uint64, secondID uint64) (bool, error) {
	uid := userID.ConvertionID(firstID)
	paths := path.Join(uid.PathDbUserId(), Configs.NameFDB())
	sql, err := sqldrv.Open("sqlite3", paths)
	if err != nil {
		return false, fmt.Errorf("ошибка соединения с БД при проверке наличия запроса в друзья. Подробности: %w", err) 
	}
	defer sql.Close()
	_, err = sql.Exec(isComing, secondID, incoming)
	if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса при проверке наличия запроса в друзья. Подробности: %w", err) 
	}
	return true, nil
}

func (r Request) update() (responces.UpdateFriend, error) {
	uf := responces.UpdateFriend{}
	uf.Status = responces.Responce{}.InternalError("внутрення ошибка при обновлении сведений о друге")
	pathUser, err := r.Sessia.UserFolder()
	if err != nil {
		return uf, fmt.Errorf("ошибка получение папки пользователя. Подробности: %w", err)
	}
	pathDB := Utility.Combine([]string{pathUser, Configs.NameFDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return uf, fmt.Errorf("сбой соединения с базой данных. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(updateFriends, r.Description, r.Id)
	if err != nil {
		return uf, fmt.Errorf("сбой обновления записи в базе данных. Подробности: %w", err)
	}
	uf.Status = responces.Responce{}.OK("Описание друга обновлено")
	return uf, nil
}

// User1 - удаляющий пользователь
// User2 - удаляемый пользователь
// Преобразование запросов User1 > incoming request, User2> outcoming request
func (r Request) remove() (responces.RemoveFriendResponce, error) {
	responce := responces.RemoveFriendResponce{}
	userDB1, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("ошибка чтения пользователя из сессии. Подробности: %w", err)
	}
	userDB1 = Utility.Combine([]string{userDB1, Configs.NameFDB()})
	uid := userID.ConvertionID(r.Id)
	userDB2 := uid.PathDbUserId()
	userDB2 = Utility.Combine([]string{userDB2, Configs.NameFDB()})

	sql1, err := sqldrv.Open("sqlite3", userDB1)
	if err != nil {
		return responce, fmt.Errorf("сбой соединения с БД user1. Подробности: %w", err)
	}
	defer sql1.Close()
	sql2, err := sqldrv.Open("sqlite3", userDB2)
	if err != nil {
		return responce, fmt.Errorf("сбой соединения с БД user2. Подробности: %w", err)
	}
	defer sql2.Close()
	row2, err := sql2.Query(getDescFriend, r.Id)
	if err != nil {
		return responce, fmt.Errorf("сбой получения описания из БД user2. Подробности: %w", err)
	}
	defer row2.Close()
	desc2 := ""
	if row2.Next() {
		row2.Scan(&desc2)
	}
	//Удаляем друзей из БД у друг друга сначала.
	_, err = sql1.Exec(removeFreinds, r.Id)
	if err != nil {
		return responce, fmt.Errorf("ошибка удаление из друзей. Подробности: %w", err)
	}
	_, err = sql2.Exec(removeFreinds, uid.IDConversion())
	if err != nil {
		return responce, fmt.Errorf("ошибка удаление из друзей. Подробности: %w", err)
	}
	//генерируем исходящую и входящую заявку в друзья.
	req := Request{
		Operation:   appends,
		Id:          r.Id,
		Coming:      outcoming,
		Description: desc2,
		Date:        time.Now(),
	}
	_, err = req.append()
	if err != nil {
		return responce, err
	}
	responce.Status = responce.Status.OK("Пользователь удалён из друзей.")

	return responce, nil
}

func (r *Request) getFriend() (responces.GetFriends, error) {
	tk := responces.GetFriends{
		Status: responces.Responce{}.InternalError("ошибка при получении сведений о друзей"),
	}
	if r.Take == 0 || r.Take >= 50 {
		tk.Status = responces.Responce{}.BadRequest("некорректное значение аргумента Take")
		return tk, fmt.Errorf("некорректное значение аргумента Take")
	}

	pathUser, err := r.Sessia.UserFolder()
	if err != nil {
		return tk, fmt.Errorf("ошибка чтения списка друзей. Подробности: %w", err)
	}
	pathDB := Utility.Combine([]string{pathUser, Configs.NameFDB()})
	sql, err := sqldrv.Open("sqlite3", pathDB)
	if err != nil {
		return tk, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(takeFriends, r.Take, r.Skip)
	if err != nil {
		return tk, fmt.Errorf("сбой запроса записей в БД. Подробности: %w", err)
	}
	defer rows.Close()
	Ids := []uint64{}
	for rows.Next() {
		var id uint64
		err := rows.Scan(&id)
		if err != nil {
			return tk, fmt.Errorf("ошибка чтения ID друзей. Подробности: %w", err)
		}
		Ids = append(Ids, id)
	}

	user, err := user.FromSessia(r.Sessia)
	if err != nil {
		return tk, fmt.Errorf("ошибка получение пользователя из сессии. Подробности: %w", err)
	}
	for _, id := range Ids {
		f, err := LoadFriend(user, id)
		if err != nil {
			return tk, fmt.Errorf("ошибка загрузки сведений о друге. Подробности: %w", err)
		}
		tk.TakeFriend = append(tk.TakeFriend, f.Convert())
	}
	tk.Status = responces.Responce{}.OK("ОК")
	return tk, nil
}

func (r Request) getComing(typecoming uint8) (responces.ComingResponce, error) {
	comings := []responces.Friend{}
	responce := responces.ComingResponce{
		Status: responces.Responce{}.InternalError("внутренняя ошибка"),
	}
	userdb, err := r.Sessia.UserFolder()
	if err != nil {
		return responce, fmt.Errorf("")
	}
	sql, err := sqldrv.Open("sqlite3", Utility.Combine([]string{userdb, Configs.NameFDB()}))
	if err != nil {
		return responce, fmt.Errorf("ошибка соединения с БД запросов. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(getReqByComing, typecoming, r.Take, r.Skip)
	if err != nil {
		return responce, fmt.Errorf("ошибка получения запросов. Подробности: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		coming := Friend{}
		rows.Scan(&coming.DateAdd, &coming.ID, &coming.Description)
		err = coming.loadData(coming.ID)
		if err != nil {
			return responce, fmt.Errorf("ошибка чтения выгрузки записей запросов. Подробности: %w", err)
		}
		comings = append(comings, coming.Convert())
	}
	responce = responces.ComingResponce{
		Status:  responce.Status.OK("Список сформирован."),
		Comings: comings,
	}
	return responce, nil
}

func (r Request) Do() (responces.Resp, error) {
	var err error
	var answer responces.Resp
	switch r.Operation {
	case accept:
		answer, err = r.accept()
	case appends:
		answer, err = r.append()
	case cancel:
		answer, err = r.cancel()
	case remove:
		answer, err = r.remove()
	case update:
		answer, err = r.update()
	case getFriend:
		answer, err = r.getFriend()
	case getIncoming:
		answer, err = r.getComing(incoming)
	case getOutcoming:
		answer, err = r.getComing(outcoming)
	case addBlackList:
		answer, err = r.addBlock()
	case removeBlackList:
		answer, err = r.delBlock()
	case getListBlackList:
		answer, err = r.listBlocks()
	default:
		break
	}

	return answer, err
}
