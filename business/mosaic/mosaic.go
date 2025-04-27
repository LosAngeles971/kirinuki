package mosaic

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

import (
	"fmt"
	"sync"

	"github.com/LosAngeles971/kirinuki/business/multistorage"
	"github.com/LosAngeles971/kirinuki/business/config"
	"github.com/LosAngeles971/kirinuki/business/helpers"
	log "github.com/sirupsen/logrus"
)

const (
	STATE_COMPLETED = 0
	STATE_FAILED    = 1
	STATE_MISSING   = 2
)

type MosaicOption func(*Mosaic)

// Mosaic provides abstraction to file CRUD methods, 
// handling the underlying complexity of splitting and rebuilding of file's chunks
// over distributed storage
type Mosaic struct {
	ms          *multistorage.MultiStorage // storage system used to handle Kirinuki files
	max_threads int                   // maximum number of parallel threads for upload/download methods
}

func New(ms *multistorage.MultiStorage, opts ...MosaicOption) *Mosaic {
	m := &Mosaic{
		ms:          ms,
		max_threads: 4,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (m *Mosaic) getTarget(chunks []*Chunk) (*Chunk, string) {
	for _, c := range chunks {
		for _, nn := range m.ms.Names() {
			status, ok := c.status[nn]
			if ok && status == STATE_MISSING {
				return c, nn
			}
		}
	}
	return nil, ""
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

func (m *Mosaic) uploadChunk(c *Chunk, sName string) {
	c.err = nil
	log.Debugf("uploading of chunk %s from file %s ...", c.Name, c.filename)
	c.Checksum, c.err = helpers.GetFileHash(c.filename)
	if c.err == nil {
		c.err = m.ms.Upload(sName, c.filename, c.Name)
	}
}

func (m *Mosaic) download(chunk *Chunk) error {
	for _, sName := range m.ms.Names() {
		log.Debugf("downloading of chunk %s to %s ...", chunk.Name, chunk.filename)
		ck, err := m.ms.Download(sName, chunk.Name, chunk.filename)
		if err == nil {
			if len(chunk.Checksum) > 0 {
				if chunk.Checksum == ck {
					return nil
				} else {
					log.Errorf("failed download chunk %s from %s -> expected hash %s not %s", chunk.Name, sName, chunk.Checksum, ck)
				}
			} else {
				// checksum check only if c.Checksum is set
				// Indeed, Table Of Content cannot have c.Checksum set
				return nil
			}
		} else {
			log.Errorf("failed download chunk %s from %s -> %v", chunk.Name, sName, chunk.err)
		}
	}
	return fmt.Errorf("failed to download chunk %s from all targets", chunk.Name)
}

func (m *Mosaic) Upload(chunks []*Chunk) error {
	log.Debugf("uploading [%v] chunks", len(chunks))
	for _, c := range chunks {
		for _, nn := range m.ms.Names() {
			c.status[nn] = STATE_MISSING
		}
	}
	for {
		end, completed := m.isComplete(chunks)
		if end {
			if completed {
				return nil
			} else {
				return fmt.Errorf("failed to upload chunks")
			}
		}
		var wg sync.WaitGroup
		for i := 0; i < m.max_threads; i++ {
			c, sName := m.getTarget(chunks)
			if c != nil && sName != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					m.uploadChunk(c, sName)
				}()
				if c.err != nil {
					return c.err
				}
				c.status[sName] = STATE_COMPLETED
				log.Debugf("upload of chunk [%s] from [%s] completed", c.Name, sName)
			}
		}
		wg.Wait()
	}
}

func (m *Mosaic) Download(chunks []*Chunk) error {
	log.Debugf("downloading [%v] chunks", len(chunks))
	for _, c := range chunks {
		if c.filename == "" {
			c.filename = config.GetTmp() + "/" + c.Name
		}
		err := m.download(c)
		if err != nil {
			return err
		}
	}
	return nil
}
