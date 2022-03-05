package storage

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/local"
	log "github.com/sirupsen/logrus"
)

type filesystem struct {
	name   string
	path   string
	impl   stow.Location
}

func NewFilesystem(name string, s StorageItem) (filesystem, error) {
	ll := filesystem{
		name: name,
		path: s.Cfg["path"],
	}
	if _, err := os.Stat(ll.path); os.IsNotExist(err) {
		err := os.Mkdir(ll.path, 0700)
		if err != nil {
			return filesystem{}, err
		}
	}
	loc, err := stow.Dial(local.Kind, stow.ConfigMap{"path": ll.path})
	if err != nil {
		return filesystem{}, err
	}
	ll.impl = loc
	return ll, nil
}

func (l filesystem) Name() string {
	return l.name
}

func (l filesystem) Get(filename string) ([]byte, error) {
	c, err := l.impl.Container(l.path)
	if err != nil {
		log.Fatalf("failed to get container %s", l.path)
		return []byte{}, err
	}
	i, err := c.Item(filename)
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

func (l filesystem) Put(filename string, data []byte) error {
	c, err := l.impl.Container(l.path)
	if err != nil {
		log.Fatalf("Failed to get container %s", l.path)
		return err
	}
	r := bytes.NewReader(data)
	_, err = c.Put(filename, r, int64(len(data)), nil)
	return err
}