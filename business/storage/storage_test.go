package storage

import (
	"os"
	"testing"
)

func TestLocal(t *testing.T) {
	ll, err := newStorage("test", ConfigItem{
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