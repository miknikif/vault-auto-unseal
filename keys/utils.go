package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/miknikif/vault-auto-unseal/common"
)

const (
	AES_KEY_SIZE_128 = 16
	AES_KEY_SIZE_192 = 24
	AES_KEY_SIZE_256 = 32
)

type AESKey string

func generateAESKey(keySize int) (AESKey, error) {
	var key AESKey
	if keySize != AES_KEY_SIZE_128 && keySize != AES_KEY_SIZE_192 && keySize != AES_KEY_SIZE_256 {
		return "", errors.New("AES Key size must be 16/24/32 to select AES-128, AES-192, or AES-256")
	}

	bs := make([]byte, keySize)

	if _, err := rand.Read(bs); err != nil {
		return "", nil
	}

	key = AESKey(hex.EncodeToString(bs))

	return key, nil
}

func encryptDataWithAES(key AESKeyModel, aesPayload *AESPayload) error {
	l, _ := common.GetLogger()
	l.Debug("encryptDataWithAES - started", "key", key, "payload", aesPayload)
	if err := aesPayload.validatePlaintext(); err != nil {
		return err
	}

	bspt := []byte(aesPayload.Plaintext)

	bsKey, err := hex.DecodeString(string(key.AESKey))
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(bsKey)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ct := aesGCM.Seal(nonce, nonce, bspt, nil)

	aesPayload.Payload = common.EncToB64(hex.EncodeToString(ct))
	aesPayload.Pref = "vault"
	aesPayload.Version = key.Version

	if err := aesPayload.validateCiphertext(); err != nil {
		return err
	}

	l.Debug("encryptDataWithAES - finished", "key", key, "payload", aesPayload)
	return nil
}

func decryptDataWithAES(key AESKeyModel, aesPayload *AESPayload) error {
	l, _ := common.GetLogger()
	l.Debug("decryptDataWithAES - started", "key", key, "payload", aesPayload)
	if err := aesPayload.validateCiphertext(); err != nil {
		return err
	}

	bskey, err := hex.DecodeString(string(key.AESKey))
	if err != nil {
		return err
	}
	payload, err := common.DecFromB64(aesPayload.Payload)
	if err != nil {
		return err
	}
	enc, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(bskey)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, bsct := enc[:nonceSize], enc[nonceSize:]

	pt, err := aesGCM.Open(nil, nonce, bsct, nil)
	if err != nil {
		return err
	}

	aesPayload.Plaintext = string(pt[:])
	aesPayload.Pref = "vault"
	aesPayload.Version = key.Version

	if err := aesPayload.validatePlaintext(); err != nil {
		return err
	}

	l.Debug("decryptDataWithAES - finished", "key", key, "payload", aesPayload)
	return nil
}

func createNewAESKeyModel(ver int, size int) (AESKeyModel, error) {
	key, err := generateAESKey(size)
	if err != nil {
		return AESKeyModel{}, err
	}
	return AESKeyModel{
		Name:    int(time.Now().Unix()),
		Version: ver,
		AESKey:  key,
	}, nil
}

func findKeyVersion(keys []AESKeyModel, version int) (AESKeyModel, error) {
	for _, key := range keys {
		if key.Version == version {
			return key, nil
		}
	}
	return AESKeyModel{}, errors.New("Specified version of the key not found")
}

// func encryptData(key InternalKey, ek *EncryptWithTransitKeyRequest) (string, error) {
// 	ct, err := encryptDataWithAES(key, ek.Plaintext)
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	b64 := base64.StdEncoding.EncodeToString([]byte(ct))
// 	return fmt.Sprintf("vault:v1:%s", b64), nil
// }
//
// func decryptData(key InternalKey, dk *DecryptWithTransitKeyRequest) (string, error) {
// 	ct_data := strings.Split(dk.Ciphertext, ":")
// 	b64 := ct_data[len(ct_data)-1]
// 	ct, err := base64.StdEncoding.DecodeString(b64)
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	pt, err := decryptDataWithAES(key, string(ct))
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return string(pt), nil
// }
