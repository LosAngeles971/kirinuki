package internal

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/sirupsen/logrus"
)

const (
	Test_email    = "losangeles971@gmail.com"
	Test_password = "losangeles971@gmail.com"
)

func Setup() {
	logrus.SetLevel(logrus.DebugLevel)
	tmp := os.TempDir() + "/tmp"
	_ = os.Mkdir(tmp, 0755)
}

func GetStorage(tDir string, t *testing.T) *storage.MultiStorage {
	base := os.TempDir() + "/" + tDir
	if _, err := os.Stat(base); err == nil {
		os.RemoveAll(base)
	}
	_ = os.Mkdir(base, os.ModePerm)
	sm, _ := storage.NewMultiStorage()
	sm.AddLocal(tDir, base)
	return sm
}

func Clean(tDir string) {
	os.RemoveAll(os.TempDir() + "/" + tDir)
}

func GetTmp() string {
	return os.TempDir() + "/tmp"
}