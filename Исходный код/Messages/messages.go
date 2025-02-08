package message

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	friends "ServerApp/UserData/Friends"
	user "ServerApp/UserData/User"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	AddToDump = "INSERT INTO DMessages " +
		"VALUES (?,?,?,?,?,?,?,?,?);"
	AddToOrig    = "INSERT INTO OMessages VALUES (?,?,?);"
	EdOrig       = "UPDATE OMessages SET Date = ?, DMKey = ? WHERE OMKEY = ?;"
	DelOrig      = "DELETE FROM OMessages WHERE OMKEY = ?;"
	AddMedia     = "INSERT INTO Media(Hash,OrderBy,DMKey) VALUES (?,?,?);"
	FindMembers  = "SELECT IdUser FROM Members;"
	FindMessageForUpdate = "SELECT DMessages.IdUser, DMessages.Date, DMessages.ForwardedKey, DMessages.Read, DMessages.Important, " +
		"DMessages.DeletedAt, DMessages.UpdateAt FROM DMessages, OMessages " +
		"WHERE DMessages.DMKey == OMessages.DMKey AND OMessages.OMKey = ? GROUP BY DMessages.DMKey LIMIT 1;"
	internalerror = "Внутренняя ошибка"
)

func (mr *MessageRequest) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &mr)
	if err != nil {
		return err
	}
	return nil
}
func (mr *MessageRequest) Compose() ([]byte, error) {
	data, err := json.Marshal(&mr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Главный метод, выполняет все необходимые функции для добавления, редактирования, удаления у пользователей сообщений
func (mr *MessageRequest) WorkMsgToDB() (responces.MesagesResponce, error) {
	resp, err := mr.fillStruct()
	if err != nil {
		return resp, fmt.Errorf("ошибка наполнения структуры MessageRequest. Подробности: %w", err)
	}
	println("Y")
	if mr.Message.DMessages.Content == "" && mr.Operation!= "delete" && mr.Operation!= "deleteOnlyMe" && len(mr.Message.Hash) == 0{
		resp.Status = resp.Status.BadRequest("Пустое сообщение.")
		return resp, fmt.Errorf("пустое сообщение. Подробности: %w", err)
	}
	println("DD")
	//Только пользователь
	switch mr.Operation {
	case "deleteOnlyMe":
		resp, err := mr.addToPerson()
		if err != nil {
			return resp, err
		}
		return resp, nil
	}
	//Беседа
	if !TypeDialog(mr.Message.DMessages.PkDialog) {
		result, err := CheckUserInDialog(mr.Message.DMessages.idUser, mr.Message.DMessages.PkDialog)
		if err != nil {
			return resp, err
		}
		if !result{
			resp.Status = resp.Status.BadRequest("Пользователь не участник диалога.")
			return resp, nil
		}
		resp, err = mr.workToDialog()
		if err != nil {
			return resp, err
		}
		return resp, nil
	}

	//Личная переписка
	resp, err = mr.addToPersons()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (mr *MessageRequest) fillStruct() (responces.MesagesResponce, error) {
	uid, err := mr.Sessia.ParseID()
	if err != nil {
		return responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}, err
	}
	mr.Message.DMessages.idUser = uid.IDConversion()
	resp, err := mr.ValidDialogName()
	if err != nil {
		return resp, err
	}
	mr.Message.DMessages.senderUser = mr.Message.DMessages.idUser
	// //Заполнение OMKey, если есть
	mr.OMessage.OMKey = mr.Message.DMessages.DMKey
	//Генерация ключа DMKey
	timeTemp := Utility.GetTime()
	mr.OMessage.DMessages.Date = timeTemp
	mr.OMessage.Date = timeTemp
	mr.OMessage.DMessages.DMKey = mr.OMessage.genKeyMSG()
	
	//Генерация ключа OMKey
	if mr.Operation == "add"{

		mr.OMessage.Date = timeTemp + "OMKey"

		mr.OMessage.OMKey = mr.OMessage.genKeyMSG()
	}
	mr.Message.DMessages.Content = strings.TrimSpace(mr.Message.DMessages.Content)
	return responces.MesagesResponce{MsgKey: "-1"}, nil
}

// Добавить сообщения в папки пользователей личной переписки
func (mr *MessageRequest) addToPersons() (responces.MesagesResponce, error) {
	resp, err := mr.addToPerson()
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления 1 пользователю. Подробности: %v", err)
	}
	mrUserTwo := MessageRequest{}
	mrUserTwo = *mr
	mrUserTwo.Message.DMessages.PkDialog = strconv.FormatUint(uint64(mr.Message.DMessages.idUser), 10)
	mrUserTwo.Message.DMessages.idUser, err = strconv.ParseUint(mr.Message.DMessages.PkDialog, 10, 32)
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления преобразования в id пользователя. Подробности: %v", err)
	}
	resp, err = mrUserTwo.addToPerson()
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления 2 пользователю. Подробности: %v", err)
	}
	return resp, err
}

