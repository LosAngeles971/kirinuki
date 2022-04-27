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
package toc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
)

// TOC handles the table of content of the Kirinuki files
type TableOfContent struct {
	Lastupdate int64                `json:"lastupdate"`
	Kfiles     []*kirinuki.Kirinuki `json:"kfiles"`
}

type Option func(*TableOfContent) error

// TOCWithData is used to load an existent table of content
func WithData(data []byte) Option {
	return func(t *TableOfContent) error {
		return json.Unmarshal(data, &t)
	}
}

func WithFilename(sFile string) Option {
	return func(t *TableOfContent) error {
		data, err := ioutil.ReadFile(sFile)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, &t)
	}
}

func New(opts ...Option) (*TableOfContent, error) {
	t := &TableOfContent{
		Lastupdate: time.Now().UnixNano(),
		Kfiles:     []*kirinuki.Kirinuki{},
	}
	for _, opt := range opts {
		err := opt(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (t *TableOfContent) Size() int {
	return len(t.Kfiles)
}

// Exist returns true if the Kirinuki file with the given name exists
func (t TableOfContent) Exist(name string) bool {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return true
		}
	}
	return false
}

// Get returns the Kirinuki file with the given name and true if the file exists
// This method returns the Kirinuki file without the data
func (t TableOfContent) Get(name string) (*kirinuki.Kirinuki, bool) {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return k, true
		}
	}
	return nil, false
}

func (t *TableOfContent) Add(k *kirinuki.Kirinuki) bool {
	if t.Exist(k.Name) {
		return false
	}
	t.Kfiles = append(t.Kfiles, k)
	return true
}

func (t TableOfContent) Find(pattern string) []kirinuki.Kirinuki {
	rr := []kirinuki.Kirinuki{}
	for _, k := range t.Kfiles {
		match, _ := regexp.MatchString(pattern, k.Name)
		if match {
			rr = append(rr, *k)
		}
	}
	return rr
}

func (t TableOfContent) Save(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0755)
}

func GetChunks(email string, password string, ss []storage.Storage, tempDir string) []*mosaic.Chunk {
	chunks := []*mosaic.Chunk{}
	for i := range ss {
		name := enigma.GetHash([]byte(fmt.Sprintf("%s_%s_%v", email, password, i)))
		c := mosaic.NewChunk(i, name, mosaic.WithFilename(tempDir + "/" + name))
		// FIX ME: toc got a full mesh, if you add a new target you got some errors a new need a redistribution
		c.Targets = ss
		chunks = append(chunks, c)
	}
	return chunks
}