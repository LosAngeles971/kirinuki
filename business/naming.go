package business

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

type naming struct {
	letterRunes     []rune
	chunk_name_size int
}

func newNaming() naming {
	return naming{
		letterRunes:     []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		chunk_name_size: 48,
	}
}

/*
Name generation for TOC's chunkgs
Name generation uses a not reversible function (sha256), using session data and chunk's index as inputs
*/
func (n naming) getNameForTOCChunk(session *Session, index int) string {
	data := sha256.Sum256([]byte(fmt.Sprintf("%s_%s_%v", session.GetEmail(), session.GetPassword(), index)))
	return base64.StdEncoding.EncodeToString(data[:])
}

/*
Name generation for a generic file's chunkgs
Chunk's name is the result of a random choice of "SIZE_OF_CHUNK_NAME" chars
*/
func (n naming) getNameForChunk() string {
	b := make([]rune, n.chunk_name_size)
	for i := range b {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(n.letterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = n.letterRunes[nBig.Int64()]
	}
	return string(b)
}
