package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/sftp"

	log "github.com/sirupsen/logrus"
)

type StowStorage struct {
	name      string
	kind      string
	container string
	cfg       stow.ConfigMap
}

func NewStowStorage(name string, ci ConfigItem) (Storage, error) {
	s := StowStorage{
		name: name,
		cfg:  stow.ConfigMap{},
		kind: ci.Type,
	}
	switch ci.Type {
	case local.Kind:
		s.container = ci.Cfg["path"]
		s.cfg["path"] = ci.Cfg["path"]
		if _, err := os.Stat(s.container); os.IsNotExist(err) {
			err := os.Mkdir(s.container, 0700)
			if err != nil {
				return StowStorage{}, err
			}
		}
	case s3.Kind:
		s.container = ci.Cfg["bucket"]
		s.cfg[s3.ConfigAccessKeyID] = ci.Cfg["accesskey"]
		s.cfg[s3.ConfigSecretKey] = ci.Cfg["secretkey"]
		s.cfg[s3.ConfigRegion] = ci.Cfg["region"]
		s.cfg[s3.ConfigDisableSSL] = ci.Cfg["disable_ssl"]
		if len(ci.Cfg["endpoint"]) > 0 {
			s.cfg[s3.ConfigEndpoint] = ci.Cfg["endpoint"]
		}
	case sftp.Kind:
		s.container = ci.Cfg["directory"]
		s.cfg[sftp.ConfigBasePath] = ci.Cfg["base_path"]
		s.cfg[sftp.ConfigHost] = ci.Cfg["host"]
		s.cfg[sftp.ConfigPort] = ci.Cfg["port"]
		s.cfg[sftp.ConfigUsername] = ci.Cfg["username"]
		s.cfg[sftp.ConfigPassword] = ci.Cfg["password"]
	default:
		return StowStorage{}, fmt.Errorf("unrecognized type of storage %s", ci.Type)
	}
	return s, nil
}

func (s StowStorage) Name() string {
	return s.name
}

// Get returns a file (if exist) from the storage target
func (s StowStorage) Get(name string) ([]byte, error) {
	loc, err := stow.Dial(s.kind, s.cfg)
	if err != nil {
		return nil, err
	}
	defer loc.Close()
	c, err := loc.Container(s.container)
	if err != nil {
		log.Errorf("failed to get container %s", s.container)
		return nil, err
	}
	i, err := c.Item(name)
	if err != nil {
		return []byte{}, err
	}
	r, err := i.Open()
	if err != nil {
		return []byte{}, err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// Put saves a file to the storage target
func (s StowStorage) Put(name string, data []byte) error {
	loc, err := stow.Dial(s.kind, s.cfg)
	if err != nil {
		return err
	}
	defer loc.Close()
	c, err := loc.Container(s.container)
	if err != nil {
		log.Errorf("failed to get container %s", s.container)
		return err
	}
	r := bytes.NewReader(data)
	_, err = c.Put(name, r, int64(len(data)), nil)
	return err
}

func (s StowStorage) Download(name string, filename string) (string, error) {
	loc, err := stow.Dial(s.kind, s.cfg)
	if err != nil {
		return "", err
	}
	defer loc.Close()
	c, err := loc.Container(s.container)
	if err != nil {
		return "", fmt.Errorf("failed to get container %s -> %v", s.container, err)
	}
	i, err := c.Item(name)
	if err != nil {
		return "", err
	}
	r, err := i.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()
	sh := enigma.NewStreamHash(r)
	// FIX ME: avoid to use memory
	dd, err := ioutil.ReadAll(sh.GetReader())
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filename, dd, 0755)
	if err != nil {
		return "", err
	}
	return sh.GetHash(), nil
}

func (s StowStorage) Upload(filename string, name string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	loc, err := stow.Dial(s.kind, s.cfg)
	if err != nil {
		return err
	}
	defer loc.Close()
	c, err := loc.Container(s.container)
	if err != nil {
		log.Errorf("failed to get container %s", s.container)
		return err
	}
	_, err = c.Put(name, r, info.Size(), nil)
	return err
}
