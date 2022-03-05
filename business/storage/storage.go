/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This file is in charge of uploading/downloading files into different storages

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package storage

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Storage interface {
	Name() string
	Get(name string) ([]byte, error)
	Put(name string, data []byte) error
}

type StorageItem struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"`
}

type StorageMap struct {
	Map map[string]StorageItem `yaml:"map" json:"map"`
}

type StorageMapOption func(*StorageMap) error

func WithYAMLData(data []byte) StorageMapOption {
	return func(m *StorageMap) error {
		return yaml.Unmarshal(data, m)
	}
}

func WithJSONData(data []byte) StorageMapOption {
	return func(m *StorageMap) error {
		return json.Unmarshal(data, m)
	}
}

func NewStorageMap(opts ...StorageMapOption) (*StorageMap, error) {
	m := &StorageMap{
		Map: map[string]StorageItem{},
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *StorageMap) Add(name string, si StorageItem) {
	m.Map[name] = si
}

func (m *StorageMap) Get(name string) (Storage, error) {
	s, ok := m.Map[name]
	if !ok {
		return nil, fmt.Errorf("storage %s does not exist", name)
	}
	switch s.Type {
	case "filesystem":
		return NewFilesystem(name, s)
	default:
		return nil, fmt.Errorf("unrecognized type of storage %s", s.Type)
	}
}

func (m *StorageMap) Size() int {
	return len(m.Map)
}

func (m *StorageMap) Array() []Storage {
	rr := []Storage{}
	for name := range m.Map {
		ss, err := m.Get(name)
		if err == nil {
			rr = append(rr, ss)
		}
	}
	return rr
}