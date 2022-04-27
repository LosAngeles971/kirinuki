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
	_ "embed"
	"io/ioutil"
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/sirupsen/logrus"
)

//go:embed test_file1.png
var test_file1 []byte

func TestKirinuki(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/kirinuki"
	_ = os.Mkdir(base, os.ModePerm)
	target, err := storage.NewStowStorage("kirinuki", storage.ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	sFile := base + "/tobe_uploaded"
	tFile := base + "/tobe_downloaded.png"
	err = ioutil.WriteFile(sFile, test_file1, 0755)
	if err != nil {
		t.Fatal(err)
	}
	h1, err := enigma.GetFileHash(sFile)
	if err != nil {
		t.Fatal(err)
	}
	kk := NewKirinuki("test", WithRandomkey())
	err = kk.Upload(sFile, []storage.Storage{target})
	if err != nil {
		t.Fatal(err)
	}
	err = kk.Download(tFile, []storage.Storage{target})
	if err != nil {
		t.Fatal(err)
	}
	h2, err := enigma.GetFileHash(tFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("mismatch %s - %s", h1, h2)
	}
	os.RemoveAll(base + "/")
}
