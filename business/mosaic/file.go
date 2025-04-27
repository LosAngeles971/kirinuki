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
package mosaic

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/LosAngeles971/kirinuki/business/config"
	"github.com/LosAngeles971/kirinuki/business/helpers"
	"github.com/LosAngeles971/kirinuki/business/multistorage"

	log "github.com/sirupsen/logrus"
)

const (
	key_size    = 32        // size of the symmetric key
	nameSize    = 24        // size of the files' names
	buffer_size = 16 * 1024 // buffer's size during encryption/decryption of files
)

// File: a single file into Kirinuki domain
type File struct {
	Encrypted bool `json:"encrypted"`         // date of last modification
	Date         int64    `json:"date"`         // date of last modification
	Name         string   `json:"name"`         // file's name
	Chunks       []*Chunk `json:"chunks"`       // list of chunks into which the file is destructured
	Symmetrickey string   `json:"symmetrickey"` // symmetric encryption key used to encrypt the file
	Checksum     string   `json:"checksum"`     // checksum of the original file
}

type FileOption func(*File)

/* // Usage of a random generated symmetric key to encrypt the file
func WithRandomkey() FileOption {
	return func(f *File) {
		f.Symmetrickey = helpers.GetRndEncodedKey()
	}
}

// Usage of a provided symmetric key (the table of content uses this option)
func WithEncodedKey(key string) FileOption {
	return func(k *File) {
		k.Symmetrickey = key
	}
}

// Usage of provide chunks
// Only the "table of content" uses this option
func WithChunks(chunks []*Chunk) FileOption {
	return func(f *File) {
		f.Chunks = chunks
	}
} */

// It creates a new File object
// It may include the data (in case of WithChunks option) or not
func newFile(name string, targets []string, opts ...FileOption) *File {
	k := &File{
		Encrypted: false,
		Name: name,
		Date: time.Now().UnixNano(),
	}
	for _, opt := range opts {
		opt(k)
	}
	k.Chunks = getCrushMap(targets)
	return k
}

// splitting of external file into chunks
func (file *File) Split(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()
	tot := int64(0)
	for i, c := range file.Chunks {
		var chunk_size int64
		if i == len(file.Chunks)-1 {
			chunk_size = size - tot
		} else {
			chunk_size = size / int64(len(file.Chunks))
			tot += chunk_size
		}
		buf := make([]byte, chunk_size)
		n, err := f.Read(buf)
		if err != nil || n != int(chunk_size) {
			return err
		}
		err = ioutil.WriteFile(c.GetFilename(), buf, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// merging external chunks into external file
func (file *File) Merge(filename string) error {
	log.Debugf("merging #chunks %v to file %s [%s]...", len(file.Chunks), file.Name, filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	sort.Slice(file.Chunks, func(i, j int) bool {
		return file.Chunks[i].Index < file.Chunks[j].Index
	})
	for _, c := range file.Chunks {
		chunk, err := os.Open(c.GetFilename())
		if err != nil {
			return err
		}
		n, err := io.Copy(f, chunk)
		if err != nil || n == 0 {
			log.Debugf("failed merging chunk %v %s [%s] to filename %s - bytes %v - err -> %v", c.Index, c.Name, c.GetFilename(), filename, n, err)
		}
		chunk.Close()
		log.Debugf("merged chunk %v %s [%s] to filename %s - bytes %v", c.Index, c.Name, c.GetFilename(), filename, n)
	}
	log.Debugf("merging #chunks %v to file %s [%s] completed", len(file.Chunks), file.Name, filename)
	return nil
}

// uploading an external file to the storage
func (f *File) Upload(filename string, ms *multistorage.MultiStorage) error {
	log.Debugf("uploading file %s as %s to storage...", filename, f.Name)
	// Table of Content comes with own key
	if len(f.Symmetrickey) < 1 {
		f.Symmetrickey = helpers.GetRndEncodedKey()
	}
	var err error
	f.Checksum, err = helpers.GetFileHash(filename)
	if err != nil {
		return err
	}
	log.Debugf("kirinuki file %s from %s for hash %s", f.Name, filename, f.Checksum)
	ff := config.GetTmp() + "/" + helpers.GetFilename(nameSize)
	err = helpers.EncryptFile(filename, ff, f.Symmetrickey)
	if err != nil {
		return fmt.Errorf("failed to encrypt file -> %v", err)
	}
	log.Debugf("kirinuki file %s from %s encrypted with key %s", f.Name, filename, f.Symmetrickey)
	// Table of Content comes with own chunks
	if f.Chunks == nil || len(f.Chunks) < 1 {
		f.setCrushMap(ms.Names())
	}
	f.Split(ff)
	mm := mosaic.New(ms)
	err = mm.Upload(f.Chunks)
	log.Debugf("uploaded file %s as %s to storage with error %v", filename, f.Name, err)
	return err
}

// downloading an external file from the storage
func (f *File) Download(filename string, ms *multistorage.MultiStorage) error {
	log.Debugf("downloading file %s from storage to local file %s...", f.Name, filename)
	mm := mosaic.New(ms)
	err := mm.Download(f.Chunks)
	if err != nil {
		return fmt.Errorf("failed download -> %v", err)
	}
	mergeFile := config.GetTmp() + "/" + helpers.GetFilename(nameSize)
	err = f.Merge(mergeFile)
	if err != nil {
		return fmt.Errorf("failed to merge chunks to %s -> %v", mergeFile, err)
	}
	log.Debugf("merged chunks of file %s from storage to local temporary file %s", f.Name, mergeFile)
	log.Debugf("decrypting file %s with key %s", mergeFile, f.Symmetrickey)
	err = helpers.DecryptFile(mergeFile, filename, f.Symmetrickey)
	if err != nil {
		return err
	}
	log.Debugf("decrypted local file %s to local file %s", mergeFile, filename)
	if len(f.Checksum) > 0 {
		h, err := helpers.GetFileHash(filename)
		if err != nil {
			return err
		}
		if f.Checksum != h {
			return fmt.Errorf("kirinuki file %s expects hash %s not %s", f.Name, f.Checksum, h)
		}
	} else {
		log.Warningf("kirinuki file %s does not have checksum", f.Name)
	}
	log.Debugf("download file %s from storage to local file %s completed", f.Name, filename)
	return nil
}
