package internal

import (
	"os"

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

func Clean(tDir string) {
	os.RemoveAll(os.TempDir() + "/" + tDir)
}

func GetTmp() string {
	return os.TempDir() + "/tmp"
}