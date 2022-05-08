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

	"github.com/LosAngeles971/kirinuki/internal"
)

//go:embed minio.json
var minioFile []byte

//go:embed stow.win.local.json
var localFile []byte

func doTest(sName string, sFile []byte, t *testing.T) {
	internal.Setup()
	hh := internal.GetHash(sftpFile)
	sm, err := NewMultiStorage(WithJSONData(sFile))
	if err != nil {
		t.Fatalf("failed storage map init -> %v", err)
	}
	s, err := sm.get(sName)
	if err != nil {
		t.Fatalf("failed to ge storage %s -> %v", sName, err)
	}
	err = s.Put("testfile", sftpFile)
	if err != nil {
		t.Fatalf("failed put file [%v]", err)
	}
	dd, err := s.Get("testfile")
	if err != nil {
		t.Fatalf("failed get file [%v]", err)
	}
	if internal.GetHash(dd) != hh {
		t.Fatalf("expected hash %s not %s", hh, internal.GetHash(dd))
	}
	f1 := fmt.Sprintf("%s/testfile", os.TempDir())
	h1, err := s.Download("testfile", f1)
	if h1 != hh {
		t.Fatalf("mismatch %s %s", h1, hh)
	}
	if err != nil {
		t.Fatalf("failed download file [%v]", err)
	}
	df, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatalf("failed to read file %s -> %v", f1, err)
	}
	if internal.GetHash(df) != hh {
		t.Fatalf("df - expected hash %s not %s", hh, internal.GetHash(df))
	}
	err = s.Upload(f1, "testfile2")
	if err != nil {
		t.Fatalf("failed to upload %s -> %v", f1, err)
	}
	du, err := s.Get("testfile2")
	if err != nil {
		t.Fatalf("failed get file [%v]", err)
	}
	if internal.GetHash(du) != hh {
		t.Fatalf("du - expcted hash %s not %s", hh, internal.GetHash(du))
	}
}

func TestMinio(t *testing.T) {
	doTest("minio", minioFile, t)
}

func TestLocal(t *testing.T) {
	doTest("local", localFile, t)
}