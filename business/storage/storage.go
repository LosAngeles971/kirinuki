/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

single storage management by means of stow

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package storage

import (
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

// Storage allows to get/put files to a generic storage supported by stow
type Storage struct {
	name      string
	kind      string
	container string
	cfg       stow.ConfigMap
}

// newStorage creates a Storage struct starting from a configuration item, the latter bears configurations 
// depending on the specific target (filesystem, s3, ...)
func newStorage(name string, ci ConfigItem) (Storage, error) {
	s := Storage{
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
				return Storage{}, err
			}
		}
	case s3.Kind:
		s.container = ci.Cfg["bucket"]
		s.cfg[s3.ConfigAccessKeyID] = ci.Cfg["accesskey"]
		s.cfg[s3.ConfigSecretKey] = ci.Cfg["secretkey"]
		s.cfg[s3.ConfigRegion] = ci.Cfg["region"]
	case sftp.Kind:
		s.container = ci.Cfg["directory"]
		s.cfg[sftp.ConfigHost] = ci.Cfg["host"]
		s.cfg[sftp.ConfigUsername] = ci.Cfg["username"]
		s.cfg[sftp.ConfigPassword] = ci.Cfg["password"]
	default:
		return Storage{}, fmt.Errorf("unrecognized type of storage %s", ci.Type)
	}
	return s, nil
}

// Name returns the (arbitrary) name of the specific storage target
func (s Storage) Name() string {
	return s.name
}

// Get returns a file (if exist) from the storage target
func (s Storage) Get(name string) ([]byte, error) {
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
func (s Storage) Put(filename string, data []byte) error {
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
	_, err = c.Put(filename, r, int64(len(data)), nil)
	return err
}
