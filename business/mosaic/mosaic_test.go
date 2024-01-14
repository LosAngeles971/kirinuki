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

import (
	_ "embed"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/stretchr/testify/require"
)

//go:embed test_file1.png
var test_file1 []byte

// TestMosaic: it verifies upload and download of Kirinuki files
func TestMosaic(t *testing.T) {
	tsm := storage.NewTestLocalMultistorage("mosaic")
	sChunk := NewChunk(1, "file", WithFilename(storage.GetTmp() + "/tobe_uploaded"))
	tChunk := NewChunk(1, "file", WithFilename(storage.GetTmp() + "/tobe_downloaded"))
	err := ioutil.WriteFile(sChunk.filename, test_file1, 0755)
	require.Nil(t, err)
	h1, err := storage.GetFileHash(sChunk.filename)
	require.Nil(t, err)
	mm := New(tsm.GetMultiStorage())
	err = mm.Upload([]*Chunk{sChunk})
	require.Nil(t, err)
	err = mm.Download([]*Chunk{tChunk})
	require.Nil(t, err)
	h2, err := storage.GetFileHash(tChunk.filename)
	require.Nil(t, err)
	require.Equal(t, h1, h2)
	tsm.Clean()
}
