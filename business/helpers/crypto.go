package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	key_size    = 32        // size of the symmetric key
	nameSize    = 24        // size of the files' names
	buffer_size = 16 * 1024 // buffer's size during encryption/decryption of files
)

func GetRndEncodedKey() string {
	key := sha256.Sum256(GetRndBytes(key_size))
	return hex.EncodeToString(key[:])
}

// 
func EncryptData(plain []byte, ekey string) ([]byte, error) {
	key, err := hex.DecodeString(ekey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode symmetric key - %v", err)
	}
	cph, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cph)
    if err != nil {
        panic(err)
    }
	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}
	return gcm.Seal(iv, iv, plain, nil), nil
}

// 
func DecryptData(cccc []byte, ekey string) ([]byte, error) {
	key, err := hex.DecodeString(ekey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode symmetric key - %v", err)
	}
	cph, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(cph)
    if err != nil {
        panic(err)
    }
	// Since we know the ciphertext is actually nonce+ciphertext
    // And len(nonce) == NonceSize(). We can separate the two.
    nonceSize := gcm.NonceSize()
    nonce, ccc := cccc[:nonceSize], cccc[nonceSize:]
	return gcm.Open(nil, nonce, ccc, nil)
}

// Local encryption of a file
func EncryptFile(sFile, tFile, ekey string) error {
	log.Debugf("encrypting ( %s )  to ( %s )...", sFile, tFile)
	key, err := hex.DecodeString(ekey)
	if err != nil {
		return fmt.Errorf("failed to decode symmetric key - %v", err)
	}
	in, err := os.Open(sFile)
	if err != nil {
		return fmt.Errorf("failed to open file ( %s ) - %v", sFile, err)
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}
	stream := cipher.NewCTR(block, iv)
	inBuf := make([]byte, buffer_size)
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

// Local decryption of a file
func DecryptFile(sFile, tFile, ekey string) error {
	log.Debugf("decrypting %s to %s ...", sFile, tFile)
	key, err := hex.DecodeString(ekey)
	if err != nil {
		return fmt.Errorf("failed to decode symmetric key - %v", err)
	}
	in, err := os.Open(sFile)
	if err != nil {
		return fmt.Errorf("failed to open file ( %s ) - %v", sFile, err)
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(key)
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
	inBuf := make([]byte, buffer_size)
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