package dialogs

import (
	"ServerApp/Configs"
	msg "ServerApp/Messages"
	responces "ServerApp/Responces"
	friends "ServerApp/UserData/Friends"
	user "ServerApp/UserData/User"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	fillDialogToAll = "INSERT INTO ALLDIALOGS" +
		" VALUES (?,?);"
	fillDialogToHistory = "INSERT OR REPLACE INTO HistoryDialogsUser" +
		" VALUES (?,?);"
	FindMembers = "SELECT * FROM Members WHERE Hash = ;"
	AddMembers  = "INSERT INTO MEMBERS" +
		" VALUES (?, ?, ?, ?);"
	FillDialog    = "INSERT INTO Dialog(Name, Privacy, CreatedAt) VALUES (?,?,?)"
	UpdateName    = "UPDATE Dialog SET Name = ?;"
	UpdatePhoto   = "UPDATE Dialog SET Photo = ?;"
	UpdatePrivacy = "UPDATE Dialog SET Privacy = ?;"
	DeleteMember = "DELETE FROM Members WHERE IdUser = ?;"
)

func (dr *DialogRequest) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &dr)
	if err != nil {
		return err
	}
	return nil
}
func (dr *DialogRequest) Compose() ([]byte, error) {
	data, err := json.Marshal(&dr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (kr *KickResponce) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &kr)
	if err != nil {
		return err
	}
	return nil
}
func (kr *KickResponce) Compose() ([]byte, error) {
	data, err := json.Marshal(&kr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Главный метод, создаст полностью диалог
func (dr *DialogRequest) CopyDialog() (responces.DialogResponce, error) {
	resp := responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}
	uid, err := dr.Sessia.ParseID()
	if err != nil {
		return resp, err
	}
	dr.Dialog.idUserCreator = uid.IDConversion()
	dr.delCreatorInArray()
	if len(dr.Dialog.IdUsers) < 1{
		resp.Status = resp.Status.BadRequest("Количество пользователей меньше 1")
		return resp, fmt.Errorf("ошибка количества пользователей диалога. Подробности %w", err)
	}
	resp, err = dr.checkValidUsers()
	if err != nil {
		return resp, fmt.Errorf("ошибка проверки валидности пользователей диалога. Подробности %w", err)
	}
	if (resp != responces.DialogResponce{}) {
		return resp, nil
	}
	idCreator := strconv.FormatUint(dr.Dialog.idUserCreator, 10)
	dr.Dialog.CreatedAt = Utility.GetTime()
	hash := Utility.SHA256(idCreator + dr.Dialog.CreatedAt + dr.Dialog.Name)
	dr.Dialog.Hash = hash
	//Проверка на наличие БД для всех диалогов
	//pathDialogAll := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), namealldialogdb})
	pathDialogAll := Configs.ChatsDB()
	exist, err := Utility.Exists(pathDialogAll)
	if !exist {
		err = CreateAllDialogsDB()
	}
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Проверка наличия БД, в случае отсуствия - создает
	err = dr.Dialog.CheckDialog()
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Наполнение диалога
	err = dr.Dialog.fillDialog()
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Создания сообщений диалога
	err = dr.Dialog.MessageDbToDialog()
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Добавление диалога в БД всех диалогов
	err = dr.Dialog.AddDialogToAll()
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Добавление диалога в историю пользователей
	members := dr.Dialog.IdUsers
	members = append(members, dr.Dialog.idUserCreator)
	err = user.AddHistoryDialog(members, dr.Dialog.Hash)
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	//Добавление мемберов к диалогу + создание у них переписки диалога
	err = dr.Dialog.addMembersDialog()
	if err != nil {
		return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, err
	}
	return responces.DialogResponce{Status: (responces.Responce{}).OK("Диалог начат."), Hash: dr.Dialog.Hash}, nil
}

func (d Dialog) CheckDialog() error {
	//Проверка на наличие БД базового диалога
	//pathDialogBase := Utility.Combine([]string{Configs.Path{}.FolderDefDB().Path(), namedialogdb})
	pathDialogBase := Configs.SettingsDefDDB()
	exist, err := Utility.Exists(pathDialogBase)
	if !exist {
		err = d.createBaseDialog()
	}
	if err != nil {
		return fmt.Errorf("ошибка проверки наличия базового dialog.db. Подробности %w", err)
	}
	//Копирование диалога
	pathDialog := Utility.Combine([]string{Configs.Dialogs(), d.Hash, Configs.SettingDDB()})
	exist, err = Utility.Exists(pathDialog)
	if exist {
		return fmt.Errorf("данный диалог уже существует")
	}
	if err != nil {
		return fmt.Errorf("ошибка проверки наличия dialog.db. Подробности %w", err)
	}
	_, err = Utility.CreateFile(pathDialog)
	if err != nil {
		return fmt.Errorf("ошибка создания директории диалога. Подробности %w", err)
	}
	err = Utility.CopyFile(pathDialogBase, pathDialog)
	if err != nil {
		return fmt.Errorf("ошибка копирования dialog.db. Подробности %w", err)
	}
	return nil
}

func (d Dialog) fillDialog() error {
	//pathdialog := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), d.Hash, namedialogdb})
	pathdialog := Utility.Combine([]string{Configs.Dialogs(), d.Hash, Configs.SettingDDB()})
	Instance, err := openDB(pathdialog)
	if err != nil {
		return fmt.Errorf("ошибка открытия sql dialog.db. Подробности %w", err)
	}
	defer Instance.Close()
	_, err = Instance.Exec(FillDialog, d.Name, d.Private, d.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы dialog в dialog.db. Подробности %w", err)
	}
	return nil
}

func (d *Dialog) AddDialogToAll() error {
	//path := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), namealldialogdb})
	path := Configs.ChatsDB()
	instance, err := openDB(path)
	if err != nil {
		return fmt.Errorf("ошибка открытия sql allDialogs.db. Подробности %w", err)
	}
	_, err = instance.Exec(fillDialogToAll, d.Hash, d.Private)
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы всех диалогов. Подробности: %w", err)
	}
	defer instance.Close()
	return nil
}

