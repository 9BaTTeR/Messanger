package Configs

import (
	"fmt"
	"io"
	"os"
	"path"
)

type Path string

func SettingsDefDDB() string {
	return path.Join(DefaultsDB(), SettingDDB())
}

func InternalRSA() string {
	return path.Join(defFiles, internalRSA)
}

func UserDir() string {
	return instance.Pathes.UserDir
}

func DataBases() string {
	return instance.Pathes.DataBases
}

func Media() string {
	return instance.Pathes.Media
}

func Dialogs() string {
	return instance.Pathes.DialogDir
}

func DefaultsDB() string {
	result := path.Join(instance.Pathes.DataBases, defDefaultsDB)
	return result
}

func Logs() string {
	return defLogs
}

func configs() string {
	return path.Join(defFiles, config)
}

func FriendsDefDB() string {
	return path.Join(DefaultsDB(), nameFDB)
}
func BlackListDefDB() string {
	return path.Join(DefaultsDB(), nameBLDB)
}

func UserDefDB() string {
	return path.Join(DefaultsDB(), userDB)
}

func HistoryDefDB() string {
	return path.Join(DefaultsDB(), historyDB)
}

func SessiaDefDB() string {
	return path.Join(DefaultsDB(), sessiaDB)
}

func MessagesDefDB() string {
	return path.Join(DefaultsDB(), messagesDB)
}

func ChatsDB() string {
	return path.Join(instance.Pathes.DialogDir, alldialogs)
}

func CertsDB() string {
	return path.Join(instance.Pathes.DataBases, certsDB)

}

func EmailDB() string {
	return path.Join(instance.Pathes.DataBases, emailDB)
}

func MailRSA() string {
	return path.Join(defFiles, instance.Mail.PathRSA)
}
func MailTemplate() string {
	return path.Join(defFiles, instance.Mail.TemplateHTML)
}

func GetMailKey() string {
	return readfile(MailRSA())
}
func GetAuthHMTL() string {
	return readfile(MailTemplate())
}

func LastID() string {
	return path.Join(defLogs, lastID)
}
func ShutdownLog() string {
	return path.Join(defLogs, endServer)
}

func ServerLog() string {
	return path.Join(defLogs, server)
}

func readfile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	data := make([]byte, 64)
	output := ""
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		output += string(data[:n])
	}
	return output
}

// const (
// 	pathuserdata          = "../UserDir"
// 	pathToDB              = "../DataBases"
// 	defaultsDB            = pathToDB + "/Defaults"
// 	mail_private_key_path = "mail_secure.key"
// 	htmlauth_path         = "UserData/MailSender/EmailVerifications.html"
// 	logpath               = "service.log"
// 	userID                = "../UserDir/id.yml"
// 	emaildb               = "../DataBases/email.db"
// 	media                 = "../Media"
// 	DialogsDB             = "../DialogsDB"
// )
