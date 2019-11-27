package cipherx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type Cipherx struct {
	gcm cipher.AEAD
}

func NewCipherx(secretKey string) (*Cipherx, error) {
	passphrase := createHash(secretKey)
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipherx{gcm}, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *Cipherx) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := c.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (c *Cipherx) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