// Добавить сообщения в папку 1 пользователя личной переписки
func (mr *MessageRequest) addToPerson() (responces.MesagesResponce, error) {
	//Добавляем первому пользователю
	var err error
	mr.sqlInstance, err = openDB(mr.pathUserDialogDb())
	if err != nil {
		return responces.MesagesResponce{}, fmt.Errorf("ошибка открытия БД пользователю. Подробности: %v", err)
	}
	resp, err := mr.workMsg()
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления пользователю сообщения. Подробности: %v", err)
	}
	err = user.AddHistoryDialog([]uint64{mr.Message.DMessages.idUser}, mr.Message.DMessages.PkDialog)
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления диалога в историю пользователю. Подробности: %v", err)
	}
	return resp, err
}

func (mr *MessageRequest) pathUserDialogDb() string {
	uid := uid.ConvertionID(mr.Message.DMessages.idUser)
	pathuser := Utility.Combine([]string{uid.PathDbUserId(), "Dialogs",  mr.Message.DMessages.PkDialog, Configs.NameMDB()})
	return pathuser
}

// Добавить сообщения в папки пользователей беседы
func (mr *MessageRequest) workToDialog() (responces.MesagesResponce, error) {
	m, err := AllMembers(mr.Message.DMessages.PkDialog)
	if err != nil {
		return responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}, err
	}
	for i := 0; i < len(m.IdUsers); i++ {
		mrdialogs := mr
		mrdialogs.Message.DMessages.idUser = m.IdUsers[i]
		mrdialogs.sqlInstance, err = openDB(mrdialogs.pathUserDialogDb())
		if err != nil {
			return responces.MesagesResponce{}, err
		}
		resp, err := mrdialogs.workMsg()
		if err != nil {
			return resp, err
		}
	}
	resp, err := mr.workMsgToDialog()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Добавить сообщения в каталог беседы
func (mr MessageRequest) workMsgToDialog() (responces.MesagesResponce, error) {
	pathDialog := Utility.Combine([]string{Configs.Dialogs(), mr.Message.DMessages.PkDialog, Configs.NameMDB()})

	var err error
	mr.sqlInstance, err = openDB(pathDialog)
	if err != nil {
		return responces.MesagesResponce{}, nil
	}
	defer mr.sqlInstance.Close()
	resp, err := mr.workMsg()
	if err != nil {
		return resp, fmt.Errorf("ошибка добавления сообщения в каталог беседы. Подробности: %v", err)
	}
	return resp, nil
}

// Добавляет к 1 пользователю сообщения (DMessage, OMessage, при наличии Media)
func (mr MessageRequest) workMsg() (responces.MesagesResponce, error) {
	//Создание сообщения
	resp, err := mr.workMsgs(*mr.sqlInstance)
	if err != nil {
		return resp, fmt.Errorf("ошибка работы с сообщениями. Подробности %v", err)
	}
	if len(mr.Message.Hash) == 0 {
		return resp, nil
	}
	//Создание медиа
	err = mr.Message.addMedia(*mr.sqlInstance)
	if err != nil {
		return responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}, err
	}
	return resp, nil
}

// Создает сообщение у пользователя
func (mr *MessageRequest) workMsgs(s sql.DB) (responces.MesagesResponce, error) {
	resp := responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}
	err := mr.Message.DMessages.checkMessagesDb()
	if err != nil {
		return resp, err
	}
	if err != nil {
		return resp, fmt.Errorf("ошибка подключения sql. Подробности: %v", err)
	}
	//Добавление
	switch mr.Operation {
	case "add":
		err = mr.addMsg(s)
		if err != nil {
			return resp, fmt.Errorf("ошибка вставки записи. Подробности: %v", err)
		}
		resp.Status = resp.Status.OK("Успешно отправлено")
	//Редактирование
	case "edit":
		resp, err = mr.editMsg(s)
		if err != nil {
			return resp, fmt.Errorf("ошибка обновления записи Подробности: %v", err)
		}
	case "delete", "deleteOnlyMe":
		resp, err = mr.delMsg(s)
		if err != nil {
			return resp, fmt.Errorf("ошибка удаления записи. Подробности: %v", err)
		}
	}
		
	resp.MsgKey = mr.OMessage.OMKey
	return resp, nil
}

