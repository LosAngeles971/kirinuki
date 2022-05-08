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

import (
	_ "embed"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/internal"
)

//go:embed test_file1.png
var test_file1 []byte

// TestMosaic verifies upload and download of Kirinuki files
func TestMosaic(t *testing.T) {
	internal.Setup()
	sm := storage.GetTmp("mosaic")
	sChunk := NewChunk(1, "file", WithFilename(internal.GetTmp() + "/tobe_uploaded"))
	tChunk := NewChunk(1, "file", WithFilename(internal.GetTmp() + "/tobe_downloaded"))
	err := ioutil.WriteFile(sChunk.filename, test_file1, 0755)
	if err != nil {
		t.Fatal(err)
	}
	h1, err := internal.GetFileHash(sChunk.filename)
	if err != nil {
		t.Fatal(err)
	}
	mm := New(sm)
	err = mm.Upload([]*Chunk{sChunk})
	if err != nil {
		t.Fatalf("failed upload -> %v", err)
	}
	err = mm.Download([]*Chunk{tChunk})
	if err != nil {
		t.Fatalf("failed download -> %v", err)
	}
	h2, err := internal.GetFileHash(tChunk.filename)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("mismatch %s %s", h1, h2)
	}
	internal.Clean("mosaic")
}
