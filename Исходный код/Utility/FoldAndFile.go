package Utility

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type TypeFile struct {
	types string
}

func GetTime() string {
	time := time.Now().Format("2006-01-02 15:04:05.00")
	return time
}

func CreateFolders(path string) error {
	return os.MkdirAll(path, 0777)
}

const (
	originalfile = "raw"
	media        = "../Media"
)

func MediaExists(hash string) (bool, error) {
	path := Combine([]string{media, hash})
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), originalfile) {
			return true, nil
		}
	}
	return false, nil
}

func RemoveFolder(path string) error {
	return os.RemoveAll(path)
}

func DeleteFiles(path string) error {
	os.RemoveAll(path)
	return nil
}

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func CreateFile(path string) (*os.File, error) {
	fp := filepath.Dir(path)
	exf, err := Exists(fp)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки путей каталога. Подробности: %s", err)
	}
	if !exf {
		err := os.MkdirAll(fp, 0777)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания папок. Подробности: %s", err)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания файлов. %v", err)
	}

	return f, nil
}

func MoveFile(pathfile string, pathtofile string) error {
	CopyFile(pathfile, pathtofile)
	err := os.Remove(pathfile)
	if err != nil {
		return fmt.Errorf("ошибка удаления исходного файла: %s", err)
	}
	return nil
}

func CopyFile(pathfile string, pathtofile string) error {
	inputFile, err := os.Open(pathfile)
	if err != nil {
		return fmt.Errorf("не найден исходный файл: %s", err)
	}
	fp := filepath.Dir(pathtofile)
	exf, err := Exists(fp)
	if err != nil {
		return fmt.Errorf("ошибка проверки путей каталога. Подробности: %s", err)
	}
	if !exf {
		err := os.MkdirAll(fp, 0777)
		if err != nil {
			return fmt.Errorf("ошибка создания папок. Подробности: %s", err)
		}
	}
	outputFile, err := os.Create(pathtofile)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("невозможно скопировать файл: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("невозможна запись файла: %s", err)
	}
	return nil
}

func ReadFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	result := string(file)
	return result
}

func AppendToFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func RewriteFile(path string, content string) error {
	fp := filepath.Dir(path)
	exf, err := Exists(fp)
	if err != nil {
		return fmt.Errorf("ошибка проверки путей каталога")
	}
	if !exf {
		err := os.MkdirAll(fp, 0777)
		if err != nil {
			return fmt.Errorf("ошибка создания папок")
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func Combine(paths []string) string {
	pathes := paths[0]
	for _, item := range paths[1:] {
		pathes = path.Join(pathes, item)
	}
	return pathes
}

