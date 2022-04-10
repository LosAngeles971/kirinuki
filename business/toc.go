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
package business

import (
	"encoding/json"
	"regexp"
	"time"
)

// TOC handles the table of content of the Kirinuki files
type TOC struct {
    Lastupdate	int64				`json:"lastupdate"`
    Kfiles		[]*Kirinuki 		`json:"kfiles"`
}

type TOCOption func(*TOC) error

// TOCWithData is used to load an existent table of content
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

// Exist returns true if the Kirinuki file with the given name exists
func (t TOC) exist(name string) bool {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return true
		}
	}
	return false
}

// Get returns the Kirinuki file with the given name and true if the file exists
// This method returns the Kirinuki file without the data
func (t TOC) get(name string) (*Kirinuki, bool) {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return k, true
		}
	}
	return nil, false
}

func (t *TOC) add(k *Kirinuki) bool {
	if t.exist(k.Name) {
		return false
	}
	t.Kfiles = append(t.Kfiles, k)
	return true
}

func (t TOC) find(pattern string) []Kirinuki {
	rr := []Kirinuki{}
	for _, k := range t.Kfiles {
		match, _ := regexp.MatchString(pattern, k.Name)
		if match {
			rr = append(rr, *k)
		}
	}
	return rr
}