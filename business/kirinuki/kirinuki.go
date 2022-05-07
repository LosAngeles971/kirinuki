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
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	log "github.com/sirupsen/logrus"
)

func GetFilename(size int) string {
	dd := enigma.GetRndBytes(size)
	return hex.EncodeToString(dd)
}

type Kirinuki struct {
	ms       *storage.MultiStorage
	nameSize int
	tempDir string
}

type Option func(*Kirinuki)

func WithTempDir(tempDir string) Option {
	return func(k *Kirinuki) {
		k.tempDir = tempDir
	}
}

func New(ms *storage.MultiStorage, opts ...Option) *Kirinuki {
	k := &Kirinuki{
		ms:       ms,
		nameSize: 48,
		tempDir: os.TempDir(),
	}
	for _, o := range opts {
		o(k)
	}
	return k
}

func (k *Kirinuki) setCrushMap(f *File) {
	log.Debugf("setting crush map for file %s ...", f.Name)
	f.Chunks = []*mosaic.Chunk{}
	names := k.ms.Names()
	// create a chunk for every available storage
	for i := 0; i < len(names); i++ {
		name := GetFilename(k.nameSize/2)
		c := mosaic.NewChunk(i, name, mosaic.WithFilename(k.tempDir + "/" + name))
		c.TargetNames = k.ms.Names()
		f.Chunks = append(f.Chunks, c)
		log.Debugf("crush map for file %s chunks %v [%s] - size %v - #targets %v", f.Name, c.Index, c.Name, c.Real_size, len(c.TargetNames))
	}
	log.Debugf("crush map for file %s #chunks %v", f.Name, len(f.Chunks))
}

func (k *Kirinuki) splitFile(filename string, file *File) error {
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

func (k *Kirinuki) mergeChunks(file *File, filename string) error {
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

func (k *Kirinuki) Upload(filename string, f *File) error {
	log.Debugf("uploading file %s as %s to storage...", filename, f.Name)
	// Table of Content comes with own key
	if len(f.Symmetrickey) < 1 {
		ee := enigma.New(enigma.WithRandomkey())
		f.Symmetrickey = ee.GetEncodedKey()
	}
	var err error
	f.Checksum, err = enigma.GetFileHash(filename)
	if err != nil {
		return err
	}
	log.Debugf("kirinuki file %s from %s for hash %s", f.Name, filename, f.Checksum)
	ff := k.tempDir + "/" + GetFilename(k.nameSize)
	ee := enigma.New(enigma.WithEncodedkey(f.Symmetrickey))
	err = ee.EncryptFile(filename, ff)
	if err != nil {
		return err
	}
	log.Debugf("kirinuki file %s from %s encrypted with key %s", f.Name, filename, f.Symmetrickey)
	// Table of Content comes with own chunks
	if f.Chunks == nil || len(f.Chunks) < 1 {
		k.setCrushMap(f)
	}
	k.splitFile(ff, f)
	mm := mosaic.New(k.ms, mosaic.WithTempDir(k.tempDir))
	err = mm.Upload(f.Chunks)
	log.Debugf("uploaded file %s as %s to storage with error %v", filename, f.Name, err)
	return err
}

func (k *Kirinuki) Download(f *File, filename string) error {
	log.Debugf("downloading file %s from storage to local file %s...", f.Name, filename)
	mm := mosaic.New(k.ms, mosaic.WithTempDir(k.tempDir))
	err := mm.Download(f.Chunks)
	if err != nil {
		return err
	}
	ff := k.tempDir + "/" + GetFilename(k.nameSize)
	err = k.mergeChunks(f, ff)
	if err != nil {
		return err
	}
	log.Debugf("merged chunks of file %s from storage to local temporary file %s", f.Name, ff)
	log.Debugf("decrypting file %s with key %s", ff, f.Symmetrickey)
	ee := enigma.New(enigma.WithEncodedkey(f.Symmetrickey))
	err = ee.DecryptFile(ff, filename)
	if err != nil {
		return err
	}
	log.Debugf("decrypted local file %s to local file %s", ff, filename)
	if len(f.Checksum) > 0 {
		h, err := enigma.GetFileHash(filename)
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
