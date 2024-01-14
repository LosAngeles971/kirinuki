package main

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestMain(m *testing.M) {
	tmp := os.TempDir() + "/kirinuki_tmp"
	os.Setenv(storage.KIRINUKI_TMP, tmp)
	_ = os.Mkdir(tmp, os.ModePerm)
    code := m.Run() 
	os.RemoveAll(os.TempDir() + "/kirinuki_tmp")
    os.Exit(code)
}