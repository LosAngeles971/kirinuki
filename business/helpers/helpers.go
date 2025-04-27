package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func DeleteLocalFile(filename string) {
	if stat, err := os.Stat(filename); err == nil {
		if !stat.IsDir() {
			os.Remove(filename)
		}
	}
}

func GetRndBytes(size int) []byte {
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}
	return key
}

func GetFilename(size int) string {
	dd := GetRndBytes(size)
	return hex.EncodeToString(dd)
}

// It returns the hash of a array of bytes
func GetHash(data []byte) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// It returns the hash of a local file
func GetFileHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil || n == 0 {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func CreateRandomFile(sFile string, size int) error {
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return err
	}
	return os.WriteFile(sFile, data, 0755)
}