package business

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestGateway(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	g, err := New(test_email, test_password, true, sm)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range k_data_tests {
		err = g.Upload(tt.name, tt.data, true)
		if err != nil {
			t.Errorf("failed upload %s due to %v", tt.name, err)
		}
		data, err := g.Download(tt.name)
		if err != nil {
			t.Errorf("failed download %s due to %v", tt.name, err)
		}
		if len(tt.data) != len(data) {
			t.Fatalf("wrong size, expected %v not %v", len(tt.data), len(data))
		}
		for i := range tt.data {
			if tt.data[i] != data[i] {
				t.Fatal("rebuild corrupted")
			}
		}
	}
}