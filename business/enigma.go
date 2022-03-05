package business

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	rr "math/rand"
	"time"
)

type Enigma struct {
	letterRunes []rune
	keyString   string
	prefix      string
}

type EnigmaOption func(*Enigma)

func WithKeystring(keystring string) EnigmaOption {
	return func(e *Enigma) {
		e.keyString = keystring
	}
}

func WithMainkey(email string, password string) EnigmaOption {
	return func(e *Enigma) {
		e.keyString = e.hash([]byte(email + password + e.prefix))
	}
}

func WithPrefix(prefix string) EnigmaOption {
	return func(e *Enigma) {
		e.prefix = prefix
	}
}

func NewEnigma(opts ...EnigmaOption) *Enigma {
	e := &Enigma{
		letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		prefix:      "Kirinuki",
		keyString:   "",
	}
	for _, opt := range opts {
		opt(e)
	}
	if len(e.keyString) == 0 {
		b := make([]rune, 32)
		rr.Seed(time.Now().UnixNano())
		for i := range b {
			b[i] = e.letterRunes[rr.Intn(len(e.letterRunes))]
		}
		e.keyString = string(b)
	}
	return e
}

func (e *Enigma) hash(data []byte) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func (e *Enigma) encrypt(plaintext []byte) ([]byte, error) {
	key, _ := hex.DecodeString(e.keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (e *Enigma) decrypt(enc []byte) ([]byte, error) {
	key, _ := hex.DecodeString(e.keyString)
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}
	return plaintext, nil
}
