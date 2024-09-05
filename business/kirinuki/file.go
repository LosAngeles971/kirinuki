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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

const (
	key_size    = 32        // size of the symmetric key
	nameSize    = 24        // size of the files' names
	buffer_size = 16 * 1024 // buffer's size during encryption/decryption of files
)

// File: it represents a single file into Kirinuki domain
type File struct {
	Date         int64           `json:"date"`         // date of last upload
	Name         string          `json:"name"`         // file's name
	Chunks       []*mosaic.Chunk `json:"chunks"`       // list of file's chunks
	Symmetrickey string          `json:"symmetrickey"` // symmetric encryption key used to encrypt the file
	Checksum     string          `json:"checksum"`     // checksum of the original file
}

type FileOption func(*File)

// Usage of a random generated symmetric key to encrypt the file
func WithRandomkey() FileOption {
	return func(f *File) {
		f.setRandomKey()
	}
}

// Usage of an already existent symmetric key (the table of content uses this option)
func WithEncodedKey(key string) FileOption {
	return func(k *File) {
		k.Symmetrickey = key
	}
}

// Usage of already existent chunks (table of content uses this option)
func WithChunks(chunks []*mosaic.Chunk) FileOption {
	return func(f *File) {
		f.Chunks = chunks
	}
}

// NewFile: it creates a new File into the world of Kirinuki
func NewFile(name string, opts ...FileOption) *File {
	k := &File{
		Name: name,
		Date: time.Now().UnixNano(),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

// this method is used to assign a random symmetric key to a File deserialized by the table of content
func (f *File) setRandomKey() {
	key := sha256.Sum256(storage.GetRndBytes(key_size))
	f.Symmetrickey = hex.EncodeToString(key[:])
}

// encrypting external file into external file using the file's symmetric key
func (f *File) Encrypt(sFile string, tFile string) error {
	log.Debugf("encrypting %s to %s ...", sFile, tFile)
	key, err := hex.DecodeString(f.Symmetrickey)
	if err != nil {
		return fmt.Errorf("failed to decode symmetric key -> %v", err)
	}
	in, err := os.Open(sFile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}
	stream := cipher.NewCTR(block, iv)
	inBuf := make([]byte, buffer_size)
	for {
		n, err := in.Read(inBuf)
		if err == io.EOF {
			out.Write(iv)
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}
		stream.XORKeyStream(inBuf, inBuf[:n])
		out.Write(inBuf[:n])
	}
}

// decrypting external file into external file using the file's symmetric key
func (f *File) Decrypt(sFile string, tFile string) error {
	log.Debugf("decrypting %s to %s ...", sFile, tFile)
	key, err := hex.DecodeString(f.Symmetrickey)
	if err != nil {
		return fmt.Errorf("failed to decode symmetric key -> %v", err)
	}
	in, err := os.Open(sFile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer out.Close()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	fi, err := in.Stat()
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = in.ReadAt(iv, msgLen)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	inBuf := make([]byte, buffer_size)
	for {
		n, err := in.Read(inBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}
		if n > int(msgLen) {
			n = int(msgLen)
		}
		msgLen -= int64(n)
		stream.XORKeyStream(inBuf, inBuf[:n])
		out.Write(inBuf[:n])
	}
}

// Split: the func splits the file into chunks
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
		err = os.WriteFile(c.GetFilename(), buf, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Merge: the func merges a set of chunks into a coherent file
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

// setCrushMap: the func assigns storage targets to every chunk of file
func (f *File) setCrushMap(targets []string) {
	log.Debugf("setting crush map for file %s ...", f.Name)
	f.Chunks = []*mosaic.Chunk{}
	// create a chunk for every available storage
	for i := 0; i < len(targets); i++ {
		name := storage.GetFilename(nameSize)
		c := mosaic.NewChunk(i, name, mosaic.WithFilename(storage.GetTmp()+"/"+name))
		c.TargetNames = targets
		f.Chunks = append(f.Chunks, c)
		log.Debugf("crush map for file %s chunks %v [%s] - size %v - #targets %v", f.Name, c.Index, c.Name, c.Real_size, len(c.TargetNames))
	}
	log.Debugf("crush map for file %s #chunks %v", f.Name, len(f.Chunks))
}

// uploading an external file to the storage
func (f *File) Upload(filename string, ms *storage.MultiStorage) error {
	log.Debugf("uploading file %s as %s to storage...", filename, f.Name)
	// Table of Content comes with own key
	if len(f.Symmetrickey) < 1 {
		f.setRandomKey()
	}
	var err error
	f.Checksum, err = storage.GetFileHash(filename)
	if err != nil {
		return err
	}
	log.Debugf("kirinuki file %s from %s for hash %s", f.Name, filename, f.Checksum)
	ff := storage.GetTmp() + "/" + storage.GetFilename(nameSize)
	err = f.Encrypt(filename, ff)
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
func (f *File) Download(filename string, ms *storage.MultiStorage) error {
	log.Debugf("downloading file %s from storage to local file %s...", f.Name, filename)
	mm := mosaic.New(ms)
	err := mm.Download(f.Chunks)
	if err != nil {
		return fmt.Errorf("failed download -> %v", err)
	}
	mergeFile := storage.GetTmp() + "/" + storage.GetFilename(nameSize)
	err = f.Merge(mergeFile)
	if err != nil {
		return fmt.Errorf("failed to merge chunks to %s -> %v", mergeFile, err)
	}
	log.Debugf("merged chunks of file %s from storage to local temporary file %s", f.Name, mergeFile)
	log.Debugf("decrypting file %s with key %s", mergeFile, f.Symmetrickey)
	err = f.Decrypt(mergeFile, filename)
	if err != nil {
		return err
	}
	log.Debugf("decrypted local file %s to local file %s", mergeFile, filename)
	if len(f.Checksum) > 0 {
		h, err := storage.GetFileHash(filename)
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
