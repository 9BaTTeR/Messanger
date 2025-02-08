package userservice

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type HeaderImg struct {
	header []byte
}

func (hi HeaderImg) JPG() {
	hi.header = []byte{255, 216, 255}
}
func (hi HeaderImg) PNG() {
	hi.header = []byte{137, 80, 78}
}
func (hi HeaderImg) GIF() {
	hi.header = []byte{71, 73, 70}
}
func (hi HeaderImg) HeadFile() string {
	return string(hi.header)
}

type SizeOf uint64

const (
	Byte  SizeOf = 1
	KByte SizeOf = 1024
	Mbyte SizeOf = 1024
	GByte SizeOf = 1024
)

type typeImg uint64

const (
	varmaxweigth uint = 5242880
)

type Image struct {
	BinImage []byte
}

const (
	JPG typeImg = 1 + iota
	PNG
	GIF
	Undefined typeImg = 0
)

func StringImage(image string) (result []byte) {
	temp := strings.Replace(image, "[", "", 1)
	temp = strings.Replace(temp, "]", "", 1)
	massive := strings.Split(temp, ",")
	for _, k := range massive {
		res, _ := strconv.ParseUint(strings.Replace(k, "\"", "", 2), 10, 8)
		result = append(result, byte(res))
	}
	return result
}

func CreateMedia(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
func CheckImg(file []byte) (error, bool) {
	var tp typeImg
	tp.imgType(file)

	if tp == Undefined {
		return fmt.Errorf("неизвестный формат файла"), false
	}
	if uint(len(file)) > varmaxweigth {
		return fmt.Errorf("максимальный размер файла 5 МБ"), false
	}
	return nil, true
}

func (ti *typeImg) imgType(header []byte) {
	head := HeaderImg{}
	head.GIF()
	if bytes.Equal(header[:3], []byte(head.HeadFile())) {
		*ti = JPG
		return
	}
	head.PNG()
	if bytes.Equal(header[:3], []byte(head.HeadFile())) {
		*ti = PNG
		return
	}
	head.GIF()
	if bytes.Equal(header[:3], []byte(head.HeadFile())) {
		*ti = GIF
		return
	}
	*ti = Undefined
}

func GetMd5(bytes []byte) string {
	hash := md5.New()
	hash.Write([]byte(bytes))
	return hex.EncodeToString(hash.Sum(nil))
}

func Copy(source string, target string) error {

	if source == target {
		return nil
	}
	fin, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла. Подробности: " + err.Error())
	}
	defer fin.Close()

	fout, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("ошибка создания файла. Подробности: " + err.Error())
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return fmt.Errorf("ошибка копирования файла. Подробности: " + err.Error())
	}
	return nil
}

func Move(source string, target string) {
	os.Rename(source, target)
}