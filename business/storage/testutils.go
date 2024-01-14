package storage

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	Test_email    = "losangeles971@gmail.com"
	Test_password = "losangeles971@gmail.com"
)

type TestLocalMultistorage struct {
	sm *MultiStorage
	base string
}

func SetTestEnv() {
	tmp := os.TempDir() + "/kirinuki_tmp"
	os.Setenv(KIRINUKI_TMP, tmp)
	_ = os.Mkdir(tmp, os.ModePerm)
}

func CleanTestEnv() {
	os.RemoveAll(os.TempDir() + "/kirinuki_tmp")
}

// NewLocalMultistorage: it creates a multistorage for testing purpose.
// This multistorage only uses local filesystem by means of a single Local Storage Target
func NewTestLocalMultistorage(tDir string) *TestLocalMultistorage {
	sm, _ := NewMultiStorage()
	tsm := &TestLocalMultistorage{
		sm: sm,
		base: os.TempDir() + "/" + tDir,
	}
	logrus.SetLevel(logrus.DebugLevel)
	SetTestEnv()
	if _, err := os.Stat(tsm.base); err == nil {
		os.RemoveAll(tsm.base)
	}
	_ = os.Mkdir(tsm.base, os.ModePerm)
	sm.AddLocal(tDir, tsm.base)
	return tsm
}

func (tsm *TestLocalMultistorage) Clean() {
	os.RemoveAll(tsm.base)
}

func (tsm *TestLocalMultistorage) GetMultiStorage() *MultiStorage {
	return tsm.sm
}

func CreateFile(sFile string, size int) error {
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return err
	}
	return ioutil.WriteFile(sFile, data, 0755)
}

