/*
 * Created on Sat Apr 09 2022
 * Author @LosAngeles971
 *
 * The MIT License (MIT)
 * Copyright (c) 2022 @LosAngeles971
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software
 * and associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
 * TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package business

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type Enigma struct {
	letterRunes []rune
	keyString   [32]byte
	prefix      string
}

type EnigmaOption func(*Enigma)

func withMainkey(email string, password string) EnigmaOption {
	return func(e *Enigma) {
		e.keyString = sha256.Sum256([]byte(email + password + e.prefix))
	}
}

func withRandomkey() EnigmaOption {
	return func(e *Enigma) {
		key := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			panic(err.Error())
		}
	}
}

func withEncodedkey(key string) EnigmaOption {
	return func(e *Enigma) {
		key, _ := hex.DecodeString(key)
		copy(e.keyString[:], key[:32])
	}
}

func newEnigma(opts ...EnigmaOption) *Enigma {
	e := &Enigma{
		letterRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		prefix:      "Kirinuki",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Enigma) hash(data []byte) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func (e *Enigma) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.keyString[:])
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
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(e.keyString[:])
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

func (e *Enigma) getEncodedKey() string {
	return hex.EncodeToString(e.keyString[:])
}