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

type Chunk struct {
	Name        string   `json:"name"`
	Real_size   int      `json:"real_size"`
	Index       int      `json:"int"`
	TargetNames []string `json:"targets"`
	Checksum    string   `json:"checksum"`
	err         error
	filename    string
	status      map[string]string
}

type ChunkOption func(*Chunk)

func WithFilename(filename string) ChunkOption {
	return func(s *Chunk) {
		s.filename = filename
	}
}

func NewChunk(index int, name string, opts ...ChunkOption) *Chunk {
	c := &Chunk{
		TargetNames: []string{},
		Name:        name,
		Index:       index,
		status:      map[string]string{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Chunk) GetFilename() string {
	return c.filename
}