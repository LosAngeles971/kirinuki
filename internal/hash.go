package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
)

type StreamHash struct {
	r io.Reader
	h hash.Hash
	t io.Reader
}

func NewStreamHash(r io.Reader) StreamHash {
	s := StreamHash{
		r: r,
		h: sha256.New(),
	}
	s.t = io.TeeReader(s.r, s.h)
	return s
}

func (s StreamHash) GetReader() io.Reader {
	return s.t
}

func (s StreamHash) GetHash() string {
	return hex.EncodeToString(s.h.Sum(nil))
}