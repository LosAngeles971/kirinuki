package storage

import (
	"os"
	"testing"
)

var yCfgFile string = `
  map:
    local:
      type: "filesystem"
      config:
        path: "./"
`

var jCfgFile string = `
{
	"map": {
		"local": {
			"type": "filesystem",
			"config": {
				"path": "./"
			}
		}
	}
}			
`

func TestLoad(t *testing.T) {
	m1, err := NewStorageMap(WithYAMLData([]byte(yCfgFile)))
	if err != nil {
		t.Fatal(err)
	}
	m2, err := NewStorageMap(WithJSONData([]byte(jCfgFile)))
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []*StorageMap{m1, m2,} {
		_, err =  m.Get("local")
		if err != nil {
			t.Fatal(err)
		}
		if m.Size() != 1 {
			t.Fatalf("not expected size of %v", m.Size())
		}
	}
}

func TestAdd(t *testing.T) {
	sm, err := NewStorageMap()
	if err != nil {
		t.Fatal(err)
	}
	if sm.Size() != 0 {
		t.Fatal("storage array must be empty")
	}
	sm.Add("test", StorageItem{
		Type: "filesystem",
		Cfg: map[string]string{
			"path": os.TempDir(),
		},
	})
	if sm.Size() != 1 {
		t.Fatalf("wrong storage array size %v", sm.Size())
	}
	ss := sm.Array()
	if len(ss) != 1 {
		t.Fatalf("wrong array size %v", len(ss))
	}
}
