package helpers

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