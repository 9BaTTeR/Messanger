package getter

import (
	"ServerApp/Configs"
	msg "ServerApp/Messages"
	responces "ServerApp/Responces"
	friends "ServerApp/UserData/Friends"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	FindDialogInHistory = "SELECT Hash FROM HistoryDialogsUser WHERE UpdateAt > ? ORDER BY UpdateAt LIMIT ? OFFSET ?;"
	FindDialog          = "SELECT Name, Photo FROM Dialog;"
	FindLogin           = "SELECT Login, Photo FROM User;"
	FindMessages        = "SELECT DMessages.*, OMessages.OMKey FROM DMessages, OMessages WHERE DMessages.DMKey == OMessages.DMKey AND " +
		"OMessages.Date > ? GROUP BY DMessages.DMKey ORDER BY OMessages.Date LIMIT ? OFFSET ?;"
	FindDeletedMessages = "SELECT DMKey FROM DMessages WHERE Date > ? AND DeletedAt != '' ORDER BY Date LIMIT ? OFFSET ?;"
	FindMedia   = "SELECT Hash, OrderBy FROM Media WHERE DMKey = ?;"
	FindMessage = "SELECT DMessages.IdUser, DMessages.Date, DMessages.Content, DMessages.ForwardedKey, DMessages.Read, DMessages.Important, " +
		"DMessages.DeletedAt, DMessages.UpdateAt FROM DMessages, OMessages WHERE DMessages.DMKey == OMessages.DMKey AND " +
		"OMessages.OMKey = ? GROUP BY DMessages.DMKey LIMIT 1;"
	FindMessageCount = "SELECT COUNT(*) FROM OMessages;"
	FindDialogCount  = "SELECT COUNT(*) FROM HistoryDialogsUser;"
	FindAllMedia = "SELECT Hash FROM Media GROUP BY Hash LIMIT ? OFFSET ?;"
)

func (gd *GetDialogs) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &gd)
	if err != nil {
		return err
	}
	return nil
}
func (gd *GetDialogs) Compose() ([]byte, error) {
	data, err := json.Marshal(&gd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (mr *GetMessages) Compose() ([]byte, error) {
	data, err := json.Marshal(&mr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (mr *GetMessages) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &mr)
	if err != nil {
		return err
	}
	return nil
}

func (gm *GetMedia) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &gm)
	if err != nil {
		return err
	}
	return nil
}
func (gm *GetMedia) Compose() ([]byte, error) {
	data, err := json.Marshal(&gm)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (gu *GetUser) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &gu)
	if err != nil {
		return err
	}
	return nil
}
func (gu *GetUser) Compose() ([]byte, error) {
	data, err := json.Marshal(&gu)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (gamd *GetAllMediaDialog) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &gamd)
	if err != nil {
		return err
	}
	return nil
}
func (gamd *GetAllMediaDialog) Compose() ([]byte, error) {
	data, err := json.Marshal(&gamd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//Вытягивание диалогов

func (gd GetDialogs) ReqDialogs() (responces.GetDialogs, error) {
	if gd.Date == "" {
		gd.Date = "-1"
	}
	resp := responces.GetDialogs{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), Dialogs: []responces.Dialog{}}
	if gd.Take == 0 || gd.Take >= 50 {
		resp.Status = responces.Responce{}.BadRequest("некорректное значение аргумента Take")
		return resp, fmt.Errorf("некорректное значение аргумента Take")
	}
	uid, err := gd.Sessia.ParseID()
	if err != nil {
		return resp, fmt.Errorf("не удается определить id пользователя. Подробности: %w", err)
	}
	pathHistory := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameHDB()})
	instance, err := openDB(pathHistory)
	if err != nil {
		return resp, fmt.Errorf("не удается открыть подключение sql. Подробности: %w", err)
	}
	defer instance.Close()
	resp, err = gd.findDialogs(*instance)
	if err != nil {
		return resp, fmt.Errorf("не удается найти диалог. Подробности: %w", err)
	}
	resp.CountDialog, err = gd.findDialogCount(*instance)
	if err != nil {
		return resp, fmt.Errorf("не удается количество диалогов. Подробности: %w", err)
	}
	return resp, nil
}
func (gd GetDialogs) findDialogs(s sql.DB) (responces.GetDialogs, error) {
	resp := responces.GetDialogs{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), Dialogs: []responces.Dialog{}}
	rows, err := s.Query(FindDialogInHistory, gd.Date, gd.Take, gd.Skip)
	if err != nil {
		return resp, fmt.Errorf("ошибка получения информации dialogs.db Подробности: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		d := responces.Dialog{}
		err := rows.Scan(&d.Hash)
		if err != nil {
			return resp, fmt.Errorf("ошибка вытягивания информации dialogs.db. Подробности: %w", err)
		}
		resp.Dialogs = append(resp.Dialogs, d)
	}

	resp, err = gd.findInfo(resp)
	if err != nil {
		return resp, fmt.Errorf("ошибка поиска информации диалога. Подробности: %w", err)
	}
	resp.Status = resp.Status.OK("Список сформирован.")
	return resp, nil
}

