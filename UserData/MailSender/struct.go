package MailSender

type AuthData struct {
	user     string
	password string
}
type acceptemail struct {
	Name      string
	Code      string
	Aboutuser aboutuser
}
type aboutuser struct {
	IP4      string //IP4 устройства
	Datetime string //Время запроса
	Device   string //Сведения об устройстве
}
