/*
 * Created on Sun Apr 10 2022
 * Author @LosAngeles971
 *
 * This software is licensed under GNU General Public License v2.0
 * Copyright (c) 2022 @LosAngeles971
 *
 * The GNU GPL is the most widely used free software license and has a strong copyleft requirement.
 * When distributing derived works, the source code of the work must be made available under the same license.
 * There are multiple variants of the GNU GPL, each with different requirements.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
 * TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package enigma

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	key_size = 32
)

//const V1 byte = 0x1

type Enigma struct {
	keyString   [key_size]byte
	prefix      string
	buffer_size int
	iv_size     int
	hmacSize    int
}

type EnigmaOption func(*Enigma)

func GetRndBytes(size int) []byte {
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}
	return key
}

func WithMainkey(email string, password string) EnigmaOption {
	return func(e *Enigma) {
		e.keyString = sha256.Sum256([]byte(email + password + e.prefix))
	}
}

func WithRandomkey() EnigmaOption {
	return func(e *Enigma) {
		e.keyString = sha256.Sum256(GetRndBytes(key_size))
	}
}

func WithEncodedkey(key string) EnigmaOption {
	return func(e *Enigma) {
		key, _ := hex.DecodeString(key)
		copy(e.keyString[:], key[:key_size])
	}
}

func New(opts ...EnigmaOption) *Enigma {
	e := &Enigma{
		prefix:      "Kirinuki",
		buffer_size: 16 * 1024,
		iv_size:     16,
		hmacSize:    sha512.Size,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Enigma) EncryptData(data []byte) ([]byte, error) {
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
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Enigma) DecryptData(enc []byte) ([]byte, error) {
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

func (e *Enigma) GetEncodedKey() string {
	return hex.EncodeToString(e.keyString[:])
}

func (e *Enigma) EncryptFile(sFile string, tFile string) error {
	log.Debugf("encrypting %s to %s ...", sFile, tFile)
	in, err := os.Open(sFile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(e.keyString[:])
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}
	stream := cipher.NewCTR(block, iv)
	inBuf := make([]byte, e.buffer_size)
	for {
		n, err := in.Read(inBuf)
		if err == io.EOF {
			out.Write(iv)
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}
		stream.XORKeyStream(inBuf, inBuf[:n])
		out.Write(inBuf[:n])
	}
}

func (e *Enigma) DecryptFile(sFile string, tFile string) error {
	log.Debugf("decrypting %s to %s ...", sFile, tFile)
	in, err := os.Open(sFile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(e.keyString[:])
	if err != nil {
		return err
	}
	fi, err := in.Stat()
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = in.ReadAt(iv, msgLen)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	inBuf := make([]byte, e.buffer_size)
	for {
		n, err := in.Read(inBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}
		if n > int(msgLen) {
			n = int(msgLen)
		}
		msgLen -= int64(n)
		stream.XORKeyStream(inBuf, inBuf[:n])
		out.Write(inBuf[:n])
	}
}
