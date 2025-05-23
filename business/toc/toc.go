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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/LosAngeles971/kirinuki/business/config"
	"github.com/LosAngeles971/kirinuki/business/helpers"
	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/multistorage"
)

const (
	toc_name = "TableOfContent"
)

// TOC handles the table of content of the Kirinuki files
type TableOfContent struct {
	Lastupdate int64            `json:"lastupdate"`
	Files      []*kirinuki.File `json:"files"`
	ms         *multistorage.MultiStorage
}

type Option func(*TableOfContent) error

func WithFilename(sFile string) Option {
	return func(t *TableOfContent) error {
		data, err := multistorage.LoadDataFromTemp(sFile)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, &t)
	}
}

func New(ms *multistorage.MultiStorage, opts ...Option) (*TableOfContent, error) {
	t := &TableOfContent{
		Lastupdate: time.Now().UnixNano(),
		Files:      []*kirinuki.File{},
		ms: ms,
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
	return len(t.Files)
}

// Exist returns true if the Kirinuki file with the given name exists
func (t TableOfContent) Exist(name string) bool {
	for _, k := range t.Files {
		if name == k.Name {
			return true
		}
	}
	return false
}

// Get returns the Kirinuki file with the given name and true if the file exists
// This method returns the Kirinuki file without the data
func (t TableOfContent) Get(name string) (*kirinuki.File, bool) {
	for _, f := range t.Files {
		if name == f.Name {
			return f, true
		}
	}
	return nil, false
}

func (t *TableOfContent) Add(f *kirinuki.File) bool {
	if t.Exist(f.Name) {
		return false
	}
	t.Files = append(t.Files, f)
	return true
}

func (t TableOfContent) Find(pattern string) []*kirinuki.File {
	rr := []*kirinuki.File{}
	for _, f := range t.Files {
		match, _ := regexp.MatchString(pattern, f.Name)
		if match {
			rr = append(rr, f)
		}
	}
	return rr
}

func (t TableOfContent) save(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0755)
}

func (t *TableOfContent) getCrushMap(email string, password string) []*mosaic.Chunk {
	chunks := []*mosaic.Chunk{}
	// create one chunk for every storage
	for i := range t.ms.Names() {
		name := helpers.GetHash([]byte(fmt.Sprintf("%s_%s_%v", email, password, i)))
		c := mosaic.NewChunk(i, name, mosaic.WithFilename(config.GetTmp()+"/"+name))
		// FIX ME: toc got a full mesh, if you add a new target you got some errors a new need a redistribution
		c.TargetNames = t.ms.Names()
		chunks = append(chunks, c)
	}
	return chunks
}

func (t *TableOfContent) getKey(email string, password string) string {
	key := sha256.Sum256([]byte(email + password))
	return hex.EncodeToString(key[:])
}

func (t *TableOfContent) Store(email string, password string) error {
	tocFile := config.GetTmp() + "/" + helpers.GetFilename(24)
	err := t.save(tocFile)
	if err != nil {
		return err
	}
	chunks := t.getCrushMap(email, password)
	f := kirinuki.NewFile(toc_name, kirinuki.WithEncodedKey(t.getKey(email, password)), kirinuki.WithChunks(chunks))
	err = f.Upload(tocFile, t.ms)
	helpers.DeleteLocalFile(tocFile)
	return err
}

func (t *TableOfContent) Load(email string, password string) error {
	chunks := t.getCrushMap(email, password)
	f := kirinuki.NewFile(toc_name, kirinuki.WithEncodedKey(t.getKey(email, password)), kirinuki.WithChunks(chunks))
	tocFile := config.GetTmp() + "/" + helpers.GetFilename(24)
	err := f.Download(tocFile, t.ms)
	if err != nil {
		return fmt.Errorf("failed to download toc -> %v", err)
	}
	data, err := os.ReadFile(tocFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, t)
}