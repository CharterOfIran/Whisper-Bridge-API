package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
)

// EncryptedPayload ساختار نهایی برای ذخیره در ردیس
type EncryptedPayload struct {
	EncryptedKey  string `json:"k"`
	EncryptedData string `json:"d"`
	IV            string `json:"i"`
}

// EncryptHybrid پیاده‌سازی رمزنگاری ترکیبی طبق نقشه راه
func EncryptHybrid(plainText string, publicKeyPEM string) (*EncryptedPayload, error) {
	// 1. Parse RSA Public Key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	// 2. Generate Random AES-256 Key
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	// 3. Encrypt Data with AES-GCM
	blockAES, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(blockAES)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nil, iv, []byte(plainText), nil)

	// 4. Wrap AES Key with RSA-OAEP
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, aesKey, nil)
	if err != nil {
		return nil, err
	}

	return &EncryptedPayload{
		EncryptedKey:  base64.StdEncoding.EncodeToString(encryptedKey),
		EncryptedData: base64.StdEncoding.EncodeToString(ciphertext),
		IV:            base64.StdEncoding.EncodeToString(iv),
	}, nil
}