func (mr *MessageRequest) addMsg(s sql.DB) error {
	mr.Message.DMessages.DMKey = mr.OMessage.DMessages.DMKey
	mr.OMessage.Date = mr.OMessage.DMessages.Date
	mr.Message.DMessages.Date = mr.OMessage.Date
	err := mr.Message.DMessages.addToDMessages(s)
	if err != nil {
		return fmt.Errorf("ошибка вставки записи в таблицу DUMPMessages. Подробности: %v", err)
	}
	err = mr.OMessage.addToOMessages(s)
	if err != nil {
		return fmt.Errorf("ошибка вставки записи в таблицу OMessages. Подробности: %v", err)
	}
	return nil
}

func (mr *MessageRequest) editMsg(s sql.DB) (responces.MesagesResponce, error) {
	mr.OMessage.DMessages.UpdateAt = Utility.GetTime()
	resp, err := mr.prepareForChangeMsg(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка подготовки собщения для редактирования. Подробности: %v", err)
	}
	err = mr.OMessage.editOMessages(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка редактирования записи в таблице OMessages. Подробности: %v", err)
	}
	resp.Status = resp.Status.OK("Успешно отредактировано")
	return resp, nil
}

func (mr *MessageRequest) delMsg(s sql.DB) (responces.MesagesResponce, error) {
	mr.OMessage.DMessages.DeletedAt = Utility.GetTime()
	mr.OMessage.DMessages.DMKey = mr.OMessage.OMKey
	resp, err := mr.prepareForChangeMsg(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка подготовки собщения для удаления. Подробности: %v", err)
	}
	err = mr.OMessage.delOmessages(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка удаления записи в таблице OMessages. Подробности: %v", err)
	}
	resp.Status = resp.Status.OK("Успешно удалено")
	return resp, nil
}

func (mr *MessageRequest) prepareForChangeMsg(s sql.DB) (responces.MesagesResponce, error) {
	resp := responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}
	mr.OMessage.DMessages.senderUser = mr.Message.DMessages.senderUser
	err := mr.OMessage.infoForChange(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка получения инфорормации OMessages. Подробности: %v", err)
	}
	if mr.OMessage.DMessages.idUser != mr.OMessage.DMessages.senderUser && mr.Operation != "deleteOnlyMe"{
		resp.Status = resp.Status.BadRequest("пользователь пытался удалить сообщение, которое ему не принадлежит")
		return resp, fmt.Errorf("пользователь пытался удалить сообщение, которое ему не принадлежит")
	}
	mr.OMessage.DMessages.Content = mr.Message.DMessages.Content
	println(mr.OMessage.DMessages.DMKey)
	err = mr.OMessage.DMessages.addToDMessages(s)
	if err != nil {
		return resp, fmt.Errorf("ошибка вставки записи в таблицу DUMPMessages. Подробности: %v", err)
	}
	return resp, nil
}

// Проверяет наличие переписки, в случае отсуствия создает ее, отдает путь к ней.
func (dm DMessages) checkMessagesDb() error {
	uid := uid.ConvertionID(dm.idUser)
	pathuser :=  Utility.Combine([]string{uid.PathDbUserId(), "Dialogs", dm.PkDialog, Configs.NameMDB()}) 
	result, err := Utility.Exists(pathuser)
	if err != nil {
		return err
	}
	if result {
		return nil
	}
	_, err = Utility.CreateFile(pathuser)
	if err != nil {
		return err
	}
	pathdefault := Configs.MessagesDefDB()
	result, err = Utility.Exists(pathdefault)
	if err != nil {
		return err
	}
	if !result {
		err = CreateMessages()
		if err != nil {
			return err
		}
	}
	err = CopyMessagesDB(pathdefault, pathuser)
	if err != nil {
		return fmt.Errorf("ошибка копирования бд, подробности: %v", err)
	}
	return nil
}

// SQL методы
// Запрос на добавление медиа
func (m *Media) addMedia(s sql.DB) error {
	pathMedia := Configs.Media()
	for i := 0; i < len(m.Hash); i++ {
		path := pathMedia + "/" + m.Hash[i]
		exist, err := Utility.Exists(path)
		if err != nil {
			return fmt.Errorf("ошибка проверки наличия Media. Подробности: %v", err)
		}
		if !exist {
			return fmt.Errorf("медиа не существует. Подробности: %v", err)
		}
		_, err = s.Exec(AddMedia, m.Hash[i], m.Order[i], m.DMessages.DMKey)
		if err != nil {
			return fmt.Errorf("ошибка вставки записи в таблицу Media. Подробности: %v", err)
		}
	}
	return nil
}

