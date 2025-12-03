package licenser

import (
	crypto_rand "crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/nacl/box"
)

func init() {
	boxLocal = &BoxKeys{
		ServerPub: serverCertPub,
	}
	if KeyExists(root, nameModuleKey) {
		if ValueExists(root, nameModuleKey, nameLicenseKey) {
			licenseKeyValue, _ = ReadStringValueWithDefault(root, nameModuleKey, nameLicenseKey, "")
		}
		if ValueExists(root, nameModuleKey, nameClientCertPubKey) {
			boxLocal.LocalPub, _ = ReadStringValueWithDefault(root, nameModuleKey, nameClientCertPubKey, "")
		}
		if ValueExists(root, nameModuleKey, nameClientCertPrivKey) {
			boxLocal.LocalPriv, _ = ReadStringValueWithDefault(root, nameModuleKey, nameClientCertPrivKey, "")
		}
	} else {
		key, err := CreateKey(root, nameModuleKey)
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
		if err := WriteStringValue(root, nameModuleKey, nameClientCertPrivKey, boxLocal.LocalPriv); err != nil {
			panic(err)
		}
		if err := WriteStringValue(root, nameModuleKey, nameClientCertPubKey, boxLocal.LocalPub); err != nil {
			panic(err)
		}
	}
}
