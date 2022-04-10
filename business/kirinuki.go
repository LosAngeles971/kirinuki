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
	"fmt"
	"time"

	"github.com/LosAngeles971/kirinuki/business/storage"
	log "github.com/sirupsen/logrus"
)

type Kirinuki struct {
	Date         int64    `json:"date"`
	Name         string   `json:"name"`
	Encryption   bool     `json:"encryption"`
	Padding      bool     `json:"padding"`
	Chunks       []*chunk `json:"chunks"`
	Symmetrickey string   `json:"symmetrickey"`
	Replicas     int      `json:"replicas"`
	Checksum     string   `json:"checksum"`
}

type KirinukiOption func(*Kirinuki)

func WithRandomkey() KirinukiOption {
	return func(k *Kirinuki) {
		k.Encryption = true
		ee := newEnigma(withRandomkey())
		k.Symmetrickey = ee.getEncodedKey()
	}
}

func WithEncodedKey(key string) KirinukiOption {
	return func(k *Kirinuki) {
		k.Encryption = true
		k.Symmetrickey = key
	}
}

func NewKirinukiTOC(email string, password string) *Kirinuki {
	return NewKirinuki(newEnigma().hash([]byte(email)), WithEncodedKey(newEnigma(withMainkey(email, password)).getEncodedKey()))
}

// NewKirinukiFile creates a KirinukiFile from a generic file
func NewKirinuki(name string, opts ...KirinukiOption) *Kirinuki {
	k := &Kirinuki{
		Name: name,
		Encryption: false,
		Padding:    false,
		Date:       time.Now().UnixNano(),
		Replicas:   1,
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

func (k *Kirinuki) addData(orig []byte) error {
	if len(orig) < 1 {
		return fmt.Errorf("cannot add empty data to Kirinuki file %s", k.Name)
	}
	k.Checksum = newEnigma().hash(orig)
	var data []byte
	var err error
	if k.Encryption {
		ee := newEnigma(withEncodedkey(k.Symmetrickey))
		data, err = ee.encrypt(orig)
		if err != nil {
			return err
		}
	} else {
		data = orig
	}
	chunks, err := splitFile(data, getChunksNumberForKFile(data))
	if err != nil {
		return err
	}
	k.Chunks = []*chunk{}
	for index := range chunks {
		ch, err := newChunk(index, withChunkData(chunks[index]))
		if err != nil {
			log.Errorf("failed at chunk %v", index)
			return err
		}
		k.Chunks = append(k.Chunks, ch)
	}
	return nil
}

// Build rebuilfs the original data file from fullfilled array of chunks
func (k *Kirinuki) Build() ([]byte, error) {
	data := []byte{}
	for index := 0; index < len(k.Chunks); index++ {
		chunk := k.Chunks[index]
		data = append(data, chunk.data...)
	}
	var orig []byte
	var err error
	if k.Encryption {
		ee := newEnigma(withEncodedkey(k.Symmetrickey))
		orig, err = ee.decrypt(data)
		if err != nil {
			return nil, err
		}
	} else {
		orig = data
	}
	ck := newEnigma().hash(orig)
	if ck != k.Checksum {
		return nil, fmt.Errorf("expected checksum [%s] not [%s]", k.Checksum, ck)
	}
	return orig, nil
}

// setCrushMap creates the map between chunks and storages, depending on the specific association algorithm
func (kfile *Kirinuki) setCrushMap(ss []storage.Storage) {
	for _, c := range kfile.Chunks {
		c.setTargets(ss)
	}
}
