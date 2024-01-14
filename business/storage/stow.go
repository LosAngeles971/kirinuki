package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

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

func NewStowStorage(name string, def StorageDefinition) (Storage, error) {
	s := StowStorage{
		name: name,
		cfg:  stow.ConfigMap{},
		kind: def.Type,
	}
	switch def.Type {
	case local.Kind:
		s.container = def.Cfg["path"]
		s.cfg["path"] = def.Cfg["path"]
		if _, err := os.Stat(s.container); os.IsNotExist(err) {
			err := os.Mkdir(s.container, 0700)
			if err != nil {
				return StowStorage{}, err
			}
		}
	case s3.Kind:
		s.container = def.Cfg["bucket"]
		s.cfg[s3.ConfigAccessKeyID] = def.Cfg["accesskey"]
		s.cfg[s3.ConfigSecretKey] = def.Cfg["secretkey"]
		s.cfg[s3.ConfigRegion] = def.Cfg["region"]
		s.cfg[s3.ConfigDisableSSL] = def.Cfg["disable_ssl"]
		if len(def.Cfg["endpoint"]) > 0 {
			s.cfg[s3.ConfigEndpoint] = def.Cfg["endpoint"]
		}
	case sftp.Kind:
		s.container = def.Cfg["directory"]
		s.cfg[sftp.ConfigBasePath] = def.Cfg["base_path"]
		s.cfg[sftp.ConfigHost] = def.Cfg["host"]
		s.cfg[sftp.ConfigPort] = def.Cfg["port"]
		s.cfg[sftp.ConfigUsername] = def.Cfg["username"]
		s.cfg[sftp.ConfigPassword] = def.Cfg["password"]
	default:
		return StowStorage{}, fmt.Errorf("unrecognized type of storage %s", def.Type)
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
	DeleteLocalFile(filename)
	loc, err := stow.Dial(s.kind, s.cfg)
	if err != nil {
		return "", err
	}
	defer loc.Close()
	c, err := loc.Container(s.container)
	if err != nil {
		return "", fmt.Errorf("failed to get container %s -> %v", s.container, err)
	}
	log.Debugf("storage %s is downloading %s/%s to %s", s.name, s.container, name, filename)
	i, err := c.Item(name)
	if err != nil {
		return "", err
	}
	r, err := i.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()
	sh := NewStreamHash(r)
	// FIX ME: avoid to use memory
	dd, err := ioutil.ReadAll(sh.GetReader())
	if err != nil {
		return "", err
	}
	if len(dd) < 1 {
		return "", fmt.Errorf("download file %s got 0 bytes", filename)
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
	if info.Size() < 1 {
		return fmt.Errorf("cannot upload file %s it got 0 bytes", filename)
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
	log.Debugf("uploaded file %s bytes %v", filename, info.Size())
	return err
}
