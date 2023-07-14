package command

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

func generateAESKey(keySize int) (string, error) {
	if keySize != 32 && keySize != 24 && keySize != 16 {
		return "", errors.New("AES Key size must be 16/24/32 to select AES-128, AES-192, or AES-256")
	}

	bs := make([]byte, keySize)

	if _, err := rand.Read(bs); err != nil {
		return "", nil
	}

	key := hex.EncodeToString(bs)

	return key, nil
}

func encryptDataWithAES(key InternalKey, pt string) (string, error) {
	aesKey := key.aesKey
	bspt := []byte(pt)

	bsKey, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(bsKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ct := aesGCM.Seal(nonce, nonce, bspt, nil)
	return fmt.Sprintf("%x", ct), nil
}

func decryptDataWithAES(key InternalKey, ct string) (string, error) {
	bskey, err := hex.DecodeString(key.aesKey)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(bskey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, bsct := enc[:nonceSize], enc[nonceSize:]

	pt, err := aesGCM.Open(nil, nonce, bsct, nil)
	if err != nil {
		return "", err
	}

	return string(pt), nil
}

func encryptData(key InternalKey, ek *EncryptWithTransitKeyRequest) (string, error) {
	ct, err := encryptDataWithAES(key, ek.Plaintext)

	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(ct))
	return fmt.Sprintf("vault:v1:%s", b64), nil
}

func decryptData(key InternalKey, dk *DecryptWithTransitKeyRequest) (string, error) {
	ct_data := strings.Split(dk.Ciphertext, ":")
	b64 := ct_data[len(ct_data)-1]
	ct, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		return "", err
	}

	pt, err := decryptDataWithAES(key, string(ct))

	if err != nil {
		return "", err
	}

	return string(pt), nil
}
