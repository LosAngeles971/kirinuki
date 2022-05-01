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

func New(ms *storage.MultiStorage) *Kirinuki {
	return &Kirinuki{
		ms:       ms,
		nameSize: 48,
		tempDir: os.TempDir(),
	}
}

func (k *Kirinuki) setCrushMap(f *File) {
	f.Chunks = []*mosaic.Chunk{}
	names := k.ms.Names()
	// create a chunk for every available storage
	for i := 0; i < len(names); i++ {
		name := GetFilename(k.nameSize/2)
		c := mosaic.NewChunk(i, name, mosaic.WithFilename(k.tempDir + "/" + name))
		c.TargetNames = k.ms.Names()
		f.Chunks = append(f.Chunks, c)
	}
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
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, c := range file.Chunks {
		log.Debugf("merging %s to %s ...", c.Name, filename)
		chunk, err := os.Open(c.GetFilename())
		if err != nil {
			return err
		}
		n, err := io.Copy(f, chunk)
		if err != nil || n == 0 {
			return fmt.Errorf("error copying %v bytes to %s -> %v", n, filename, err)
		}
		chunk.Close()
	}
	return nil
}

func (k *Kirinuki) Upload(filename string, f *File) error {
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
	// Table of Content comes with own chunks
	if f.Chunks == nil || len(f.Chunks) < 1 {
		k.setCrushMap(f)
	}
	k.splitFile(ff, f)
	mm := mosaic.New(k.ms, mosaic.WithTempDir(k.tempDir))
	err = mm.Upload(f.Chunks)
	return err
}

func (k *Kirinuki) Download(f *File, filename string) error {
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
	ee := enigma.New(enigma.WithEncodedkey(f.Symmetrickey))
	err = ee.DecryptFile(ff, filename)
	if err != nil {
		return err
	}
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
	return nil
}
