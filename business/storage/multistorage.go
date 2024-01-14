package storage

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

import (
	"encoding/json"
	"fmt"

	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/sftp"
)

const (
	TEMP_STORAGE = "temp"
)

// StorageDefinition defines a specific storage target
type StorageDefinition struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"` // configuration parameters
}

// StorageDefinitions contains a list of defined storage targets
type StorageDefinitions struct {
	Map map[string]StorageDefinition `yaml:"map" json:"map"` 
}

// MultiStorage: it allows to deal with all defined storage targets (upload, download, delete, ...)
type MultiStorage struct {
	targets []Storage // list of defined storage targets
}

type MultiStorageOption func(*MultiStorage) error

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

// Add: it allows to manually add a storage target defined by ci
func (m *MultiStorage) Add(name string, def StorageDefinition) error {
	var ss Storage
	var err error
	switch def.Type {
	case local.Kind, s3.Kind:
		ss, err = NewStowStorage(name, def)
	case sftp.Kind:
		ss, err = NewSFTP(name, def.Cfg), nil
	default:
		return fmt.Errorf("unrecognized type of storage %s", def.Type)
	}
	if err != nil {
		return err
	}
	m.targets = append(m.targets, ss)
	return nil
}

func (m *MultiStorage) LoadByJSON(data []byte) error {
	defs := StorageDefinitions{}
	err := json.Unmarshal(data, &defs)
	if err != nil {
		return err
	}
	for name, def := range defs.Map {
		err = m.Add(name, def)
		if err != nil {
			return err
		}
	}
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
	return m.Add(name, StorageDefinition{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
}

// Names: it returns the list of storage targets' names
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