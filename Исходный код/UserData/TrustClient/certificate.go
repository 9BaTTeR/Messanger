package trustclient

import (
	"ServerApp/Configs"
	rr "ServerApp/Responces"
	"ServerApp/Utility"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

func (c *Certificate) Parse(source string) error {
	err := json.Unmarshal([]byte(source), &c)
	if err != nil {
		return err
	}
	return nil
}
func (c Certificate) Compose() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}

const (
	lifetime = time.Millisecond * 50
)

const (
	insertRow = "INSERT INTO CERTS\n" +
		"VALUES(?,?,?,?,?,?)"
	lastRowID    = "select rowid from CERTS ORDER BY rowid DESC limit 1;"
	rowDataCerts = "select HashCert,MinVersion,MaxVersion, NameClient FROM CERTS WHERE NameClient like ? LIMIT 1;"
)

func (c *Certificate) Append() (rr.VerifedResponce, error) {
	c.verifed = true
	pathDB := Configs.CertsDB()
	exists, err := Utility.Exists(pathDB)
	vr := rr.VerifedResponce{}
	vr.Responce = vr.Responce.InternalError("сбой регистрации сертификата")
	if err != nil {

		return vr, fmt.Errorf("невозможно проверить наличие БД. Подробности: %w", err)
	}
	if !exists {
		err := createDB()
		if err != nil {
			return vr, fmt.Errorf("сбой генерации БД. Подробности: %w", err)
		}
	}
	exists, needupd, err := c.alreadyHave()
	if err != nil {
		return vr, fmt.Errorf("ошибка проверка сертификата на наличие. Подробности: %w", err)
	}
	if needupd {
		vr.Responce = vr.Responce.BadRequest("Обновите версию для сертификата.")
		return vr, nil
	}
	if exists {
		vr.Responce = vr.Responce.BadRequest("У вас уже есть сертификат для приложения с таким именем и версией.")
		return vr, nil
	}
	err = c.lastID()
	if err != nil {
		return vr, fmt.Errorf("сбой получения последней записи ID. Подробности: %w", err)
	}
	err = c.Certificate()
	if err != nil {
		return vr, fmt.Errorf("сбой генерации сертификата. Подробности: %w", err)
	}
	sql, err := openDB()
	sql.SetConnMaxLifetime(lifetime)
	if err != nil {
		sql.Close()
		return vr, fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	_, err = sql.Exec(insertRow, c.Hash, c.Cert, c.Name, c.MinVersion, c.MinVersion, c.verifed)
	if err != nil {
		return vr, fmt.Errorf("ошибка записи сертификата. Подробности: %w", err)
	}
	vr.Responce = rr.Responce{}.OK("Действие выполнено, сертификат записан.")
	return vr, nil
}

// Первый bool - сертификат есть.
// Второй bool - сертификат требует обновления версии.
func (c *Certificate) alreadyHave() (bool, bool, error) {
	rawVersion := c.MinVersion
	sql, err := openDB()
	sql.SetConnMaxLifetime(lifetime)
	if err != nil {
		return false, false, fmt.Errorf("ошибка соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(rowDataCerts, c.Name)
	if err != nil {
		return false, false, fmt.Errorf("ошибка чтения записей из БД. Подробности: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return false, false, nil
	}
	err = rows.Scan(&c.Hash, &c.MinVersion, &c.MaxVersion, &c.Name)
	rows.Close()
	if err != nil {
		return false, false, fmt.Errorf("ошибка парсинга значения. Подробности: %w", err)
	}
	defer rows.Close()
	if rawVersion > c.MaxVersion {
		return true, true, nil
	}
	if rawVersion < c.MinVersion {
		return true, true, nil
	}
	return true, false, nil
}

func (c *Certificate) lastID() error {
	if c.id != 0 {
		return fmt.Errorf("индекс уже определён")
	}
	sql, err := openDB()
	sql.SetConnMaxLifetime(lifetime)
	if err != nil {
		return fmt.Errorf("сбой соединения с БД. Подробности: %w", err)
	}
	defer sql.Close()
	rows, err := sql.Query(lastRowID)
	if err != nil {
		return fmt.Errorf("сбой чтения БД. Подробности: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		c.id = 1
		return nil
	}
	err = rows.Scan(&c.id)
	if err != nil {
		return fmt.Errorf("ошибка парсера записи из БД. Подробности: %w", err)
	}
	c.id++
	return nil
}

func genRSA(r *rsa.PrivateKey) error {
	if len(r.Primes) > 0 {
		return fmt.Errorf("уже определён. Опередёлнный ключ %+v", r)
	}
	temp := rand.Reader
	r, err := rsa.GenerateKey(temp, 512)
	if err != nil {
		return fmt.Errorf("ошибка генерации RSA ключа. Подробности: %s.", err)
	}
	brsa, err := x509.MarshalPKCS8PrivateKey(r)
	if err != nil {
		return fmt.Errorf("ошибка сериализации RSA ключа. Подробности: %w", err)
	}
	file, err := Utility.CreateFile(Configs.InternalRSA())
	if err != nil {
		return fmt.Errorf("ошибка создание файла ключа. Подробности: %w", err)
	}

	_, err = file.Write(brsa)
	if err != nil {
		return fmt.Errorf("ошибка записи RSA ключа. Подробности: %w", err)
	}
	defer file.Close()
	return nil
}

func loadRSA() (*rsa.PrivateKey, error) {
	rsaresult := &rsa.PrivateKey{}
	exists, err := Utility.Exists(Configs.InternalRSA())
	if err != nil {
		return nil, fmt.Errorf("невозможно проверить наличие ключа. Подробности: %w", err)
	}
	if !exists {
		err = genRSA(rsaresult)
		if err != nil {
			return nil, fmt.Errorf("ошибка при генерации RSA ключа. Подробности: %w", err)
		}
		return rsaresult, nil
	}
	rsatext := []byte(Utility.ReadFile(Configs.InternalRSA()))
	temprsa, err := x509.ParsePKCS8PrivateKey(rsatext)
	if err != nil {
		return rsaresult, fmt.Errorf("ошибка чтения RSA ключа из файла. Подробности: %w", err)
	}
	rsaresult = temprsa.(*rsa.PrivateKey)
	return rsaresult, nil

}

func (cf Certificate) template() *x509.Certificate {
	template := x509.Certificate{
		Issuer: pkix.Name{Organization: []string{"RollaDie"},
			CommonName: "RollaDie"},
		SerialNumber: big.NewInt(int64(cf.id)),
		Subject: pkix.Name{
			Organization: []string{cf.Name},
			CommonName:   cf.Name,
		},
		EmailAddresses:        []string{cf.Email},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 6, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		PublicKeyAlgorithm:    x509.RSA,
	}
	return &template
}

func (cf *Certificate) ParseCert(path string) error {
	certtext := Utility.ReadFile("test.pem")
	cert, err := x509.ParseCertificate([]byte(certtext))
	if err != nil {
		return fmt.Errorf("ошибка чтения сертификата. Подробности: %w", err)
	}
	cert.Verify(x509.VerifyOptions{})
	return nil
}

func (cf *Certificate) Certificate() error {
	r := &rsa.PrivateKey{}
	r, err := loadRSA()
	if err != nil {
		return fmt.Errorf("ошибка загрузки RSA ключа. Подробности: %w", err)
	}
	cert, err := x509.CreateCertificate(rand.Reader, cf.template(), cf.template(), &r.PublicKey, r)
	if err != nil {
		return fmt.Errorf("ошибка генерации сертификата. Подробности: %w", err)
	}
	if err != nil {
		return fmt.Errorf("ошибка генерации сертификата. Подробности: %w", err)
	}
	cf.Cert = cert
	cf.Hash = hex.EncodeToString(Utility.MD5BHash(cf.Cert))
	return nil
}
