package kirinuki

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
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/stretchr/testify/require"
)

func TestSplitMerge(t *testing.T) {
	storage.SetTestEnv()
	file := NewFile("split-merge")
	file.Chunks = []*mosaic.Chunk{
		mosaic.NewChunk(1, "c1", mosaic.WithFilename(storage.GetTmp()+"/split1")),
		mosaic.NewChunk(1, "c2", mosaic.WithFilename(storage.GetTmp()+"/split2")),
		mosaic.NewChunk(1, "c3", mosaic.WithFilename(storage.GetTmp()+"/split3")),
	}
	splitFile := storage.GetTmp() + "/split.png"
	mergeFile := storage.GetTmp() + "/merge.png"
	err := storage.CreateFile(splitFile, 100000)
	require.Nil(t, err)
	h1, _ := storage.GetFileHash(splitFile)
	err = file.Split(splitFile)
	require.Nil(t, err)
	err = file.Merge(mergeFile)
	require.Nil(t, err)
	h2, err := storage.GetFileHash(mergeFile)
	require.Nil(t, err)
	require.Equal(t, h1, h2)
	storage.CleanTestEnv()
}

func TestConfidentiality(t *testing.T) {
	storage.SetTestEnv()
	size := 50000
	data := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, data)
	require.Nil(t, err)
	checksum := storage.GetHash(data)
	sFile := storage.GetTmp() + "/plain.png"
	err = ioutil.WriteFile(sFile, data, 0755)
	require.Nil(t, err)
	f := NewFile("plain", WithRandomkey())
	tFile := storage.GetTmp() + "/crypted.png"
	ttFile := storage.GetTmp() + "/decrypted.png"
	h1, err := storage.GetFileHash(sFile)
	require.Nil(t, err)
	require.Equal(t, checksum, h1)
	err = f.Encrypt(sFile, tFile)
	require.Nil(t, err)
	err = f.Decrypt(tFile, ttFile)
	require.Nil(t, err)
	h2, err := storage.GetFileHash(ttFile)
	require.Nil(t, err)
	require.Equal(t, h1, h2)
}

func TestIO(t *testing.T) {
	storage.SetTestEnv()
	name := "source"
	fName := storage.GetTmp() + "/" + name
	err := storage.CreateFile(fName, 100000)
	require.Nil(t, err)
	checksum, _ := storage.GetFileHash(fName)
	tsm := storage.NewTestLocalMultistorage("kirinuki")
	file := NewFile(name)
	err = file.Upload(fName, tsm.GetMultiStorage())
	require.Nil(t, err)
	destFile := storage.GetTmp() + "/dest.png" 
	err = file.Download(destFile, tsm.GetMultiStorage())
	require.Nil(t, err)
	h, err := storage.GetFileHash(destFile)
	require.Nil(t, err)
	require.Equal(t, checksum, h)
	storage.CleanTestEnv()
	tsm.Clean()
}
