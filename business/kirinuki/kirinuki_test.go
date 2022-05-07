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
package kirinuki

import (
	"crypto/rand"
	_ "embed"
	"io"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/internal"
)

func TestSplitFile(t *testing.T) {
	k := New(storage.GetTmp("split"))
	file := NewKirinuki("split-merge")
	file.Chunks = []*mosaic.Chunk{
		mosaic.NewChunk(1, "c1", mosaic.WithFilename(internal.GetTmp() + "/split1")),
		mosaic.NewChunk(1, "c2", mosaic.WithFilename(internal.GetTmp() + "/split2")),
		mosaic.NewChunk(1, "c3", mosaic.WithFilename(internal.GetTmp() + "/split3")),
	}
	splitFile := internal.GetTmp() + "/split.png" 
	mergeFile := internal.GetTmp() + "/merge.png" 
	err := ioutil.WriteFile(splitFile, test_file1, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = k.splitFile(splitFile, file)
	if err != nil {
		t.Fatal(err)
	}
	err = k.mergeChunks(file, mergeFile)
	if err != nil {
		t.Fatal(err)
	}
	h1, err := enigma.GetFileHash(splitFile)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := enigma.GetFileHash(mergeFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("mismatch %s - %s", h1, h2)
	}
	internal.Clean("split")
}

func TestIO(t *testing.T) {
	internal.Setup()
	size := 50000
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err.Error())
	}
	checksum := enigma.GetHash(data)
	name := "source"
	fName := internal.GetTmp() + "/" + name
	err := ioutil.WriteFile(fName, data, 0755)
	if err != nil {
		t.Fatal(err)
	}
	k := New(storage.GetTmp("kirinuki"), WithTempDir(internal.GetTmp()))
	file := NewKirinuki(name)
	err = k.Upload(fName, file)
	if err != nil {
		t.Fatal(err)
	}
	destFile := internal.GetTmp() + "/dest.png" 
	err = k.Download(file, destFile)
	if err != nil {
		t.Fatalf("download failed -> %v", err)
	}
	h, err := enigma.GetFileHash(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if h != checksum {
		t.Fatalf("wrong hash %s instead of %s", h, checksum)
	}
	internal.Clean("kirinuki")
}