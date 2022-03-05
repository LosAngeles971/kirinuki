/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Testing uploading/downloading Kirinuki's files

This testing session DOES NOT need external data, but it interacts with the system:



+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestMosaic(t *testing.T) {
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
	for _, tt := range k_data_tests {
		k1, err := NewKirinuki(WithKirinukiData(tt.name, tt.data))
		if err != nil {
			t.Fatal(err)
		}
		err = putKiriuki(k1, sm.Array())
		if err != nil {
			t.Fatal(err)
		}
		data, err := getKirinuki(k1, sm.Array())
		if err != nil {
			t.Fatal(err)
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