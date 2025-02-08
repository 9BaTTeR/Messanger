package Utility

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func MD5BHash(text []byte) []byte {
	hash := md5.Sum(text)
	return hash[:]
}

func SHA256(text string) string {
	hash := sha3.New256()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum([]byte(text))[:])
}

func SHA384(text string) string {
	hash := sha3.New384()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum([]byte(text))[:])
}

func SHA512(text string) string {
	hash := sha3.New512()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func SaltPass(pass string) (result string) {
	result = SHA512(pass)
	return
}

func BASE64(bytes []byte) string{
	return b64.StdEncoding.EncodeToString(bytes)
}

func DecodeBASE64(base64 string) ([]byte, error){
	result, err := b64.StdEncoding.DecodeString(base64)
	if err != nil{
		return result, fmt.Errorf("ошибка кодирования в BASE64. Подробности: %+v", err)
	}
	return result, nil
}
