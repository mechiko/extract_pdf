package licenser

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.design/x/clipboard"
	"golang.org/x/crypto/nacl/box"
)

// при запуске программы проверяем наличие и работоспособность реестра
// иначе не будет работать потом
// генерируем ключевую пару
func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	boxLocal = &BoxKeys{
		ServerPub: serverCertPub,
	}
	if keyExists(root, nameModuleKey) {
		if valueExists(root, nameModuleKey, nameLicenseKey) {
			licenseKeyValue, _ = readStringValueWithDefault(root, nameModuleKey, nameLicenseKey, "")
		}
		if valueExists(root, nameModuleKey, nameClientCertPubKey) {
			boxLocal.LocalPub, _ = readStringValueWithDefault(root, nameModuleKey, nameClientCertPubKey, "")
		}
		if valueExists(root, nameModuleKey, nameClientCertPrivKey) {
			boxLocal.LocalPriv, _ = readStringValueWithDefault(root, nameModuleKey, nameClientCertPrivKey, "")
		}
	} else {
		key, err := createKey(root, nameModuleKey)
		if err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameLicenseKey, ""); err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameClientCertPrivKey, ""); err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameClientCertPubKey, ""); err != nil {
			panic(err)
		}
	}
	if boxLocal.LocalPriv == "" || boxLocal.LocalPub == "" {
		// создаем ключи и прописываем
		publicKey, privateKey, err := box.GenerateKey(crypto_rand.Reader)
		if err != nil {
			panic(err)
		}
		boxLocal.LocalPub = hex.EncodeToString(publicKey[:])
		boxLocal.LocalPriv = hex.EncodeToString(privateKey[:])
		if err := writeStringValue(root, nameModuleKey, nameClientCertPrivKey, boxLocal.LocalPriv); err != nil {
			panic(err)
		}
		if err := writeStringValue(root, nameModuleKey, nameClientCertPubKey, boxLocal.LocalPub); err != nil {
			panic(err)
		}
	}
}

func initBox() (err error) {
	boxLocal = &BoxKeys{
		ServerPub: serverCertPub,
	}
	if keyExists(root, nameModuleKey) {
		if valueExists(root, nameModuleKey, nameLicenseKey) {
			licenseKeyValue, err = readStringValueWithDefault(root, nameModuleKey, nameLicenseKey, "")
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		if valueExists(root, nameModuleKey, nameClientCertPubKey) {
			boxLocal.LocalPub, err = readStringValueWithDefault(root, nameModuleKey, nameClientCertPubKey, "")
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		if valueExists(root, nameModuleKey, nameClientCertPrivKey) {
			boxLocal.LocalPriv, err = readStringValueWithDefault(root, nameModuleKey, nameClientCertPrivKey, "")
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	} else {
		key, err := createKey(root, nameModuleKey)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := key.SetStringValue(nameLicenseKey, ""); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := key.SetStringValue(nameClientCertPrivKey, ""); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := key.SetStringValue(nameClientCertPubKey, ""); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	if boxLocal.LocalPriv == "" || boxLocal.LocalPub == "" {
		// создаем ключи и прописываем
		publicKey, privateKey, err := box.GenerateKey(crypto_rand.Reader)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		boxLocal.LocalPub = hex.EncodeToString(publicKey[:])
		boxLocal.LocalPriv = hex.EncodeToString(privateKey[:])
		if err := writeStringValue(root, nameModuleKey, nameClientCertPrivKey, boxLocal.LocalPriv); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := writeStringValue(root, nameModuleKey, nameClientCertPubKey, boxLocal.LocalPub); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}