// Запрос на получение информации сообщения до изменения.
func (om *OMessages) infoForChange(s sql.DB) error {
	println(om.OMKey)
	rows := s.QueryRow(FindMessageForUpdate, om.OMKey)
	var update sql.NullString
	var delete sql.NullString
	err := rows.Scan(&om.DMessages.idUser, &om.DMessages.Date, &om.DMessages.ForwardedKey, &om.DMessages.Read,
		&om.DMessages.Important, &delete, &update)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных запроса FindMessageForUpdate. Подробности: %v", err)
	}
	if om.DMessages.UpdateAt == "" {
		om.DMessages.UpdateAt = update.String
	}
	if om.DMessages.DeletedAt == "" {
		om.DMessages.DeletedAt = delete.String
	}
	return nil
}

// Запрос на изменения оригинального сообщения.
func (om *OMessages) editOMessages(s sql.DB) error {
	_, err := s.Exec(EdOrig, Utility.GetTime(), om.DMessages.DMKey, om.OMKey)
	if err != nil {
		return err
	}
	return nil
}

// Запрос на удаление сообщения. Не готов
func (om *OMessages) delOmessages(s sql.DB) error {
	_, err := s.Exec(DelOrig, om.OMKey)
	if err != nil {
		return err
	}
	return nil
}

// Запрос на создание OMessages.
func (om *OMessages) addToOMessages(s sql.DB) error {
	_, err := s.Exec(AddToOrig, om.OMKey, om.Date, om.DMessages.DMKey)
	if err != nil {
		return err
	}
	return nil
}

// Запрос на создание DMessages.
func (dm DMessages) addToDMessages(s sql.DB) error {

	_, err := s.Exec(AddToDump, dm.DMKey, dm.senderUser, dm.Date, dm.Content, dm.ForwardedKey,
		dm.Read, dm.Important, dm.DeletedAt, dm.UpdateAt)
	if err != nil {
		return err
	}
	return nil
}

// Генерация ключа сообщения.
func (om OMessages) genKeyMSG() string {
	data := om.Date 
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Найти всех участников диалога
func AllMembers(IdDialog string) (Members, error) {
	pathdialog := Utility.Combine([]string{Configs.Dialogs(), IdDialog})
	result, err := Utility.Exists(pathdialog)
	if err != nil {
		return Members{}, err
	}
	if !result {
		return Members{}, fmt.Errorf("диалога не существует")
	}
	pathdialog = Utility.Combine([]string{pathdialog, Configs.SettingDDB()})
	instance, err := openDB(pathdialog)
	if err != nil {
		return Members{}, err
	}
	defer instance.Close()
	rows, err := instance.Query(FindMembers)
	if err != nil {
		return Members{}, err
	}
	defer rows.Close()
	m := Members{}
	for rows.Next() {
		var member sql.NullInt64
		err = rows.Scan(&member)
		if err != nil {
			return Members{}, err
		}
		m.IdUsers = append(m.IdUsers, uint64(member.Int64))
	}
	return m, nil
}

// Проверка на личную переписку, на беседу
func TypeDialog(pkDialog string) bool {
	_, err := strconv.ParseUint(pkDialog, 10, 32)
	return err == nil
}

// Проверка названия диалога
func (mr MessageRequest) ValidDialogName() (responces.MesagesResponce, error) {
	resp := responces.MesagesResponce{Status: (responces.Responce{}).InternalError(internalerror), MsgKey: "-1"}
	if mr.Message.DMessages.PkDialog == "" {
		resp.Status = resp.Status.BadRequest("id диалога был пустым")
		return resp, fmt.Errorf("id диалога был пустым")
	}
	if !TypeDialog(mr.Message.DMessages.PkDialog) {
		return resp, nil
	}
	value, err := strconv.ParseUint(mr.Message.DMessages.PkDialog, 10, 32)
	if err != nil {
		return resp, err
	}
	if value == mr.Message.DMessages.idUser {
		resp.Status = resp.Status.BadRequest("отправка сообщения самому себе")
		return resp, fmt.Errorf("отправка сообщения самому себе")
	}
	u := user.User{}
	u.Sessia = mr.Sessia
	f, err := friends.LoadFriend(u, value)
	if err != nil {
		return resp, err
	}
	result := f.IsEmpty()
	if result {
		resp.Status = resp.Status.BadRequest("пользовать не находится в друзьях")
		return resp, fmt.Errorf("пользователь не находится в друзьях")
	}
	return resp, nil
}

func CheckUserInDialog(who uint64, idDialog string) (bool, error) {
	if TypeDialog(idDialog){
		return false, fmt.Errorf("попытка проверить наличие пользователя в диалоге, когда это переписка.")
	}
	m, err := AllMembers(idDialog)
	if err != nil{
		return false, fmt.Errorf("ошибка получения всех пользователей. Подробности: %w", err)
	}
	for i := 0; i < len(m.IdUsers); i++ {
		if who == m.IdUsers[i] { 
			return true, nil
		}
	}
	return false, nil
}
