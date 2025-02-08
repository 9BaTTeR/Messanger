package MailSender

import (
	"ServerApp/Configs"
	DeviceData "ServerApp/UserData/Sessia"
	"strings"

	"text/template"
)

// Для получения типа данных data воспользуйтесь методом Getdata()
func authHtmlWithData(templatedata acceptemail) string {
	t := template.Must(template.New("letter").Parse(Configs.GetAuthHMTL()))
	var outbuf strings.Builder
	t.Execute(&outbuf, templatedata)
	return outbuf.String()
}

// Метод получения сведений об устройстве источника запроса.
func AboutUser(ip4 string, datetime string, device DeviceData.Sessia) aboutuser {
	var out aboutuser
	out.IP4 = ip4
	out.Datetime = datetime
	out.Device = device.DeviceID.HostName
	return out
}

func (a AuthData) Support() AuthData {
	a.user = "Support"
	a.password = "258013006402582"
	return a
}
func (a AuthData) NoReply() AuthData {
	a.user = "noreply"
	a.password = "258013006402582"
	return a
}
