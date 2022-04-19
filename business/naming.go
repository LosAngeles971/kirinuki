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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type naming struct {
	chunk_name_size int
}

func newNaming() naming {
	return naming{
		chunk_name_size: 48,
	}
}


// getNameForTOCChunk generates a name for a TOC's chunk
func (n naming) getNameForTOCChunk(session *Session, index int) string {
	data := sha256.Sum256([]byte(fmt.Sprintf("%s_%s_%v", session.GetEmail(), session.GetPassword(), index)))
	return hex.EncodeToString(data[:])
}


//getNameForChunk generates a name for a generic chunk
func (n naming) getNameForChunk() string {
	dd := getRndBytes(n.chunk_name_size/2)
	return hex.EncodeToString(dd)
}