func (d Dialog) MessageDbToDialog() error {
	//pathDialogMsg := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), d.Hash, namemessagedb})
	pathDialogMsg := Utility.Combine([]string{Configs.Dialogs(), d.Hash, Configs.NameMDB()})
	//pathMsg := Utility.Combine([]string{Configs.Path{}.FolderDefDB().Path(), namemessagedb})
	pathMsg := Configs.MessagesDefDB()
	exist, err := Utility.Exists(pathMsg)
	if !exist {
		err = msg.CreateMessages()
	}
	if err != nil {
		return fmt.Errorf("ошибка создания базового message.db. Подробности %w", err)
	}
	err = Utility.CopyFile(pathMsg, pathDialogMsg)
	if err != nil {
		return fmt.Errorf("ошибка копирования message.db в диалог. Подробности %w", err)
	}
	return nil
}

func (d Dialog) addMembersDialog() error {
	err := d.membersToDialog()
	if err != nil {
		return fmt.Errorf("ошибка наполнения таблицы members в dialog.db. Подробности %w", err)
	}
	d.IdUsers = append(d.IdUsers, d.idUserCreator)
	err = d.messagesToMembers()
	if err != nil {
		return fmt.Errorf("ошибка создания message.db у пользователей диалога. Подробности %w", err)
	}
	return nil
}

func (d Dialog) membersToDialog() error {
	date := Utility.GetTime()
	//pathdialog := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), d.Hash, namedialogdb})
	pathdialog := Utility.Combine([]string{Configs.Dialogs(), d.Hash, Configs.SettingDDB()})
	Instance, err := openDB(pathdialog)
	if err != nil {
		return fmt.Errorf("ошибка открытия sql dialog.db. Подробности %w", err)
	}
	defer Instance.Close()

	_, err = Instance.Exec(AddMembers, d.idUserCreator, "Создатель", date, "true")
	if err != nil {
		return fmt.Errorf("ошибка наполнения информации о создателе диалога. Подробности %w", err)
	}

	for i := 0; i < len(d.IdUsers); i++ {
		_, err := Instance.Exec(AddMembers, d.IdUsers[i], "Пользователь", date, "true")
		if err != nil {
			return fmt.Errorf("ошибка наполнения информации о пользователях диалога. Подробности %s", err)
		}
	}

	return nil
}

