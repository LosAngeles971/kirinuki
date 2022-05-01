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
	"time"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
)

type File struct {
	Date         int64           `json:"date"`         // date of upload
	Name         string          `json:"name"`         // name into the Table of Content
	Chunks       []*mosaic.Chunk `json:"chunks"`       // chunks into the MultiStorage
	Symmetrickey string          `json:"symmetrickey"` // encryption key
	Checksum     string          `json:"checksum"`     // checksum of original file
}

type FileOption func(*File)

func WithRandomkey() FileOption {
	return func(k *File) {
		ee := enigma.New(enigma.WithRandomkey())
		k.Symmetrickey = ee.GetEncodedKey()
	}
}

func WithEncodedKey(key string) FileOption {
	return func(k *File) {
		k.Symmetrickey = key
	}
}

func WithChunks(chunks []*mosaic.Chunk) FileOption {
	return func(f *File) {
		f.Chunks = chunks
	}
}

// NewKirinukiFile creates a KirinukiFile from a generic file
func NewKirinuki(name string, opts ...FileOption) *File {
	k := &File{
		Name: name,
		Date: time.Now().UnixNano(),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}
