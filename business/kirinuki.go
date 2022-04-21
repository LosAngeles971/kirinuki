/*
 * Created on Sun Apr 10 2022
 * Author @LosAngeles971
 *
 * This software is licensed under GNU General Public License v2.0
 * Copyright (c) 2022 @LosAngeles971
 *
 * The GNU GPL is the most widely used free software license and has a strong copyleft requirement.
 * When distributing derived works, the source code of the work must be made available under the same license.
 * There are multiple variants of the GNU GPL, each with different requirements.
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
	chunkNames   []string
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

// NewKirinukiFile creates a KirinukiFile from a generic file
func NewKirinuki(name string, chunkNames []string, opts ...KirinukiOption) *Kirinuki {
	k := &Kirinuki{
		Name:       name,
		Encryption: false,
		Padding:    false,
		Date:       time.Now().UnixNano(),
		Replicas:   1,
		chunkNames: chunkNames,
	}
	k.Chunks = []*chunk{}
		for i, name := range k.chunkNames {
			k.Chunks = append(k.Chunks, newChunk(i, name))
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
	chunks, err := splitFile(data, len(k.chunkNames))
	if err != nil {
		return err
	}
	k.Chunks = []*chunk{}
	for i, name := range k.chunkNames {
		k.Chunks = append(k.Chunks, newChunk(i, name, withChunkData(chunks[i])))
	}
	return nil
}

// Build rebuilds the original data file from fullfilled array of chunks
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
	// Kirinki file for TOC does not have an expected checksum
	if len(k.Checksum) > 0 {
		ck := newEnigma().hash(orig)
		if ck != k.Checksum {
			return nil, fmt.Errorf("wrong checksum wanted [%s] having [%s]", k.Checksum, ck)
		}
	}
	return orig, nil
}

// setCrushMap creates the map between chunks and storages, depending on the specific association algorithm
func (kfile *Kirinuki) setCrushMap(ss []storage.Storage) {
	for _, c := range kfile.Chunks {
		c.setTargets(ss)
	}
}
