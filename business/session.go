/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

This module is in charge of handling the Kirinuki's session

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"fmt"

	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

type Session struct {
	email        string
	password     string
	chunksForTOC int
	toc          *TOC
	storage      *storage.StorageMap
}

type SessionOption func(*Session)

func WithStorage(m *storage.StorageMap) SessionOption {
	return func(s *Session) {
		s.storage = m
	}
}

func NewSession(email string, password string, scratch bool, opts ...SessionOption) (*Session, error) {
	s := &Session{
		email:        email,
		chunksForTOC: 3,
		password:     password,
	}
	if scratch {
		var err error
		s.toc, err = newTOC()
		if err != nil {
			return nil, err
		}
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

func (s *Session) getTOCSkeleton() (*Kirinuki, error) {
	k, err := NewKirinuki()
	if err != nil {
		return nil, err
	}
	k.Name = NewEnigma().hash([]byte(s.email))
	k.Chunks = []*chunk{}
	for index := 0; index < s.chunksForTOC; index++ {
		ch, err := newChunk(index, withChunkName(newNaming().getNameForTOCChunk(s, index)))
		if err != nil {
			return nil, err
		}
		k.Chunks = append(k.Chunks, ch)
	}
	return k, nil
}

func (s *Session) login() error {
	if s.isOpen() {
		return nil
	}
	k, err := s.getTOCSkeleton()
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

func (s *Session) kill() {
	s.toc = nil
}

func (s *Session) logout() error {
	if !s.isOpen() {
		log.Errorf("session with email %s is already closed", s.email)
	}
	k := Kirinuki{}
	k.Name = NewEnigma().hash([]byte(s.email))
	k.Replicas = s.storage.Size()
	tocdata, err := s.toc.serialize()
	if err != nil {
		return err
	}
	datas, err := splitFile(tocdata, s.chunksForTOC)
	if err != nil {
		return err
	}
	k.Chunks = []*chunk{}
	for index, data := range datas {
		ch, err := newChunk(index, withChunkData(data), withChunkName(newNaming().getNameForTOCChunk(s, index)))
		if err != nil {
			return err
		}
		k.Chunks = append(k.Chunks, ch)
	}
	return putKiriuki(&k, s.storage.Array())
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