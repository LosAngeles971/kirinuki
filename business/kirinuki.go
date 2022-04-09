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

// getChunksNumberForKFile returns the number of chunks depending on the size of the source file
func getChunksNumberForKFile(file []byte) int {
	size := len(file)
	if size < 1000 {
		return 3
	}
	if size < 10000 {
		return 5
	}
	if size < 100000 {
		return 7
	}
	if size < 1000000 {
		return 9
	}
	return 11
}

type chunk struct {
	Name      string   `yaml:"name"`
	Real_size int      `yaml:"real_size"`
	Index     int      `yaml:"int"`
	Targets   []string `yaml:"targets"`
	data      []byte
}

type chunkOption func(*chunk) error

func withChunkData(data []byte) chunkOption {
	return func(c *chunk) error {
		if len(data) < 1 {
			return fmt.Errorf("wrong size of input data %v for a chunk", len(data))
		}
		c.Real_size = len(data)
		c.data = data
		return nil
	}
}

func withChunkName(name string) chunkOption {
	return func(c *chunk) error {
		c.Name = name
		return nil
	}
}

func newChunk(index int, opts ...chunkOption) (*chunk, error) {
	c := &chunk{}
	c.Targets = []string{}
	c.Name = newNaming().getNameForChunk()
	c.Index = index
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// setTargets assigns an array of storage targets to the chunk
func (c *chunk) setTargets(tt []storage.Storage) {
	for i := range tt {
		c.Targets = append(c.Targets, tt[i].Name())
	}
}

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

type KirinukiOption func(*Kirinuki) error

func WithKirinukiData(name string, data []byte) KirinukiOption {
	return func(k *Kirinuki) error {
		if len(data) < 1 {
			return fmt.Errorf("wrong size of input data %v for a Kirinuki", len(data))
		}
		k.Name = name
		k.Checksum = NewEnigma().hash(data)
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
}

// NewKirinukiFile creates a KirinukiFile from a new file
func NewKirinuki(opts ...KirinukiOption) (*Kirinuki, error) {
	k := &Kirinuki{
		Encryption: false,
		Padding: false,
		Date: time.Now().UnixNano(),
		Replicas: 1,
	}
	for _, opt := range opts {
		err := opt(k)
		if err != nil {
			return nil, err
		}
	}
	return k, nil
}

// Build rebuilds the original file from a fullfilled KirinukiFile
func (k *Kirinuki) Build() ([]byte, error) {
	data := []byte{}
	for index := 0; index < len(k.Chunks); index++ {
		chunk := k.Chunks[index]
		data = append(data, chunk.data...)
	}
	ck := NewEnigma().hash(data)
	if ck != k.Checksum {
		return nil, fmt.Errorf("expected checksum [%s] not [%s]", k.Checksum, ck)
	}
	return data, nil
}

// setCrushMap creates the map between chunks and storages, depending on the specific association algorithm
func (kfile *Kirinuki) setCrushMap(ss []storage.Storage) {
	for _, c := range kfile.Chunks {
		c.setTargets(ss)
	}
}