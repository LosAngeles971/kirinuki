package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const (
	KIRINUKI_TMP = "KIRINUKI_TMP"
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

// GetTmp: it returns the local temporary directory
// The local temporary directory is always necessary, 
// since it is used to split e rebuild Kirinuki files.
func GetTmp() string {
	tmp, ok := os.LookupEnv(KIRINUKI_TMP)
	if ok {
		return tmp
	} else {
		return os.TempDir()
	}
}

func GetHash(data []byte) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

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