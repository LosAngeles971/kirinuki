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
package storage

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed minio.json
var minioFile []byte

//go:embed stow.win.local.json
var localFile []byte

func doTargetStorageTest(t *testing.T, s Storage) {
	hh := GetHash(sftpFile)
	err := s.Put("testfile", sftpFile)
	require.Nil(t, err)
	dd, err := s.Get("testfile")
	require.Nil(t, err)
	require.Equal(t, hh, GetHash(dd))
	f1 := fmt.Sprintf("%s/testfile", os.TempDir())
	h1, err := s.Download("testfile", f1)
	require.Nil(t, err)
	require.Equal(t, hh, h1)
	df, err := ioutil.ReadFile(f1)
	require.Nil(t, err)
	require.Equal(t, GetHash(df), hh)
	err = s.Upload(f1, "testfile2")
	require.Nil(t, err)
	du, err := s.Get("testfile2")
	require.Nil(t, err)
	require.Equal(t, GetHash(du), hh)
}

func TestMinio(t *testing.T) {
	t.Skip("missing integration test")
	sm, err := NewMultiStorage()
	require.Nil(t, err)
	err = sm.LoadByJSON(minioFile)
	require.Nil(t, err)
	s, err := sm.get("minio")
	require.Nil(t, err)
	doTargetStorageTest(t, s)
}

func TestLocal(t *testing.T) {
	sm, err := NewMultiStorage()
	require.Nil(t, err)
	err = sm.LoadByJSON(localFile)
	require.Nil(t, err)
	s, err := sm.get("local")
	require.Nil(t, err)
	doTargetStorageTest(t, s)
}