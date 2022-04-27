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
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/LosAngeles971/kirinuki/business/dust"
	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/storage"
	log "github.com/sirupsen/logrus"
)

const (
	STATE_COMPLETED = "completed"
	STATE_FAILED    = "failed"
	STATE_MISSING   = "missing"
)

func GetFilename(size int) string {
	dd := enigma.GetRndBytes(size)
	return hex.EncodeToString(dd)
}

// Mosaic is in charge of providing upload/download functionalities for Kirinuki files
type Mosaic struct {
	ss          []storage.Storage
	max_threads int
	tempDir     string
	nameSize 	int
}

type MosaicOption func(*Mosaic)

func WithTempDir(tempDir string) MosaicOption {
	return func(m *Mosaic) {
		m.tempDir = tempDir
	}
}

func WithStorage(ss []storage.Storage) MosaicOption {
	return func(m *Mosaic) {
		m.ss = ss
	}
}

func New(opts ...MosaicOption) *Mosaic {
	sm, err := storage.NewStorageMap(storage.WithTemp())
	if err != nil {
		panic(err)
	}
	m := &Mosaic{
		ss:        sm.Array(),
		max_threads: 4,
		tempDir: os.TempDir(),
		nameSize: 48,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (m *Mosaic) getMosaic() []*Chunk {
	chunks := []*Chunk{}
	for i := 0; i < len(m.ss); i++ {
		c := NewChunk(i, GetFilename(m.nameSize / 2))
		c.filename = m.tempDir + "/" + c.Name
		c.Targets = m.ss
		for _, s := range m.ss {
			c.status[s.Name()] = STATE_MISSING
		}
		chunks = append(chunks, c)
	}
	return chunks
}

func (m *Mosaic) getTarget(chunks []*Chunk) (*Chunk, storage.Storage) {
	for _, c := range chunks {
		for _, s := range m.ss {
			status, ok := c.status[s.Name()]
			if ok && status == STATE_MISSING {
				return c, s
			}
		}
	}
	return nil, nil
}

func (m *Mosaic) isComplete(chunks []*Chunk) (bool, bool) {
	completed := true
	for _, c := range chunks {
		ok := false
		for _, status := range c.status {
			if status == STATE_MISSING {
				return false, false
			}
			if status == STATE_COMPLETED {
				ok = true
			}
		}
		if !ok {
			completed = false
		}
	}
	return true, completed
}

func (m *Mosaic) download(chunk *Chunk) error {
	log.Debugf("downloading chunk %s", chunk.Name)
	for _, s := range m.ss {
		chunk.download(s)
		if chunk.err != nil {
			log.Errorf("failed download chunk %s from %s -> %v", chunk.Name, s.Name(), chunk.err)
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to download chunk %s from all targets", chunk.Name)
}

func (m *Mosaic) upload(chunks []*Chunk, filename string) ([]*Chunk, error) {
	chunkFiles := []string{}
	for _, c := range chunks {
		chunkFiles = append(chunkFiles, c.filename)
	}
	err := dust.SplitFile(filename, chunkFiles)
	if err != nil {
		return nil, err
	}
	log.Debugf("uploading [%v] chunks", len(chunks))
	for {
		end, completed := m.isComplete(chunks)
		if end {
			if completed {
				return chunks, nil
			} else {
				return chunks, fmt.Errorf("failed to upload %s", filename)
			}
		}
		var wg sync.WaitGroup
		for i :=0; i < m.max_threads; i++ {
			c, target := m.getTarget(chunks)
			if target != nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.upload(target)
				}()
				if c.err != nil {
					return chunks, c.err
				}
				c.status[target.Name()] = STATE_COMPLETED
				log.Debugf("upload of chunk [%s] from [%s] completed", c.Name, target.Name())
			}
		}
		wg.Wait()
	}
}

func (m *Mosaic) Upload(filename string) ([]*Chunk, error) {
	chunks := m.getMosaic()
	return m.upload(chunks, filename)
}

func (m *Mosaic) UploadWithChunks(chunks []*Chunk, filename string) error {
	for _, c := range chunks {
		for _, s := range m.ss {
			c.status[s.Name()] = STATE_MISSING
		}
	}
	_, err := m.upload(chunks, filename)
	return err
}

func (m *Mosaic) Download(chunks []*Chunk, filename string) error {
	log.Debugf("downloading [%v] chunks", len(chunks))
	for _, c := range chunks {
		err := m.download(c)
		if err != nil {
			return err
		}
	}
	chunkFiles := []string{}
	for _, c := range chunks {
		chunkFiles = append(chunkFiles, c.filename)
	}
	return dust.MergeFile(chunkFiles, filename)
}