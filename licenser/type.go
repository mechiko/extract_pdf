package licenser

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/sys/windows/registry"
)

// типа параметра
type Param int

// Используем константы для псевдо-enum
const (
	FsrarID Param = iota // Отчет начинается с 0
	OmsID
	CpuID
	MAC
)

// Date до какой даты действует
// Identity идентификатор это фсрар ид омс ид или цпу ид
type LicenseInfo struct {
	Date       string
	Identity   string
	AppVersion string
}

// hex encoded
type BoxKeys struct {
	ServerPub string
	LocalPub  string
	LocalPriv string
}

type Licenser struct {
	License        *LicenseInfo
	Parameter      Param
	ParameterValue string
	Date           time.Time
	box            *BoxKeys
}

var root = registry.CURRENT_USER

// var nevakodKey = `Software\NevaKOD`
var nameModuleKey = `Software\NevaKOD\KorrectKM`
var nameLicenseKey = `license`
var nameClientCertPubKey = `certpub`
var nameClientCertPrivKey = `certpriv`

const DateLayout = "2006-01-02"

var licenseKeyValue = ""

// прописываем в ехе
var serverCertPub = "c77c152fac8bd276d3c4e56a49d31759e7fe941a6b580dae5ccae5420532ef6e"

var boxLocal *BoxKeys

// для OmsID и FsrarID значение передаем здесь
// CpuID и MAC вычисляются из функции по текущему компу
func New(prm Param, prmValue string) (lic *Licenser, err error) {
	if err := initBox(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	lic = &Licenser{
		License:        &LicenseInfo{},
		Parameter:      prm,
		ParameterValue: prmValue,
		Date:           time.Time{},
		box:            boxLocal,
	}
	switch prm {
	case CpuID:
		p, err := getCpuId()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		lic.ParameterValue = p
	case MAC:
		p, err := getMac()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		lic.ParameterValue = p
	}
	if lic.ParameterValue == "" {
		return nil, fmt.Errorf("licenser parameter value empty")
	}
	// переменные box инициализируется в init()
	// если ключ не сгенирировался или как то обнулен
	if lic.box.LocalPub == "" {
		return nil, fmt.Errorf("licenser pub key value empty")
	}
	if lic.box.LocalPriv == "" {
		return nil, fmt.Errorf("licenser priv key value empty")
	}
	// переменная модуля инициализируется в init()
	if licenseKeyValue == "" {
		// вызываем диалог запроса лицензии
		licenseKeyValue, err = lic.startDialog()
		if err != nil {
			return nil, fmt.Errorf("licenser %w", err)
		}
		if err := writeStringValue(root, nameModuleKey, nameLicenseKey, licenseKeyValue); err != nil {
			return nil, fmt.Errorf("licenser %w", err)
		}
	}
	if licenseKeyValue != "" {
		// парсим лицензию из json расшифрованной строки
		jsonApiKey, err := decodeBox(licenseKeyValue)
		if err != nil {
			_ = writeStringValue(root, nameModuleKey, nameLicenseKey, "")
			return nil, fmt.Errorf("wrong api key %w", err)
		}
		if err := json.Unmarshal([]byte(jsonApiKey), lic.License); err != nil {
			_ = writeStringValue(root, nameModuleKey, nameLicenseKey, "")
			return nil, fmt.Errorf("wrong api key %w", err)
		}
		if lic.ParameterValue != lic.License.Identity {
			_ = writeStringValue(root, nameModuleKey, nameLicenseKey, "")
			return nil, fmt.Errorf("wrong api param %s license identity %s", lic.ParameterValue, lic.License.Identity)
		}
		lic.Date, err = time.Parse(DateLayout, lic.License.Date)
		if err != nil {
			_ = writeStringValue(root, nameModuleKey, nameLicenseKey, "")
			return nil, fmt.Errorf("wrong api key date %w", err)
		}
		if lic.Date.Before(time.Now()) {
			_ = writeStringValue(root, nameModuleKey, nameLicenseKey, "")
			return nil, fmt.Errorf("license out od date date %v", lic.Date)
		}
		return lic, nil
	}
	return lic, nil
}

func decodeBox(encoded string) (string, error) {
	encodedByte, err := hex.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	senderPublicKeyDecode, err := hex.DecodeString(boxLocal.ServerPub)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	recipientPrivKeyDecode, err := hex.DecodeString(boxLocal.LocalPriv)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	var decryptNonce [24]byte
	copy(decryptNonce[:], encodedByte[:24])
	decrypted, ok := box.Open(nil, encodedByte[24:], &decryptNonce, (*[32]byte)(senderPublicKeyDecode), (*[32]byte)(recipientPrivKeyDecode))
	if !ok {
		return "", fmt.Errorf("decrypt error")
	}
	return string(decrypted), nil
}
