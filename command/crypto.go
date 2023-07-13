package command

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	// "crypto/sha256"
	// "crypto/x509"
	"encoding/hex"
	// "encoding/pem"
	"errors"
	"fmt"
	"io"
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

func generateRSAKey() (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return &rsa.PrivateKey{}, errors.New("unable to init key")
	}

	return priv, nil
}

// func privateKeyToBytes(priv *rsa.PrivateKey) []byte {
// 	privBytes := pem.EncodeToMemory(
// 		&pem.Block{
// 			Type:  "RSA PRIVATE KEY",
// 			Bytes: x509.MarshalPKCS1PrivateKey(priv),
// 		},
// 	)
//
// 	return privBytes
// }
//
// func publicKeyToBytes(pub *rsa.PublicKey) []byte {
// 	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
// 	if err != nil {
// 		errors.New("Unable convert public key to []byte")
// 	}
//
// 	pubBytes := pem.EncodeToMemory(&pem.Block{
// 		Type:  "RSA PUBLIC KEY",
// 		Bytes: pubASN1,
// 	})
//
// 	return pubBytes
// }
//
// // BytesToPrivateKey bytes to private key
// func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
// 	block, _ := pem.Decode(priv)
// 	enc := x509.IsEncryptedPEMBlock(block)
// 	b := block.Bytes
// 	var err error
// 	if enc {
// 		b, err = x509.DecryptPEMBlock(block, nil)
// 		if err != nil {
// 			errors.New("unable to decrypt private pem block")
// 		}
// 	}
// 	key, err := x509.ParsePKCS1PrivateKey(b)
// 	if err != nil {
// 		errors.New("unable to parse private key")
// 	}
// 	return key
// }
//
// // BytesToPublicKey bytes to public key
// func BytesToPublicKey(pub []byte) *rsa.PublicKey {
// 	block, _ := pem.Decode(pub)
// 	enc := x509.IsEncryptedPEMBlock(block)
// 	b := block.Bytes
// 	var err error
// 	if enc {
// 		b, err = x509.DecryptPEMBlock(block, nil)
// 		if err != nil {
// 			errors.New("unable to decrypt public pem block")
// 		}
// 	}
// 	ifc, err := x509.ParsePKIXPublicKey(b)
// 	if err != nil {
// 		errors.New("unable to parse public pem block")
// 	}
// 	key, ok := ifc.(*rsa.PublicKey)
// 	if !ok {
// 		errors.New("unable to public block")
// 	}
// 	return key
// }

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

	return fmt.Sprintf("%s", pt), nil
}

// func encryptDataWithRSA(key InternalKey, pt []byte) ([]byte, error) {
// 	pub := &key.rsaKey.PublicKey
//
// 	hash := sha256.New()
// 	ct, err := rsa.EncryptOAEP(hash, rand.Reader, pub, []byte(pt), nil)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return ct, nil
// }
//
// func decryptDataWithRSA(key InternalKey, ct []byte) ([]byte, error) {
// 	priv := key.rsaKey
//
// 	hash := sha256.New()
// 	pt, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ct, nil)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return pt, nil
// }

// func encryptDataWithRSA(key InternalKey, ek *EncryptWithTransitKeyRequest) (string, error) {
// 	pt := ek.Plaintext
// 	pub := &key.privateKey.PublicKey
//
// 	bs := split_bytes([]byte(pt), 256)
// 	ct := []byte{}
//
// 	for _, v := range bs {
// 		hash := sha256.New()
// 		t, err := rsa.EncryptOAEP(hash, rand.Reader, pub, v, nil)
//
// 		if err != nil {
// 			return "", err
// 		}
//
// 		ct = append(ct, t...)
// 	}
//
// 	b64 := base64.StdEncoding.EncodeToString(ct)
// 	return b64, nil
// }
//
// func decryptDataWithRSA(key InternalKey, dk *DecryptWithTransitKeyRequest) (string, error) {
// 	ct_data := strings.Split(dk.Ciphertext, ":")
// 	b64 := ct_data[len(ct_data)-1]
// 	ct, err := base64.StdEncoding.DecodeString(b64)
//
// 	if err != nil {
// 		return "", errors.New("unable to decode base64 string")
// 	}
//
// 	priv := key.privateKey
// 	bs := split_bytes([]byte(ct), 256)
// 	pt := []byte{}
//
// 	for _, v := range bs {
// 		hash := sha256.New()
// 		t, err := rsa.DecryptOAEP(hash, rand.Reader, priv, v, nil)
//
// 		if err != nil {
// 			return "", err
// 		}
// 		pt = append(pt, t...)
// 	}
//
// 	return string(pt), nil
// }
