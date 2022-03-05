package storage

import (
	"os"
	"testing"
)

// Test save and load file
func TestFilesysteme(t *testing.T) {
	ll, err := NewFilesystem("test", StorageItem{
		Type: "local",
		Cfg: map[string]string{
			"path": os.TempDir(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	input := "test of filesystem storage"
	err = ll.Put("test1", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	output, err := ll.Get("test1")
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != input {
		t.Fatalf("expected [%s] not [%s]", input, string(output))
	}
}