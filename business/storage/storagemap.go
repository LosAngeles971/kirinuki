/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This file is in charge of uploading/downloading files into different storages

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package storage

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConfigItem struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"`
}

type ConfigMap struct {
	Map map[string]ConfigItem `yaml:"map" json:"map"`
}

type StorageMap struct {
	targets []Storage
}

type StorageMapOption func(*StorageMap) error

func WithYAMLData(data []byte) StorageMapOption {
	return func(m *StorageMap) error {
		cm := &ConfigMap{}
		err := yaml.Unmarshal(data, cm)
		if err != nil {
			return err
		}
		for name, ci := range cm.Map {
			err = m.Add(name, ci)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithJSONData(data []byte) StorageMapOption {
	return func(m *StorageMap) error {
		cm := &ConfigMap{}
		err := json.Unmarshal(data, cm)
		if err != nil {
			return err
		}
		for name, ci := range cm.Map {
			err = m.Add(name, ci)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func NewStorageMap(opts ...StorageMapOption) (*StorageMap, error) {
	m := &StorageMap{
		targets: []Storage{},
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *StorageMap) Add(name string, ci ConfigItem) error {
	ss, err := newStorage(name, ci)
	if err != nil {
		return err
	}
	m.targets = append(m.targets, ss)
	return nil
}

func (m *StorageMap) Get(name string) (Storage, error) {
	for _, ss := range m.targets {
		if ss.Name() == name {
			return ss, nil
		}
	}
	return Storage{}, fmt.Errorf("storage %s not found", name)
}

func (m *StorageMap) Size() int {
	return len(m.targets)
}

func (m *StorageMap) Array() []Storage {
	return m.targets
}