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
package business

import (
	"fmt"
	"os"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/business/toc"

	log "github.com/sirupsen/logrus"
)

// Session includes the mandatory set of information to work with Kirinuki files
type Session struct {
	email           string
	password        string
	chunksForTOC    int
	chunk_name_size int
	toc             *toc.TableOfContent
	storage         *storage.StorageMap
	tempDir         string
}

type SessionOption func(*Session)

func WithStorage(m *storage.StorageMap) SessionOption {
	return func(s *Session) {
		s.storage = m
	}
}

func WithTemp(temp string) SessionOption {
	return func(s *Session) {
		s.tempDir = temp
	}
}

func NewSession(email string, password string, opts ...SessionOption) (*Session, error) {
	s := &Session{
		email:           email,
		chunksForTOC:    3,
		password:        password,
		chunk_name_size: 48,
		tempDir:         os.TempDir(),
		toc:             nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.storage == nil {
		var err error
		s.storage, err = storage.NewStorageMap()
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func (s *Session) loadTableOfContent() (*toc.TableOfContent, error) {
	chunks := toc.GetChunks(s.email, s.password, s.storage.Array(), s.tempDir)
	if len(chunks) < 1 {
		return nil, fmt.Errorf("no chunks %v for toc", len(chunks))
	}
	ecTocFile := s.tempDir + "/" + mosaic.GetFilename(24)
	mm := mosaic.New(mosaic.WithStorage(s.storage.Array()))
	err := mm.Download(chunks, ecTocFile)
	if err != nil {
		return nil, err
	}
	ee := enigma.New(enigma.WithMainkey(s.email, s.password))
	tocFile := s.tempDir + "/" +mosaic.GetFilename(24)
	err = ee.DecryptFile(ecTocFile, tocFile)
	if err != nil {
		return nil, err
	}
	toc, err := toc.New(toc.WithFilename(tocFile))
	if err != nil {
		return nil, err
	}
	return toc, nil
}

func (s *Session) login() error {
	if s.isOpen() {
		return nil
	}
	toc, err := s.loadTableOfContent()
	if err != nil {
		return err
	}
	s.toc = toc
	return nil
}

// logout saves the current open TOC on cloud and closes the session
func (s *Session) logout() error {
	if !s.isOpen() {
		log.Errorf("session %s is already closed", s.email)
		return nil
	}
	tocFile := s.tempDir + "/" +mosaic.GetFilename(24)
	err := s.toc.Save(tocFile)
	if err != nil {
		return err
	}
	ecTocFile := s.tempDir + "/" +mosaic.GetFilename(24)
	ee := enigma.New(enigma.WithMainkey(s.email, s.password))
	err = ee.EncryptFile(tocFile, ecTocFile)
	if err != nil {
		return err
	}
	mm := mosaic.New(mosaic.WithStorage(s.storage.Array()))
	chunks := toc.GetChunks(s.email, s.password, s.storage.Array(), s.tempDir)
	if len(chunks) < 1 {
		return fmt.Errorf("no chunks %v for toc", len(chunks))
	}
	err = mm.UploadWithChunks(chunks, ecTocFile)
	if err != nil {
		return err
	}
	s.toc = nil
	return nil
}

func (s *Session) isOpen() bool {
	return s.toc != nil
}
