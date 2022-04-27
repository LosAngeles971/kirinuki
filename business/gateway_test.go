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
package business

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/business/toc"
	"github.com/sirupsen/logrus"
)

const (
	test_files    = 10
	test_email    = "losangeles971@gmail.com"
	test_password = "losangeles971@gmail.com"
)

func TestSession(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/session"
	_ = os.Mkdir(base, os.ModePerm)
	sm, err := storage.NewStorageMap()
	if err != nil {
		t.Fatal(err)
	}
	err = sm.Add("session", storage.ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := New(test_email, test_password, WithStorage(sm), WithTemp(base))
	if err != nil {
		t.Fatalf("failed to create session [%v]", err)
	}
	g.toc, err = toc.New()
	if err != nil {
		t.Fatal(err)
	}
	kName := "test"
	k := kirinuki.NewKirinuki(kName, kirinuki.WithRandomkey())
	ok := g.toc.Add(k)
	if !ok {
		t.Fatal("missed add")
	}
	err = g.Logout()
	if err != nil {
		t.Fatal(err)
	}
	err = g.Login()
	if err != nil {
		t.Fatalf("failed login [%v]", err)
	}
	if !g.isOpen() {
		t.Fatal("inconsistent state")
	}
	if !g.toc.Exist(kName) {
		t.Fatalf("missing kFile [%s]", kName)
	}
	os.RemoveAll(base)
}

func TestPutGet(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/gateway"
	_ = os.Mkdir(base, os.ModePerm)
	sm, err := storage.NewStorageMap()
	if err != nil {
		t.Fatal(err)
	}
	err = sm.Add("gateway", storage.ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": base,
		},
	})
	if err != nil {
		t.Fatalf("failed storage adding -> %v", err)
	}
	g, err := New(test_email, test_password, WithStorage(sm))
	if err != nil {
		t.Fatalf("failed gateway creation -> %v", err)
	}
	err = g.CreateTableOfContent()
	if err != nil {
		t.Fatalf("failed toc creation -> %v", err)
	}
	err = g.Logout()
	if err != nil {
		t.Fatalf("failed logout -> %v", err)
	}
	for i := 0; i < test_files; i++ {
		size := 50000
		data := make([]byte, size)
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			panic(err.Error())
		}
		checksum := enigma.GetHash(data)
		fName := base + "/" + mosaic.GetFilename(24)
		name := fmt.Sprintf("testfile%v", i)
		err = ioutil.WriteFile(fName, data, 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = g.Login()
		if err != nil {
			t.Fatal(err)
		}
		err = g.Upload(fName, name, false)
		if err != nil {
			t.Errorf("failed upload %s due to %v", name, err)
		}
		err = g.Logout()
		if err != nil {
			t.Fatal(err)
		}
		err = g.Login()
		if err != nil {
			t.Fatal(err)
		}
		dName := base + "/" + mosaic.GetFilename(24)
		err = g.Download(name, dName)
		if err != nil {
			t.Errorf("failed download %s due to %v", name, err)
		}
		dChecksum, err := enigma.GetFileHash(dName)
		if err != nil {
			t.Fatal(err)
		}
		if checksum != dChecksum {
			t.Fatalf("rebuild failed, expected hash [%v] not [%v]", checksum, dChecksum)
		}
	}
	err = g.Login()
	if err != nil {
		t.Fatal(err)
	}
	n, err := g.Size()
	if err != nil {
		t.Fatal(err)
	}
	if n != test_files {
		t.Errorf("expected %v Kirinuki files not %v", test_files, n)
	}
	for i := 0; i < test_files; i++ {
		name := fmt.Sprintf("test_file%v", i)
		ok, err := g.Exist(name)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("expected one found entry for %v", name)
		}
	}
	os.RemoveAll(base)
}