func (gd GetDialogs) findInfo(hd responces.GetDialogs) (responces.GetDialogs, error) {
	for i := 0; i < len(hd.Dialogs); i++ {
		var err error
		if !msg.TypeDialog(hd.Dialogs[i].Hash) {
			//Dialogs
			hd.Dialogs[i], err = findInfoDialog(hd.Dialogs[i])
			if err != nil {
				return hd, fmt.Errorf("не удается найти диалог. Подробности: %w", err)
			}
			continue
		}
		//messages
		hd.Dialogs[i], err = gd.findInfoPerson(hd.Dialogs[i])
		if err != nil {
			return hd, fmt.Errorf("не удается найти переписку. Подробности: %w", err)
		}
	}
	return hd, nil
}

func findInfoDialog(d responces.Dialog) (responces.Dialog, error) {
	//pathFolderDialog := Configs.Path{}.FolderDialogs().Path()
	pathFolderDialog := Configs.Dialogs()
	pathDialog := Utility.Combine([]string{pathFolderDialog, d.Hash, Configs.SettingDDB()})
	instance, err := openDB(pathDialog)
	if err != nil {
		return d, fmt.Errorf("не удается открыть подключение sql к dialog.db. Подробности: %w", err)
	}
	defer instance.Close()
	rows, err := instance.Query(FindDialog)
	if err != nil {
		return d, fmt.Errorf("не удается выполнить запрос для dialog.db. Подробности: %w", err)
	}
	defer rows.Close()
	var photo sql.NullString
	for rows.Next() {
		err := rows.Scan(&d.Name, &photo)
		if err != nil {
			return d, fmt.Errorf("не удается просканировать данные dialog.db. Подробности: %w", err)
		}
		d.Photo = photo.String
	}

	return d, nil
}

func (gd GetDialogs) findInfoPerson(d responces.Dialog) (responces.Dialog, error) {
	id, err := strconv.ParseUint(string(d.Hash), 10, 64)
	if err != nil {
		return d, fmt.Errorf("не удается определить id пользователя. Подробности: %w", err)
	}
	uid := uid.ConvertionID(id)
	pathUserdb := Utility.Combine([]string{uid.PathDbUserId(), Configs.NameUDB()})
	instance, err := openDB(pathUserdb)
	if err != nil {
		return d, fmt.Errorf("не удается открыть подключить sql к user.db. Подробности: %w", err)
	}
	defer instance.Close()
	rows, err := instance.Query(FindLogin)
	if err != nil {
		return d, fmt.Errorf("не удается выполнить запрос для user.db. Подробности: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&d.Name, &d.Photo)
		if err != nil {
			return d, fmt.Errorf("не удается просканировать данные user.db. Подробности: %w", err)
		}
	}
	return d, nil
}

//Найти количество диалогов
func (gm GetDialogs) findDialogCount(s sql.DB) (uint64, error) {
	rows := s.QueryRow(FindDialogCount)
	var count sql.NullInt64
	err := rows.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества сообщений. Подробности: %+v", err)
	}	
	return uint64(count.Int64), nil
}


//Вытягивание сообщений

func (gm GetMessages) ReqMessages() (responces.GetMessages, error) {
	resp := responces.GetMessages{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), Msg: []responces.Media{}}
	uid, err := gm.Sessia.ParseID()
	if err != nil {
		return resp, err
	}
	path := Utility.Combine([]string{uid.PathDbUserId(), "Dialogs", gm.PkDialog, Configs.NameMDB()})
	instance, err := openDB(path)
	if err != nil {
		return resp, err
	}
	defer instance.Close()
	key := false
	if len(gm.OMKey) > 3 {
		key = true
	}
	resp, err = gm.findMessage(*instance, key)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Найти все сообщения
