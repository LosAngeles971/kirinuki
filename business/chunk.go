/*
 * Created on Sun Apr 10 2022
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

	"github.com/LosAngeles971/kirinuki/business/storage"
)

type chunk struct {
	Name      string   `json:"name"`
	Real_size int      `json:"real_size"`
	Index     int      `json:"int"`
	Targets   []string `json:"targets"`
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