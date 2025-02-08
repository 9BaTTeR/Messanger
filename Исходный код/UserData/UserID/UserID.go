package userID

import (
	"ServerApp/Configs"
	"ServerApp/Utility"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Storage struct {
	MDir uint64
	SDir uint16
}

func (st *Storage) Set(code uint64) {
	sts := ConvertionID(code)
	st.MDir, st.SDir = sts.MDir, sts.SDir
}

func (st Storage) MasterDir() string {
	return strconv.FormatUint(uint64(st.MDir), 10)
}

func (st Storage) SlaveDir() string {
	return strconv.FormatUint(uint64(st.SDir), 10)
}

var id Storage = Storage{}

const (
	_uint16 = 65535
)

func Initizalize() error {
	return nil
}

func LastFreeID() Storage {
	storage := Storage{}
	if id != storage {
		storage = id
		storage.IncrementID()
		id = storage
		return storage
	}
	if !Configs.Shutdown() {
		storage = emergencyWay()
		storage.IncrementID()
		return storage
	}
	storage, err := normalWay()
	if os.IsNotExist(err) {
		storage = emergencyWay()
	}
	storage.IncrementID()
	if os.IsNotExist(err) {
		writeID(storage.IDConversion())
	}
	id = storage
	return storage
}

func (st *Storage) DecrementID() {
	tempId := st.IDConversion()
	if tempId > 0 {
		tempst := ConvertionID(uint64(tempId - 1))
		st.MDir, st.SDir = tempst.MDir, tempst.SDir
		id.MDir, id.SDir = st.MDir, st.SDir
	}

}

func (st *Storage) IncrementID() {
	tempId := st.IDConversion()
	tempst := ConvertionID(uint64(tempId + 1))
	st.MDir, st.SDir = tempst.MDir, tempst.SDir
	id.MDir, id.SDir = st.MDir, st.SDir
	if tempId%10 == 0 {
		writeID(tempId)
	}

}

func writeID(id uint64) error {
	return Utility.RewriteFile(Configs.LastID(), strconv.FormatUint(uint64(id), 10))
}

func WriteID() error {
	if id.IDConversion() == 0 {
		return nil
	}
	return writeID(id.IDConversion())

}

func normalWay() (Storage, error) {
	path := Configs.LastID()
	result, err := os.ReadFile(path)
	if err != nil {
		return Storage{}, err
	}
	if result == nil {
		return Storage{}, err
	}

	id, err := strconv.ParseUint(strings.Replace(string(result), "\n", "", 1), 10, 32)

	if err != nil {
		return Storage{}, err
	}
	storage := ConvertionID(uint64(id))
	return storage, nil
}
func emergencyWay() Storage {
	path := Configs.UserDir()
	MDir, err := getMDir(path)
	if err != nil {
		return Storage{MDir: 0, SDir: 0}
	}
	SDir, err := getSDir(path, MDir)
	if err != nil {
		return Storage{MDir: 0, SDir: 0}
	}
	return Storage{MDir: MDir, SDir: SDir}
}

func ConvertionID(id uint64) Storage {
	var dir uint64 = 0
	var file uint16 = 0
	dir = id / uint64(_uint16)
	file = uint16((float32(id)/float32(_uint16) - float32(dir)) * _uint16)
	storage := Storage{MDir: dir, SDir: file}
	return storage
}

func (st Storage) IDConversion() uint64 {
	return (st.MDir * _uint16) + uint64(st.SDir)
}

func getMDir(pathToDir string) (uint64, error) {
	f, err := os.Open(pathToDir)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	list, _ := f.Readdirnames(-1)
	max := 0
	for _, f := range list {
		temp := f
		num, _ := strconv.Atoi(temp)
		if num > max {
			max = num
		}
	}
	return uint64(max), nil
}
func getSDir(pathToDir string, Dir uint64) (uint16, error) {
	f, err := os.Open(pathToDir + "/" + strconv.Itoa(int(Dir)))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	list, err := f.Readdirnames(-1)
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, fmt.Errorf("ошибка, файлы каталога пусты: ")
	}
	max := 0
	for _, f := range list {

		temp := f
		num, _ := strconv.Atoi(temp)
		if num > max {
			max = num
		}
	}
	return uint16(max), nil
}

func (st Storage) PathDbUserId() string {
	path := Configs.UserDir()
	pathUserFolder := fmt.Sprintf("%s/%s/%s/", path, strconv.FormatUint(uint64(st.MDir), 10), strconv.FormatUint(uint64(st.SDir), 10))
	return pathUserFolder
}
