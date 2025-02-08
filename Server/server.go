package Server

// Базовый пакет
import (
	"ServerApp/Configs"
	dialog "ServerApp/Dialogs"
	getter "ServerApp/Getter"
	Logger "ServerApp/Logger"
	media "ServerApp/Media"
	message "ServerApp/Messages"
	resp "ServerApp/Responces"
	search "ServerApp/Search"
	Auth "ServerApp/UserData/Auth"
	goEmail "ServerApp/UserData/Email"
	friends "ServerApp/UserData/Friends"
	"ServerApp/UserData/Sessia"
	tc "ServerApp/UserData/TrustClient"
	RegUser "ServerApp/UserData/User"
	uid "ServerApp/UserData/UserID"
	"ServerApp/Utility"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

// string bool
const (
	_false = "false"
	_true  = "true"
)

func HowShutdown() bool {
	return Utility.ReadFile(Configs.ShutdownLog()) == _true
}

func Run_server() {
	Configs.SetShutdown(HowShutdown())
	log, err := Logger.Initialize()
	if err != nil {
		fmt.Printf("\n\n\nсбой инициализации логгера. Подробности: %w\n\n\n", err)
		return
	}

	err = Configs.Initialize()
	if err != nil {
		log.ChanLog <- fmt.Sprintf("\nсбой инициализации конфига. Подробности: %s\n", err)
	}
	Utility.RewriteFile(Configs.ShutdownLog(), _false)

	HttpServer := http.Server{
		Addr: Configs.Port(),
	}

	http2Server := http2.Server{}
	err = http2.ConfigureServer(&HttpServer, &http2Server)
	if err != nil {
		log.ChanLog <- fmt.Sprintf("ошибка конфигурирования сервера: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		echoPayload(w, r, log)
	})

	idleConnsClosed := make(chan struct{})

	go SigTerm(idleConnsClosed, &HttpServer, log)
	if Configs.IsTLS() {
		log.ChanLog <- fmt.Sprintf("\ngo бекенд: { HTTPVersion = 2 }; serving on https://%v/\n", Configs.Domain()+Configs.Port())
		err = HttpServer.ListenAndServeTLS(Configs.CertTLS(), Configs.KeyTLS()) //Запускаем TLS сервер в поток
		if err != nil {
			log.ChanLog <- fmt.Sprintf("Сервер : %v\n", err)
		}
	} else {
		log.ChanLog <- fmt.Sprintf("\ngo бекенд: { HTTPVersion = 2 }; serving on http://%v/\n", Configs.Domain()+Configs.Port())
		err = HttpServer.ListenAndServe() //Запускаем сервер в поток
		if err != nil {
			log.ChanLog <- fmt.Sprintf("Сервер : %v\n", err)
		}
	}
	Utility.RewriteFile(Configs.ShutdownLog(), _true)
	uid.WriteID()
	close(log.ChanLog)
}

func SigTerm(idleConnsClosed chan struct{}, HttpServer *http.Server, l Logger.Logger) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	if err := HttpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		l.ChanLog <- fmt.Sprintf("%w", err)
	}
	_ = uid.WriteID()
	time.Sleep(time.Second)
	close(idleConnsClosed)

}

