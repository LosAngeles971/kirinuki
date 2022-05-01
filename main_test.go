package main

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/internal"
)

func cleanup() {
	internal.Clean("tmp")
}

func TestMain(m *testing.M) {
	internal.Setup()
    code := m.Run() 
	cleanup()
    os.Exit(code)
}