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
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/sirupsen/logrus"
)

func TestUpload(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/chunk"
	_ = os.Mkdir(base, os.ModePerm)
	upload := NewChunk(1, "test")
	upload.filename = base + "/upload1"
	err := ioutil.WriteFile(upload.filename, test_file1, 0755)
	if err != nil {
		t.Fatal(err)
	}
	target, err := storage.NewStowStorage("mosaic", storage.ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	upload.upload(target)
	if upload.err != nil {
		t.Fatal(upload.err)
	}
	hh, err := enigma.GetFileHash(upload.filename)
	if err != nil {
		t.Fatal(err)
	}
	if hh != upload.Checksum {
		t.Fatalf("mismatch %s - %s", upload.Checksum, hh)
	}
	os.RemoveAll(base + "/")
}

func TestDownload(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/chunk"
	_ = os.Mkdir(base, os.ModePerm)
	tFile := base + "/download1"
	err := ioutil.WriteFile(tFile, test_file1, 0755)
	if err != nil {
		t.Fatal(err)
	}
	download := NewChunk(1, "download1")
	download.Checksum = enigma.GetHash(test_file1)
	download.filename = base + "/download2"
	target, err := storage.NewStowStorage("mosaic", storage.ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	download.download(target)
	if download.err != nil {
		t.Fatal(download.err)
	}
	os.RemoveAll(base + "/")
}
