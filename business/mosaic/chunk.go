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

// Chunk: it represents a portion of a Kirinuki files
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

// Setting the filename which contain the chunk's data
func WithFilename(filename string) ChunkOption {
	return func(s *Chunk) {
		s.filename = filename
	}
}

// NewChunk: it creates a new chunk with the index and name
func NewChunk(index int, name string, opts ...ChunkOption) *Chunk {
	c := &Chunk{
		TargetNames: []string{},          // list of storages where the chunk's data is stored on
		Name:        name,                // unique identifier of the chunk
		Index:       index,               // chunk's position into the sequence of chunks which compose an entire Kirinuki file
		status:      map[string]string{}, // a map of where the chunk's data has been successufully (or not) uploaded
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// GetFilename: it returns the filename of chunk's data
func (c *Chunk) GetFilename() string {
	return c.filename
}
