/*
 * Created on Sat Apr 09 2022
 * Author @LosAngeles971
 *
 * The MIT License (MIT)
 * Copyright (c) 2022 @LosAngeles971
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software
 * and associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
 * TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	TEMP_STORAGE = "temp"
)

// ConfigItem contains configurations for a specific storage target
type ConfigItem struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"`
}

// ConfigMap contains configurations for a plurality of storage targets
type ConfigMap struct {
	Map map[string]ConfigItem `yaml:"map" json:"map"`
}

// StorageMap is a manager of a plurality of storage targets
type StorageMap struct {
	targets []Storage
}

type StorageMapOption func(*StorageMap) error



// WithYAMLData populates a StorageMap from a YAML data
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

// WithJSONData populates a StorageMap from a JSON data
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

// WithTemp add local temporary directory as a storage target
func WithTemp() StorageMapOption {
	return func(m *StorageMap) error {
		m.Add(TEMP_STORAGE, ConfigItem{
			Type: "local",
			Cfg: map[string]string{
				"path": os.TempDir(),
			},
		})
		return nil
	}
}

// NewStorageMap creates a new StorageMap (empty or populated depending on the options)
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

func Init(filename string) ([]byte, error) {
	cm := ConfigMap{
		Map: map[string]ConfigItem{},
	}
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		return yaml.Marshal(&cm)
	} else {
		return json.Marshal(&cm)
	}
}