func (gm GetMessages) findMessage(s sql.DB, key bool) (responces.GetMessages, error) {
	resp := responces.GetMessages{}
	var err error
	if key {
		resp, err = gm.findMessageInfo(s)
	} else {
		resp, err = gm.findMessagesInfo(s)
	}
	if err != nil {
		return resp, err
	}
	resp, err = findMedia(s, resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Найти определенное сообщение
func (gm GetMessages) findMessageInfo(s sql.DB) (responces.GetMessages, error) {
	resp := responces.GetMessages{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), Msg: []responces.Media{}}
	rows := s.QueryRow(FindMessage, gm.OMKey)
	m := responces.Media{}
	err := rows.Scan(&m.DMessages.IdUser, &m.DMessages.Date, &m.DMessages.Content,
		&m.DMessages.ForwardedKey, &m.DMessages.Read, &m.DMessages.Important,
		&m.DMessages.DeleteAt, &m.DMessages.UpdateAt)
	resp.Msg = append(resp.Msg, m)
	if err == sql.ErrNoRows {
		resp.Status = resp.Status.BadRequest("Сообщения не существует.")
		return resp, fmt.Errorf("данные были пустые. Подробности: %+v", err)

	} else if err != nil {
		resp.Status = resp.Status.InternalError("Внутренняя ошибка.")
		return resp, fmt.Errorf("ошибка чтения данных. Подробности: %+v", err)
	}
	resp, err = findMedia(s, resp)
	if err != nil {
		return resp, err
	}
	resp.Status = resp.Status.OK("Сообщение сформировано.")
	return resp, nil
}
//Найти все сообщения по запросу
func (gm GetMessages) findMessagesInfo(s sql.DB) (responces.GetMessages, error) {
	resp := responces.GetMessages{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), Msg: []responces.Media{}}
	if gm.Date == "" {
		gm.Date = "-1"
	}
	if gm.Take == 0 || gm.Take >= 50 {
		resp.Status = responces.Responce{}.BadRequest("некорректное значение аргумента Take")
		return resp, fmt.Errorf("некорректное значение аргумента Take")
	}
	rows, err := s.Query(FindMessages, gm.Date, gm.Take, gm.Skip)
	if err != nil {
		return resp, err
	}
	defer rows.Close()
	for rows.Next() {
		m := responces.Media{}
		err := rows.Scan(&m.DMessages.DMKey, &m.DMessages.IdUser, &m.DMessages.Date, &m.DMessages.Content,
			&m.DMessages.ForwardedKey, &m.DMessages.Read, &m.DMessages.Important,
			&m.DMessages.DeleteAt, &m.DMessages.UpdateAt, &m.DMessages.OMKey)
		if err != nil {
			return resp, err
		}
		resp.Msg = append(resp.Msg, m)
	}
	
	resp.CountMsg, err = gm.findMessageCount(s)
	if err != nil{
		resp.Status = responces.Responce{}.InternalError("Внутренняя ошибка")
		return resp, err
	}
	resp.Status = resp.Status.OK("Список сообщений сформирован.")
	return resp, nil
}
//Найти количество сообщения
func (gm GetMessages) findMessageCount(s sql.DB) (uint64, error) {
	rows := s.QueryRow(FindMessageCount)
	var count sql.NullInt64
	err := rows.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества сообщений. Подробности: %+v", err)
	}	
	return uint64(count.Int64), nil
}

func findMedia(s sql.DB, resp responces.GetMessages) (responces.GetMessages, error) {
	for i := 0; i < len(resp.Msg); i++ {
		rows, err := s.Query(FindMedia, resp.Msg[i].DMessages.DMKey)
		if err != nil {
			return resp, err
		}
		defer rows.Close()
		m := resp.Msg[i]
		for rows.Next() {
			var hash sql.NullString
			var row sql.NullInt64
			err := rows.Scan(&hash, &row)
			if err != nil {
				return resp, err
			}
			m.Hash = append(m.Hash, hash.String)
			m.Order = append(m.Order, uint64(row.Int64))
		}
		resp.Msg[i] = m
	}
	return resp, nil
}

//Вытягивание удаленных сообщений

func (gm GetMessages) GetDelMsg() (responces.GetDelMessages, error) {
	resp := responces.GetDelMessages{Status: (responces.Responce{}).InternalError("Внутренняя ошибка"), MsgKey: []string{}}
	uid, err := gm.Sessia.ParseID()
	if err != nil {
		return resp, err
	}
	pathDialog := Utility.Combine([]string{uid.PathDbUserId(), "Dialogs", gm.PkDialog, Configs.NameMDB()})
	instance, err := openDB(pathDialog)
	if err != nil{
		return resp, fmt.Errorf("ошибка открытия БД для удаленных сообщений. Подробности: %w", err)
	}
	defer instance.Close()
	rows, err := instance.Query(FindDeletedMessages, gm.Date, gm.Take, gm.Skip)
	if err != nil{
		return resp, fmt.Errorf("ошибка выполнения поиска удаленных сообщений. Подробности: %w", err)
	}
	defer rows.Close()
	for rows.Next(){
		var DMKey sql.NullString
		err = rows.Scan(&DMKey)
		if err != nil{
			return resp, fmt.Errorf("ошибка сканирования удаленных сообщений. Подробности: %w", err)
		}
		resp.MsgKey = append(resp.MsgKey, DMKey.String)
	}
	return resp, nil
}
//Вытягивание медиа-контента