func (d Dialog) messagesToMembers() error {
	//pathDialogMsg := Utility.Combine([]string{Configs.Path{}.FolderDialogs().Path(), d.Hash, namemessagedb})
	pathDialogMsg := Utility.Combine([]string{Configs.Dialogs(), d.Hash, Configs.NameMDB()})
	for i := 0; i < len(d.IdUsers); i++ {
		uId := uid.ConvertionID(uint64(d.IdUsers[i]))
		pathToFolder := uId.PathDbUserId()
		err := d.checkDialogUser(pathToFolder)
		if err != nil {
			return fmt.Errorf("ошибка проверки диалога у пользователей диалога. Подробности %w", err)
		}
		pathToUserMsg := Utility.Combine([]string{pathToFolder, "Dialogs", d.Hash, Configs.NameMDB()})
		err = Utility.CopyFile(pathDialogMsg, pathToUserMsg)
		if err != nil {
			return fmt.Errorf("ошибка копирования message.db у пользователей диалога. Подробности %w", err)
		}
	}
	return nil
}

func (d Dialog) checkDialogUser(pathToFolder string) error {
	//Проверка на наличие диалога у пользователя
	pathDialog := Utility.Combine([]string{pathToFolder, "Dialogs", d.Hash})
	exist, err := Utility.Exists(pathDialog)
	if !exist {
		err = Utility.CreateFolders(pathDialog)
	}
	if err != nil {
		return fmt.Errorf("ошибка создания диалога у пользователя. Подробности %w", err)
	}
	return nil
}

func (dr DialogRequest) checkValidUsers() (responces.DialogResponce, error) {
	//Массив с юзеров с создателем
	users := dr.Dialog.IdUsers
	users = append(users, dr.Dialog.idUserCreator)
	//Сессия создателя
	u := user.User{}
	u.Sessia = dr.Sessia
	for i := 0; i < len(dr.Dialog.IdUsers); i++ {
		//Проверка наличия пользователя
		uId := uid.ConvertionID(users[i])
		exist, err := Utility.Exists(uId.PathDbUserId())
		if err != nil {
			return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, fmt.Errorf("ошибка проверки наличия пользователя. Подробности: %w", err)
		}
		if !exist {
			return responces.DialogResponce{Status: (responces.Responce{}).InternalError("пользователя не существует"), Hash: "-1"}, nil
		}
		//Проверка наличия в друзьях пользователя
		f, err := friends.LoadFriend(u, dr.Dialog.IdUsers[i])
		if err != nil {
			return responces.DialogResponce{Status: (responces.Responce{}).InternalError(internalerror), Hash: "-1"}, fmt.Errorf("ошибка загрузки друга пользователя. Подробности: %w", err)
		}
		exist = f.IsEmpty()
		if exist {
			return responces.DialogResponce{Status: (responces.Responce{}).InternalError("пользователь не находится в друзьях"), Hash: "-1"}, fmt.Errorf("пользователь не находится в друзьях")
		}
	}
	return responces.DialogResponce{}, nil
}

