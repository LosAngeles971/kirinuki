/*
 * Created on Sun Apr 10 2022
 * Author @LosAngeles971
 *
 * This software is licensed under GNU General Public License v2.0
 * Copyright (c) 2022 @LosAngeles971
 *
 * The GNU GPL is the most widely used free software license and has a strong copyleft requirement.
 * When distributing derived works, the source code of the work must be made available under the same license.
 * There are multiple variants of the GNU GPL, each with different requirements.
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

	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/sftp"
	"gopkg.in/yaml.v3"
)

const (
	TEMP_STORAGE = "temp"
)

// ConfigMap contains configurations for a plurality of storage targets
type ConfigMap struct {
	Map map[string]ConfigItem `yaml:"map" json:"map"`
}

// StorageMap is a manager of a plurality of storage targets
type MultiStorage struct {
	targets []Storage
}

type MultiStorageOption func(*MultiStorage) error

// WithYAMLData populates a StorageMap from a YAML data
func WithYAMLData(data []byte) MultiStorageOption {
	return func(m *MultiStorage) error {
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
func WithJSONData(data []byte) MultiStorageOption {
	return func(m *MultiStorage) error {
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

// NewStorageMap creates a new StorageMap (empty or populated depending on the options)
func NewMultiStorage(opts ...MultiStorageOption) (*MultiStorage, error) {
	m := &MultiStorage{
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

func (m *MultiStorage) Add(name string, ci ConfigItem) error {
	var ss Storage
	var err error
	switch ci.Type {
	case local.Kind, s3.Kind:
		ss, err = NewStowStorage(name, ci)
	case sftp.Kind:
		ss, err = NewSFTP(name, ci.Cfg), nil
	default:
		return fmt.Errorf("unrecognized type of storage %s", ci.Type)
	}
	if err != nil {
		return err
	}
	m.targets = append(m.targets, ss)
	return nil
}

func (m *MultiStorage) get(sName string) (Storage, error) {
	for _, ss := range m.targets {
		if ss.Name() == sName {
			return ss, nil
		}
	}
	return nil, fmt.Errorf("storage [%s] does not exist", sName)
}

func (m *MultiStorage) AddLocal(name string, base string) error {
	return m.Add(name, ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
}

func (m *MultiStorage) Names() []string {
	names := []string{}
	for _, ss := range m.targets {
		names = append(names, ss.Name())
	}
	return names
}

func (m *MultiStorage) Exist(name string) bool {
	for _, ss := range m.targets {
		if ss.Name() == name {
			return true
		}
	}
	return false
}

func (m *MultiStorage) Size() int {
	return len(m.targets)
}

func (m *MultiStorage) Get(sName string, name string) ([]byte, error) {
	ss, err := m.get(sName)
	if err != nil {
		return nil, err
	}
	return ss.Get(name)
}

func (m *MultiStorage) Put(sName string, name string, data []byte) error {
	ss, err := m.get(sName)
	if err != nil {
		return err
	}
	return ss.Put(name, data)
}

func (m *MultiStorage) Download(sName string, name string, filename string) (string, error) {
	ss, err := m.get(sName)
	if err != nil {
		return "", err
	}
	return ss.Download(name, filename)
}

func (m *MultiStorage) Upload(sName string, filename string, name string) error {
	ss, err := m.get(sName)
	if err != nil {
		return nil
	}
	return ss.Upload(filename, name)
}