func (gm GetMedia) GetMedia() (responces.GetMedia, error) {
	resp := responces.GetMedia{}
	resp.Hash = gm.Hash
	if gm.Hash == "" {
		return resp, fmt.Errorf("хэш был пустой")
	}
	path := Utility.Combine([]string{Configs.Media(), gm.Hash})
	f, err := os.Open(path)
	if err != nil {
		return resp, fmt.Errorf("ошибка открытия каталога медиа. Подробности:%v", err)
	}
	defer f.Close()
	list, err := f.Readdirnames(-1)
	if err != nil {
		return resp, fmt.Errorf("ошибка чтения файлов каталога медиа. Подробности:%v", err)
	}
	path = Utility.Combine([]string{path, list[0]})
	bytes, err := os.ReadFile(path)
	if err != nil {
		return resp, fmt.Errorf("ошибка чтения файла в каталоге медиа. Подробности:%v", err)
	}
	resp.Extension = strings.Split(list[0], ".")[1]
	resp.BytesBase64 = Utility.BASE64(bytes)
	resp.Status = resp.Status.OK("Файл сформирован.")
	return resp, nil
}

//Вытягивание данных пользователя

func (gu GetUser) GetUser() (responces.GetUser, error) {
	resp := responces.GetUser{Status: (responces.Responce{}).InternalError("Внутренняя ошибка")}
	resp.Id = gu.Id
	f, err := friends.LoadUserData(gu.Id)
	if err != nil {
		return resp, err
	}
	resp.Nickname = f.Nickname
	resp.Photo = f.Photo
	resp.Status = resp.Status.OK("Данные пользователя сформированы.")
	return resp, nil
}

//Вытягивание всех участников диалога

func (gu GetUser) GetUsersDialog() (responces.GetUsersDialog, error) {
	resp := responces.GetUsersDialog{Status: (responces.Responce{}).InternalError("Внутренняя ошибка")}
	if msg.TypeDialog(gu.PkDialog){
		uid, err := gu.Sessia.ParseID()
		if err != nil {
			return resp, err
		}
		value, err := strconv.ParseUint(gu.PkDialog, 10, 32)
		if err != nil{
			return resp, nil
		}
		resp.Status = resp.Status.OK("Данные пользователей переписки сформированы.")
		resp.UsersId = []uint64{value, uid.IDConversion()}
		return resp, nil
	}
	m, err := msg.AllMembers(gu.PkDialog)
	if err != nil {
		return resp, err
	}
	resp.UsersId = m.IdUsers
	resp.Status = resp.Status.OK("Данные пользователей диалога сформированы.")
	return resp, nil
}

//Вложения диалога

func (gamd GetAllMediaDialog) DialogMedia() (responces.GetAllMedia, error){
	resp := responces.GetAllMedia{Status: (responces.Responce{}).InternalError("Внутренняя ошибка")}
	if gamd.Take == 0 || gamd.Take >= 50 {
		return resp, fmt.Errorf("некорректное значение аргумента Take")
	}
	uid, err := gamd.Sessia.ParseID()
	if err != nil {
		return resp, err
	}
	pathDialog := Utility.Combine([]string{uid.PathDbUserId(), "Dialogs", gamd.PkDialog, Configs.NameMDB()})
	instance, err := openDB(pathDialog)
	if err != nil {
		return resp, err
	}
	defer instance.Close()
	resp, err = gamd.findAllMedia(*instance)
	if err != nil {
		return resp, err
	}
	resp.Status = resp.Status.OK("Медиа сформированы")
	return resp, nil
}
func (gamd GetAllMediaDialog) findAllMedia(s sql.DB) (responces.GetAllMedia, error){
	resp := responces.GetAllMedia{Status: (responces.Responce{}).InternalError("Внутренняя ошибка")}
	rows, err := s.Query(FindAllMedia, gamd.Take, gamd.Skip)
	if err != nil{
		return resp, fmt.Errorf("медиа не найдены. Подробности: %w", err)
	}
	defer rows.Close()
	for rows.Next(){
		var hash sql.NullString
		err = rows.Scan(&hash)
		if err != nil{
			return resp, fmt.Errorf("ошибка сканирования медиа. Подробности: %w", err)
		}
		resp.Hash = append(resp.Hash, hash.String)
	}
	return resp, nil
}

func openDB(path string) (*sql.DB, error) {
	sql, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия бд. Подробности %w", err)
	}
	return sql, nil
}