func (dr DialogRequest) ChangeDialog() (responces.NonAddResponce, error) {
	resp := responces.NonAddResponce{Status: (responces.Responce{}).InternalError(internalerror)}
	uid, err := dr.Sessia.ParseID()
	if err != nil {
		return resp, err
	}
	dr.Dialog.idUserCreator = uid.IDConversion()
	pathDialog := Utility.Combine([]string{Configs.Dialogs(), dr.Dialog.Hash, Configs.SettingDDB()})
	instance, err := openDB(pathDialog)
	if err != nil {
		return resp, fmt.Errorf("ошибка подключения БД диалога. Подробности: %w", err)
	}
	if dr.Dialog.Name != "" {
		err = dr.changeNameDialog(*instance)
		if err != nil{
			return resp, err
		}
	}
	if dr.Dialog.Photo != "" {
		err = dr.changePhoto(*instance)
		if err != nil{
			return resp, err
		}
	}
	if dr.Dialog.Private != "" {
		err = dr.changePrivacy(*instance)
		if err != nil{
			return resp, err
		}
	}
	if len(dr.Dialog.IdUsers) > 0 {
		respDialog, err := dr.checkValidUsers()
		resp.Status = respDialog.Status
		if err != nil{
			return resp, err
		}
		dr.delCreatorInArray()
		if len(dr.Dialog.IdUsers) < 1{
			resp.Status = resp.Status.BadRequest("Количество пользователей меньше 1")
			return resp, fmt.Errorf("ошибка количества пользователей диалога. Подробности %w", err)
		}
		err = dr.Dialog.addMembers(*instance)
		if err != nil{
			return resp, err
		}
		err = dr.Dialog.messagesToMembers()
		if err != nil{
			return resp, err
		}
	}
	defer instance.Close()
	resp.Status = resp.Status.OK("Данные успешно изменены.")
	return resp, nil
}

func (dr DialogRequest) changeNameDialog(s sql.DB) error {
	_, err := s.Exec(UpdateName, dr.Dialog.Name)
	if err != nil {
		return fmt.Errorf("ошибка смены названия диалога. Подробности: %w", err)
	}
	return nil
}
func (dr DialogRequest) changePhoto(s sql.DB) error {
	_, err := s.Exec(UpdatePhoto, dr.Dialog.Photo)
	if err != nil {
		return fmt.Errorf("ошибка смены названия диалога. Подробности: %w", err)
	}
	return nil
}
func (dr DialogRequest) changePrivacy(s sql.DB) error {
	_, err := s.Exec(UpdatePrivacy, dr.Dialog.Private)
	if err != nil {
		return fmt.Errorf("ошибка смены названия диалога. Подробности: %w", err)
	}
	return nil
}

func (d Dialog) addMembers(s sql.DB) error{
	date := Utility.GetTime()
	for i := 0; i < len(d.IdUsers); i++ {
		_, err := s.Exec(AddMembers, d.IdUsers[i], "Пользователь", date, "true")
		if err != nil {
			return fmt.Errorf("ошибка наполнения информации о пользователях диалога. Подробности %w", err)
		}
	}
	return nil
}

func (kr KickResponce) KickUser() (responces.NonAddResponce, error){
	resp := responces.NonAddResponce{Status: (responces.Responce{}).InternalError(internalerror)}
	result, err := msg.CheckUserInDialog(kr.Whom, kr.PkDialog)
	if err != nil{
		return resp, fmt.Errorf("ошибка проверки пользователя на наличие в диалоге. Подробности: %w", err)
	}
	if !result{
		resp.Status = resp.Status.BadRequest("Вы не состоите в диалоге")
		return resp, fmt.Errorf("вы не состоите в диалоге")
	}
	path := Utility.Combine([]string{Configs.Dialogs(), kr.PkDialog, Configs.SettingDDB()})
	instance, err := openDB(path)
	if err != nil{
		return resp, fmt.Errorf("ошибка открытия бд настроек диалога. Подробности: %w", err)
	}
	err = kr.kickUser(*instance)
	if err != nil{
		return resp, fmt.Errorf("ошибка изгнания пользователя. Подробности: %w", err)	
	}
	resp.Status = resp.Status.OK("Пользователь изгнан")
	return resp, nil
}

func (kr KickResponce) kickUser(s sql.DB) error{
	_, err := s.Exec(DeleteMember, kr.Whom)
	if err != nil {
		return fmt.Errorf("ошибка смены названия диалога. Подробности: %w", err)
	}
	return nil
}

func(dr *DialogRequest) delCreatorInArray(){
	newArr := []uint64{}
	for _, val := range dr.Dialog.IdUsers {
		if val != dr.Dialog.idUserCreator {
			newArr = append(newArr, val)
		}
	}
	dr.Dialog.IdUsers = newArr
}
