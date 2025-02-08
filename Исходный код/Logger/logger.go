package logger

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"fmt"
	"os"
	"strings"
)

func Initialize() (Logger, error) {
	l := Logger{}
	logsexists, err := Utility.Exists(Configs.Logs())
	if err != nil {
		return l, fmt.Errorf("сбой проверки наличия каталога логов. Подробности: %w", err)
	}
	if !logsexists {
		Utility.CreateFolders(Configs.Logs())
	}
	l.ChanLog = make(chan string)
	f, err := os.OpenFile(Configs.ServerLog(), os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		close(l.ChanLog)
		return Logger{}, fmt.Errorf("ошибка инициализации логгера. Подробности: %w", err)
	}
	defer f.Close()
	go log(l.ChanLog)
	return l, nil
}

func log(c chan string) {
	for msg := range c {
		msg = strings.Replace(msg, "\n", "", -1)
		err := Utility.AppendToFile(Configs.ServerLog(), fmt.Sprintf("\n[%v] – %v", Utility.GetTime(), msg))
		if err != nil {
			fmt.Printf("сбой открытия лог файла. Подробности: %w", err)
		}
		data := <-c
		if data == "exit" {
			close(c)
			break
		}
		if err != nil {
			break
		}
	}
}
