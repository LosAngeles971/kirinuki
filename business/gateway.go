package business

import (
	"fmt"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

type Gateway struct {
	keeptoc bool
	session *Session
}

func New(email string, password string, scratch bool, m *storage.StorageMap) (Gateway, error) {
	s, err := NewSession(email, password, scratch, WithStorage(m))
	if err != nil {
		return Gateway{}, err
	}
	return Gateway{
		keeptoc: false,
		session: s,
	}, nil
}

// func (g Gateway) UploadData(filename string, data []byte) error {
// 	key := GenerateKey()
// 	file, err := Encrypt(data, key)
// 	if err != nil {
// 		return err
// 	}
// 	toc, err := loadTOC(CurrentSession)
// 	if err != nil {
// 		return err
// 	}
// 	if toc.Exist(filename) {
// 		return errors.New("File already present: " + filename)
// 	}
// 	kfile, err := CreateKirinukiFile(filename, file)
// 	if err != nil {
// 		return err
// 	}
// 	kfile.Symmetrickey = key
// 	err = PutKiriukiFile(&kfile)
// 	if err != nil {
// 		return err
// 	}
// 	ok := toc.Add(&kfile)
// 	if !ok {
// 		return errors.New("bug error: file already present in TOC")
// 	}
// 	return toc.storeTOC(CurrentSession)
// }

func (g Gateway) Upload(name string, data []byte, overwrite bool) error {
	err := g.session.login()
	if err != nil {
		return err
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return err
	}
	if toc.Find(name) && !overwrite {
		return fmt.Errorf("file with name %s already exists", name)
	}
	// overwrite in any case
	k, err := NewKirinuki(WithKirinukiData(name, data))
	if err != nil {
		return err
	}
	err = putKiriuki(k, g.session.storage.Array())
	if err != nil {
		return err
	}
	toc.Add(k)
	return g.session.logout()
}

func (g Gateway) Download(name string) ([]byte, error) {
	err := g.session.login()
	if err != nil {
		return nil, err
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return nil, err
	}
	k, ok := toc.Get(name)
	if !ok {
		return nil, fmt.Errorf("file %s not present", name)
	}
	data, err := getKirinuki(k, g.session.storage.Array())
	if err != nil {
		return nil, err
	}
	if !g.keeptoc {
		err = g.session.logout()
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

func (g Gateway) Logout() error {
	if  g.session.isOpen() {
		return g.session.logout()
	}
	return nil
}