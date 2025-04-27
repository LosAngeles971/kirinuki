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
	"crypto/rand"
	"io"
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/multistorage"
	"github.com/LosAngeles971/kirinuki/business/helpers"
	"github.com/LosAngeles971/kirinuki/business/config"
	"github.com/stretchr/testify/require"
)

func TestSplitMerge(t *testing.T) {
	multistorage.SetTestEnv()
	file := NewFile("split-merge")
	file.Chunks = []*mosaic.Chunk{
		mosaic.NewChunk(1, "c1", mosaic.WithFilename(config.GetTmp()+"/split1")),
		mosaic.NewChunk(1, "c2", mosaic.WithFilename(config.GetTmp()+"/split2")),
		mosaic.NewChunk(1, "c3", mosaic.WithFilename(config.GetTmp()+"/split3")),
	}
	splitFile := config.GetTmp() + "/split.png"
	mergeFile := config.GetTmp() + "/merge.png"
	err := helpers.CreateRandomFile(splitFile, 100000)
	require.Nil(t, err)
	h1, _ := helpers.GetFileHash(splitFile)
	err = file.Split(splitFile)
	require.Nil(t, err)
	err = file.Merge(mergeFile)
	require.Nil(t, err)
	h2, err := helpers.GetFileHash(mergeFile)
	require.Nil(t, err)
	require.Equal(t, h1, h2)
	multistorage.CleanTestEnv()
}

func TestConfidentiality(t *testing.T) {
	multistorage.SetTestEnv()
	size := 50000
	data := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, data)
	require.Nil(t, err)
	checksum := helpers.GetHash(data)
	sFile := config.GetTmp() + "/plain.png"
	err = os.WriteFile(sFile, data, 0755)
	require.Nil(t, err)
	f := NewFile("plain", WithRandomkey())
	tFile := config.GetTmp() + "/crypted.png"
	ttFile := config.GetTmp() + "/decrypted.png"
	h1, err := helpers.GetFileHash(sFile)
	require.Nil(t, err)
	require.Equal(t, checksum, h1)
	err = f.Encrypt(sFile, tFile)
	require.Nil(t, err)
	err = f.Decrypt(tFile, ttFile)
	require.Nil(t, err)
	h2, err := helpers.GetFileHash(ttFile)
	require.Nil(t, err)
	require.Equal(t, h1, h2)
}

func TestIO(t *testing.T) {
	multistorage.SetTestEnv()
	name := "source"
	fName := config.GetTmp() + "/" + name
	err := helpers.CreateRandomFile(fName, 100000)
	require.Nil(t, err)
	checksum, _ := helpers.GetFileHash(fName)
	tsm := multistorage.NewTestLocalMultistorage("kirinuki")
	file := NewFile(name)
	err = file.Upload(fName, tsm.GetMultiStorage())
	require.Nil(t, err)
	destFile := config.GetTmp() + "/dest.png" 
	err = file.Download(destFile, tsm.GetMultiStorage())
	require.Nil(t, err)
	h, err := helpers.GetFileHash(destFile)
	require.Nil(t, err)
	require.Equal(t, checksum, h)
	multistorage.CleanTestEnv()
	tsm.Clean()
}
