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
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

// Session includes the mandatory set of information to work with Kirinuki files
type Session struct {
	email           string
	password        string
	chunksForTOC    int
	chunk_name_size int
	toc             *TOC
	storage         *storage.StorageMap
}

type SessionOption func(*Session)

func WithStorage(m *storage.StorageMap) SessionOption {
	return func(s *Session) {
		s.storage = m
	}
}

func NewSession(email string, password string, opts ...SessionOption) (*Session, error) {
	s := &Session{
		email:        email,
		chunksForTOC: 3,
		password:     password,
		chunk_name_size: 48,
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

func (s *Session) GetEmail() string {
	return s.email
}

func (s *Session) GetPassword() string {
	return s.password
}

// getNameForTOCChunk generates the number of chunks and their names for a generic file
func (s *Session) getChunks(file []byte) []string {
	names := []string{}
	n := 11
	size := len(file)
	if size < 1000000 {
		n = 9
	}
	if size < 100000 {
		n = 7
	}
	if size < 10000 {
		n = 5
	}
	if size < 1000 {
		n = 3
	}
	for i := 0; i <= n; i++ {
		dd := getRndBytes(s.chunk_name_size / 2)
		names = append(names, hex.EncodeToString(dd))
	}
	return names
}

// getNameForTOCChunk generates the number of chunks and their names for the TOC
func (s *Session) getChunksForTOC() []string {
	names := []string{}
	e := newEnigma()
	for i := 0; i < 3; i++ {
		name := e.hash([]byte(fmt.Sprintf("%s_%s_%v", s.email, s.password, i)))
		names = append(names, name)
	}
	return names
}

func (s *Session) createTableOfContent() error {
	var err error
	s.toc, err = newTOC()
	if err != nil {
		return err
	}
	return s.logout()
}

// getTableOfContent returns the empty Kirinuki file for an existend or new TOC
func (s *Session) getTableOfContent() (*Kirinuki, error) {
	key := newEnigma(withMainkey(s.email, s.password)).getEncodedKey()
	k := NewKirinuki("toc", s.getChunksForTOC(), WithEncodedKey(key))
	return k, nil
}

func (s *Session) login() error {
	if s.isOpen() {
		return nil
	}
	k, err := s.getTableOfContent()
	if err != nil {
		return err
	}
	data, err := getKirinuki(k, s.storage.Array())
	if err != nil {
		return err
	}
	s.toc, err = newTOC(TOCWithData(data))
	return err
}

// logout saves the current open TOC on cloud and closes the session
func (s *Session) logout() error {
	if !s.isOpen() {
		log.Errorf("session %s is already closed", s.email)
		return nil
	}
	k, err := s.getTableOfContent()
	if err != nil {
		return err
	}
	k.Replicas = s.storage.Size()
	tocdata, err := json.Marshal(s.toc)
	if err != nil {
		return err
	}
	err = k.addData(tocdata)
	if err != nil {
		return err
	}
	err = putKiriuki(k, s.storage.Array())
	if err != nil {
		return err
	}
	s.toc = nil
	return nil
}

func (s *Session) getTOC() (*TOC, error) {
	if !s.isOpen() {
		return nil, fmt.Errorf("session is not open")
	}
	return s.toc, nil
}

func (s *Session) isOpen() bool {
	return s.toc != nil
}
