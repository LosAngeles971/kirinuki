/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Implmentation of the Table of Content

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"encoding/json"
	"time"
)

type TOC struct {
    Lastupdate	int64				`json:"lastupdate"`
    Kfiles		[]*Kirinuki 		`json:"kfiles"`
}

type TOCOption func(*TOC) error

func TOCWithData(data []byte) TOCOption {
	return func(t *TOC) error {
		return json.Unmarshal(data, &t)
	}
}

func newTOC(opts ...TOCOption) (*TOC, error) {
	t := &TOC{
		Lastupdate: time.Now().UnixNano(),
		Kfiles: []*Kirinuki{},
	}
	for _, opt := range opts {
		err := opt(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (t TOC) Find(name string) bool {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return true
		}
	}
	return false
}

func (t TOC) Get(name string) (*Kirinuki, bool) {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return k, true
		}
	}
	return nil, false
}

func (t *TOC) Add(k *Kirinuki) bool {
	if t.Find(k.Name) {
		return false
	}
	t.Kfiles = append(t.Kfiles, k)
	return true
}

func (t TOC) serialize() ([]byte, error) {
	return json.Marshal(t)
}