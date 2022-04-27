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
package kirinuki

import (
	"os"
	"time"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
)

type Kirinuki struct {
	Date         int64           `json:"date"`
	Name         string          `json:"name"`
	Chunks       []*mosaic.Chunk `json:"chunks"`
	Symmetrickey string          `json:"symmetrickey"`
	Checksum     string          `json:"checksum"`
}

type KirinukiOption func(*Kirinuki)

func WithRandomkey() KirinukiOption {
	return func(k *Kirinuki) {
		ee := enigma.New(enigma.WithRandomkey())
		k.Symmetrickey = ee.GetEncodedKey()
	}
}

func WithEncodedKey(key string) KirinukiOption {
	return func(k *Kirinuki) {
		k.Symmetrickey = key
	}
}

// NewKirinukiFile creates a KirinukiFile from a generic file
func NewKirinuki(name string, opts ...KirinukiOption) *Kirinuki {
	k := &Kirinuki{
		Name:       name,
		Date:       time.Now().UnixNano(),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

func (k *Kirinuki) Upload(filename string, ss []storage.Storage) error {
	mm := mosaic.New(mosaic.WithStorage(ss))
	ff := os.TempDir() + "/" + mosaic.GetFilename(24)
	ee := enigma.New(enigma.WithEncodedkey(k.Symmetrickey))
	err := ee.EncryptFile(filename, ff)
	if err != nil {
		return err
	}
	k.Chunks, err = mm.Upload(ff)
	return err
}

func (k *Kirinuki) Download(filename string, ss []storage.Storage) error {
	mm := mosaic.New(mosaic.WithStorage(ss))
	ff := os.TempDir() + "/" + mosaic.GetFilename(24)
	err := mm.Download(k.Chunks, ff)
	if err != nil {
		return err
	}
	ee := enigma.New(enigma.WithEncodedkey(k.Symmetrickey))
	err = ee.DecryptFile(ff, filename)
	if err != nil {
		return err
	}
	return nil
}