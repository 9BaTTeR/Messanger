package Configs

import (
	"ServerApp/Utility"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var instance Parames

// default pathes
const (
	config         = "configs.yml"
	defUserDir     = "UserDir/"
	defDialogsDir  = "DialogsDB/"
	defDataBaseDir = "DataBases/"
	defMediaDir    = "Media/"
	defDefaultsDB  = "Defaults/"
	defLogs        = "Logs/"
	defFiles       = "Config/"
)

// default values
const (
	hostname     = "hostname.com"
	port         = 80
	adminport    = 4000
	defaultRSA   = "mail_secure.key"
	htmltemplate = "EmailVerifications.html"
)

// all names
const (
	nameFDB       = "friends.db"
	nameBLDB      = "blacklist.db"
	userDB        = "user.db"
	historyDB     = "dialogs.db"
	sessiaDB      = "sessia.db"
	messagesDB    = "messages.db"
	settingDialog = "settings.db"
	alldialogs    = "chats.db"
	certsDB       = "certs.db"
	emailDB       = "email.db"
	lastID        = "id.log"
	endServer     = "shutdown.log"
	server        = "server.log"
	internalRSA   = "RSA.key"
)

// Получаем путь сертификата
func CertTLS() string {
	return instance.TLS.Certs
}

// Получаем путь ключа сертификата
func KeyTLS() string {
	answer := instance.TLS.Key
	return answer
}

func IsTLS() bool {
	return instance.istls
}

// Friends.db
func NameFDB() string {
	return nameFDB
}

// Blacklist.db
func NameBLDB() string {
	return nameBLDB
}

// User.db
func NameUDB() string {
	return userDB
}

// Dialogs.db
func NameHDB() string {
	return historyDB
}

// Sessia.db
func NameSDB() string {
	return sessiaDB
}

// Messages.db
func NameMDB() string {
	return messagesDB
}

// Settings.db
func SettingDDB() string {
	return settingDialog
}

// Chats.db
func NameChDB() string {
	return alldialogs
}

// Certs.db
func NameCsDB() string {
	return certsDB
}

// Email.db
func NameEDB() string {
	return emailDB
}

func SetShutdown(state bool) {
	instance.normalshutdown = state
}

func Shutdown() bool {
	return instance.normalshutdown
}

func Port() string {
	return fmt.Sprintf(":%v", instance.Port)
}

func AdminPort() string {
	return fmt.Sprintf(":%v", instance.AdminPort)
}

func Domain() string {
	return instance.HostName
}

func Initialize() error {
	path := Utility.Combine([]string{defFiles, config})
	exists, err := Utility.Exists(path)
	if err != nil {
		return fmt.Errorf("файл конфига необнаружен. Подробности: %w", err)
	}
	if !exists {
		content, err := defaultWrite(path)
		instance = content
		return err
	}
	err = instance.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла")
	}
	validity := instance.validate()
	if !validity {
		return fmt.Errorf("конфиг невалиден")
	}
	return nil
}

func defaultWrite(path string) (Parames, error) {
	instance := Parames{
		HostName:  hostname,
		Port:      port,
		AdminPort: adminport,
		Pathes: Pathes{
			UserDir:   defUserDir,
			DataBases: defDataBaseDir,
			Media:     defMediaDir,
			DialogDir: defDialogsDir,
		},
		Mail: Mail{
			PathRSA:      defaultRSA,
			TemplateHTML: htmltemplate,
		},
	}
	yml, err := yaml.Marshal(instance)
	if err != nil {
		return Parames{}, fmt.Errorf("ошибка маршализации конфига. Подробности: %w", err)
	}
	err = Utility.RewriteFile(path, string(yml))
	if err != nil {
		return Parames{}, fmt.Errorf("ошибка записи конфига по умолчанию. Подробности: %w", err)
	}
	return instance, nil
}

func (p *Parames) ReadFile(path string) error {
	config := []byte(readfile(path))
	err := yaml.Unmarshal(config, &p)
	return err
}

func (p *Parames) validate() bool {
	if p.AdminPort == 0 {
		return false
	}
	if strings.Replace(p.HostName, " ", "", -1) == "" {
		return false
	}
	if strings.Replace(p.Mail.PathRSA, " ", "", -1) == "" {
		return false
	}
	if strings.Replace(p.Mail.TemplateHTML, " ", "", -1) == "" {
		return false
	}
	if strings.Replace(p.Pathes.DataBases, " ", "", -1) == "" {
		return false
	}
	if strings.Replace(p.Pathes.DialogDir, " ", "", 1) == "" {
		return false
	}
	if strings.Replace(p.Pathes.Media, " ", "", 1) == "" {
		return false
	}
	if strings.Replace(p.Pathes.UserDir, " ", "", 1) == "" {
		return false
	}
	if p.Port == 0 {
		return false
	}
	if strings.ReplaceAll(p.TLS.Certs, " ", "") != "" {
		p.istls = true
	}
	return true
}
