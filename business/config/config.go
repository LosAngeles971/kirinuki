package config

import (
	"os"
	"sync"
)

const (
	KIRINUKI_TMP = "KIRINUKI_TMP"
)

var lock = &sync.Mutex{}

// Storage defines a specific storage target
type StorageDef struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"` // configuration parameters
}

type Config struct {
	Storages map[string]StorageDef `yaml:"map" json:"map"` 
}

var singleton *Config

func getInstance() *Config {
    if singleton == nil {
        lock.Lock()
        defer lock.Unlock()
        if singleton == nil {
            singleton = &Config{}
        } else {
        }
    } else {
    }
    return singleton
}

func GetStorages() map[string]StorageDef {
	return getInstance().Storages
}

// func (m *MultiStorage) LoadByJSON(data []byte) error {
// 	defs := StorageDefinitions{}
// 	err := json.Unmarshal(data, &defs)
// 	if err != nil {
// 		return err
// 	}
// 	for name, def := range defs.Map {
// 		err = m.Add(name, def)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// GetTmp: it returns the local temporary directory
// The local temporary directory is always necessary, 
// since it is used to split e rebuild Kirinuki files.
func GetTmp() string {
	tmp, ok := os.LookupEnv(KIRINUKI_TMP)
	if ok {
		return tmp
	} else {
		return os.TempDir()
	}
}