func echoPayload(w http.ResponseWriter, req *http.Request, log Logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	defer req.Body.Close()
	log.ChanLog <- fmt.Sprintf("log: Client connected from %v\n", req.RemoteAddr)
	contents, err := io.ReadAll(req.Body)
	if err != nil {
		log.ChanLog <- fmt.Sprintf("%w", err)
		return
	}
	pathUrl := req.URL.Path[1:]
	verify, rp, err := tc.ConnVerify(contents, pathUrl)
	if err != nil {
		log.ChanLog <- err.Error()
		answer, err := rp.Responcer()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("\nошибка валидации соединения. Подробности: %s\n", err)
		}
		w.Write(answer)
		w.WriteHeader(int(rp.Responce.Code))
		return
	}
	if !verify {
		answer, err := rp.Responcer()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("\nошибка валидации соединения. Подробности: %s\n", err)
		}
		w.Write(answer)
		w.WriteHeader(int(rp.Responce.Code))
		return
	}
	validate, rr, err := Sessia.ValidateTocken(contents, pathUrl)
	if !validate || err != nil {
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		w.Write(answer)
		w.WriteHeader(int(rp.Responce.Code))
		return
	}
	switch pathUrl {
	case "Media":
		m := media.Media{}
		Utility.JSON(&m).Parse(string(contents))
		active, err := m.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при отправки медиа. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при отправки медиа Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при отправки медиа. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		media, err := m.CreateMedia()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}

		answer, err := resp.Resp(media).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "Auth":
		AuthData, err := Auth.TryAuth(string(contents), req.RemoteAddr)
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		responce, err := AuthData.Auth()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		answer, err := responce.Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		w.Write(answer)
		w.WriteHeader(200)
	case "Registration":
		regdata := RegUser.User{}
		err := Utility.JSON(&regdata).Parse(string(contents))
		if err != nil {
			log.ChanLog <- fmt.Sprintf("%s ошибка парсинга жсон", err)
			return
		}
		//regdata, err := UserData.ParseReg(contents, req.RemoteAddr, "Acer Aspire Test")
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		regcode, err := regdata.Registration(string(strings.Split(req.RemoteAddr, ":")[0]))
		if err != nil {
			regdata.Rollback()
			log.ChanLog <- fmt.Sprint(err.Error())
		}
		answer, err := resp.Resp(regcode).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprint(err.Error())
		}

		w.Write([]byte(answer))
		w.WriteHeader(200)

	case "Email":
		ac := goEmail.Request{}
		Utility.JSON(&ac).Parse(string(contents))
		a := req.RemoteAddr[:strings.Index(req.RemoteAddr, ":")]
		rr, err := ac.Do(a)
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "SendMessage":
		mr := message.MessageRequest{}
		Utility.JSON(&mr).Parse(string(contents))
		active, err := mr.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при отправки сообщания. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при отправки сообщания. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при отправки сообщания. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		msg, err := mr.WorkMsgToDB()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(msg).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "RequestFrBl":
		reqfr := friends.Request{}
		json.Unmarshal(contents, &reqfr)
		active, err := reqfr.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при работе с друзьями. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при работе с друзьями. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при работе с друзьями. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		reqfr.Date = time.Now()
		responce, err := reqfr.Do()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := responce.Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("\nОшибка создания ответа. Подробности: %w", err)
		}
		w.Write(answer)
		w.WriteHeader(200)
	case "CreateDialog":
		dr := dialog.DialogRequest{}
		Utility.JSON(&dr).Parse(string(contents))
		active, err := dr.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при создании диалога. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при создании диалога. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при создании диалога. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		rr, err := dr.CopyDialog()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetMessages":
		gm := getter.GetMessages{}
		Utility.JSON(&gm).Parse(string(contents))
		rr, err := gm.ReqMessages()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetDelMessages":
		gm := getter.GetMessages{}
		Utility.JSON(&gm).Parse(string(contents))
		rr, err := gm.GetDelMsg()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetDialogs":
		gd := getter.GetDialogs{}
		Utility.JSON(&gd).Parse(string(contents))
		rr, err := gd.ReqDialogs()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetMediaDialog":
		gamd := getter.GetAllMediaDialog{}
		Utility.JSON(&gamd).Parse(string(contents))
		rr, err := gamd.DialogMedia()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetMedia":
		gm := getter.GetMedia{}
		Utility.JSON(&gm).Parse(string(contents))
		rr, err := gm.GetMedia()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetUser":
		gu := getter.GetUser{}
		Utility.JSON(&gu).Parse(string(contents))
		rr, err := gu.GetUser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "GetUsersDialog":
		gu := getter.GetUser{}
		Utility.JSON(&gu).Parse(string(contents))
		rr, err := gu.GetUsersDialog()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "KickUser":
		kr := dialog.KickResponce{}
		Utility.JSON(&kr).Parse(string(contents))
		active, err := kr.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при изгнании пользователя. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изгнании пользователя. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изгнании пользователя. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		rr, err := kr.KickUser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := resp.Resp(rr).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "SearchRequest":
		req := search.Request{}
		err := Utility.JSON(&req).Parse(string(contents))
		if err != nil {
			body := resp.NonAddResponce{Status: resp.Responce{}.InternalError("внутренняя ошибка при поиске")}
			answer, _ := resp.Resp(body).Responser()
			w.Write(([]byte(answer)))
			w.WriteHeader(505)
			return
		}
		body, err := req.Do()
		if err != nil {
			body := resp.NonAddResponce{Status: resp.Responce{}.InternalError("внутренняя ошибка при поиске")}
			answer, _ := resp.Resp(body).Responser()
			w.Write(([]byte(answer)))
			w.WriteHeader(505)
			return
		}
		answer, _ := body.Responser()
		w.Write(([]byte(answer)))
		w.WriteHeader(505)
	case "Certificate":
		c := tc.Certificate{}
		Utility.JSON(&c).Parse(string(contents))

		cert, err := c.Append()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		answer, err := resp.Resp(resp.CertResponce{Status: cert.Responce, Hash: c.Hash}).Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())

		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "ChangeDialog":
		dr := dialog.DialogRequest{}
		err := json.Unmarshal(contents, &dr)
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		active, err := dr.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при изменении диалога. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изменении диалога. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изменении диалога. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		body, err := dr.ChangeDialog()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := body.Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)

	case "ChangeUser":
		c := RegUser.Request{}
		err := json.Unmarshal(contents, &c)
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		active, err := c.Sessia.IsWithActivateEmail()
		if err != nil {
			log.ChanLog <- fmt.Sprintf("ошибка при изменении пользователя. Подробности: %s", err)
			body := resp.NonAddResponce{
				Status: resp.Responce{}.InternalError("Внутренняя ошибка при проверке пользователя"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изменении пользователя. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}
		if !active {
			body := resp.NonAddResponce{
				Status: resp.Responce{}.BadRequest("Неактивирован Email"),
			}
			answer, err := body.Responser()
			log.ChanLog <- fmt.Sprintf("ошибка при изменении пользователя. Подробности: %s", err)
			w.Write(([]byte(answer)))
			w.WriteHeader(200)
			return
		}

		body, err := c.Do()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := body.Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}

		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	case "disable":

		body, err := Sessia.DisableTocken(contents)
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		answer, err := body.Responser()
		if err != nil {
			log.ChanLog <- fmt.Sprintf(err.Error())
		}
		w.Write(([]byte(answer)))
		w.WriteHeader(200)
	default:
		w.Write([]byte("Некорректный запрос"))
		w.WriteHeader(404)
	}
	//strtobyte := []byte("ответ")
	//w.Write(strtobyte)

	if err != nil {
		//fmt.Println("Внимание, ошибка чтения запроса.\n %s", err)
		http.Error(w, err.Error(), 500)
	}

}
