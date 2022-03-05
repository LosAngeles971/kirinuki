package business

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestTOC(t *testing.T) {
	sm, err := storage.NewStorageMap()
	if err != nil {
		t.Fatal(err)
	}
	sm.Add("test", storage.StorageItem{
		Type: "filesystem",
		Cfg: map[string]string{
			"path": os.TempDir(),
		},
	})
	session, err := NewSession(test_email, test_password, true, WithStorage(sm))
	if err != nil {
		t.Fatalf("failed to create a session due to %v", err)
	}
	err = session.login()
	if err != nil {
		t.Fatalf("failed to open a session due to %v", err)
	}
	toc, err := session.getTOC()
	if err != nil {
		t.Fatalf("failed to get toc from session due to %v", err)
	}
	for _, tt := range k_data_tests {
		k, err := NewKirinuki(WithKirinukiData(tt.name, tt.data))
		if err != nil {
			t.Fatal(err)
		}
		ok := toc.Add(k)
		if !ok {
			t.Fatal("File not added")
		}
		if !toc.Find(tt.name) {
			t.Fatalf("toc does not contain kirinuki %s", tt.name)
		}
	}
	err = session.logout()
	if err != nil {
		t.Fatalf("failed to logout from session due to %v", err)
	}
	err = session.login()
	if err != nil {
		t.Fatalf("failed to open a session due to %v", err)
	}
	toc2, err := session.getTOC()
	if err != nil {
		t.Fatalf("failed to get toc (2) from session due to %v", err)
	}
	for _, tt := range k_data_tests {
		if !toc2.Find(tt.name) {
			t.Fatalf("reloaded toc does not contain kirinuki %s", tt.name)
		}
	}
}