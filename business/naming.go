/*
 * Created on Sat Apr 09 2022
 * Author @LosAngeles971
 *
 * The MIT License (MIT)
 * Copyright (c) 2022 @LosAngeles971
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software
 * and associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
 * TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
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


// getNameForTOCChunk generates a name for a TOC's chunk
func (n naming) getNameForTOCChunk(session *Session, index int) string {
	data := sha256.Sum256([]byte(fmt.Sprintf("%s_%s_%v", session.GetEmail(), session.GetPassword(), index)))
	return base64.StdEncoding.EncodeToString(data[:])
}


//getNameForChunk generates a name for a generic chunk
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
