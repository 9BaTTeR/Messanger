package media

import (
	"ServerApp/Configs"
	responces "ServerApp/Responces"
	"ServerApp/Utility"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func (m *Media) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &m)
	if err != nil {
		return err
	}
	return nil
}
func (m Media) Compose() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 300 600 900 1080 :jpeg
const (
	Byte          uint = 8
	KByte         uint = 1024
	Mbyte         uint = 1024 * KByte
	GByte         uint = 1024 * Mbyte
	varmaxweigth  uint = 10 * Mbyte
	dataerror          = "Размер превышен"
	internalerror      = "Внутренняя ошибка"
	succfull           = "Медиа успешно создано"
)

func (m Media) createFolderMedia() (string, error) {
	pathMedia := Configs.Media()
	exists, err := Utility.Exists(pathMedia)

	if !exists {
		err = Utility.CreateFolders(pathMedia)
	}
	if err != nil {
		return "", fmt.Errorf("ошибка чтения каталога media. Подробности: %w", err)
	}

	if m.MediaKey == "" {
		return "", fmt.Errorf("mediaKey был пустой")
	}

	exists, err = Utility.Exists(Utility.Combine([]string{pathMedia, m.MediaKey}))
	if !exists {
		err = Utility.CreateFolders(Utility.Combine([]string{pathMedia, m.MediaKey}))
	}
	if err != nil {
		return "", fmt.Errorf("ошибка чтения каталога media. Подробности: %w", err)
	}
	return Utility.Combine([]string{pathMedia, m.MediaKey, "raw." + m.Extension}), nil
}

func (m *Media) CreateMedia() (responces.MediaResponce, error) {
	decodebytes, err := Utility.DecodeBASE64(m.BytesBase64)
	if err != nil {
		return responces.MediaResponce{Status: (responces.Responce{}).InternalError(internalerror), Key: "-1"}, err
	}
	m.Bytes = decodebytes

	err = CheckImg([]byte(m.Bytes))
	if err != nil {
		return responces.MediaResponce{Status: (responces.Responce{}).InternalError(dataerror), Key: "-1"}, err
	}

	m.MediaKey = Utility.MD5Hash(m.BytesBase64 + m.Extension)
	path, err := m.createFolderMedia()
	if err != nil {
		return responces.MediaResponce{Status: (responces.Responce{}).InternalError(internalerror), Key: "-1"}, err
	}

	f, err := os.Create(path)
	if err != nil {
		return responces.MediaResponce{Status: (responces.Responce{}).InternalError(internalerror), Key: "-1"}, err
	}
	_, err = f.Write(m.Bytes)
	if err != nil {
		return responces.MediaResponce{Status: (responces.Responce{}).InternalError(internalerror), Key: "-1"}, err
	}
	f.Close()
	return responces.MediaResponce{Status: (responces.Responce{}).OK("Медиа загружены"), Key: m.MediaKey}, err
}

func CheckImg(file []byte) error {
	if uint(len(file)) > varmaxweigth {
		return fmt.Errorf("максимальный размер файла 10 МБ")
	}
	return nil
}

func (m *Media) Resize1080p() error {
	image, _, err := image.Decode(bytes.NewReader(m.Bytes))
	if err != nil {
		return err
	}
	newImage := resize.Resize(1080, 0, image, resize.Lanczos3)
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return err
	}
	m.Bytes = buff.Bytes()
	return nil
}

func (m *Media) Resize300p() error {
	image, _, err := image.Decode(bytes.NewReader(m.Bytes))
	if err != nil {
		return err
	}
	newImage := resize.Resize(300, 0, image, resize.Lanczos3)
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return err
	}
	m.Bytes = buff.Bytes()
	return nil
}

func (m *Media) Resize600p() error {
	image, _, err := image.Decode(bytes.NewReader(m.Bytes))
	if err != nil {
		return err
	}
	newImage := resize.Resize(600, 0, image, resize.Lanczos3)
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return err
	}
	m.Bytes = buff.Bytes()
	return nil
}

func (m *Media) Resize900p() error {
	image, _, err := image.Decode(bytes.NewReader(m.Bytes))
	if err != nil {
		return err
	}
	newImage := resize.Resize(900, 0, image, resize.Lanczos3)
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return err
	}
	m.Bytes = buff.Bytes()
	return nil
}

