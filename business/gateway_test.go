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
	"testing"

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/internal"
)

const (
	test_files    = 1
)

// TestSession tests the creation, alteration, storing and re-opening of a session
func TestSession(t *testing.T) {
	internal.Setup()
	sm := storage.GetTmp("session")
	g, err := New(internal.Test_email, internal.Test_password, WithStorage(sm))
	if err != nil {
		t.Fatal(err)
	}
	g.SetEmptyTableOfContent()
	kName := "test"
	k := kirinuki.NewFile(kName, kirinuki.WithRandomkey())
	// Adding something to the TableOfContent
	if !g.toc.Add(k) {
		t.Fatal("cannog add a kirinuki file")
	}
	// Saving the TableOfContent
	err = g.Logout()
	if err != nil {
		t.Fatalf("failed to logout -> %v", err)
	}
	// Opening the TableOfContent
	err = g.Login()
	if err != nil {
		t.Fatalf("failed re-login [%v]", err)
	}
	if !g.isOpen() {
		t.Fatal("inconsistent state")
	}
	// Check if the kirinuki file is still there into the TableOfContent
	if !g.toc.Exist(kName) {
		t.Fatalf("missing kFile [%s]", kName)
	}
	internal.Clean("session")
}

func TestIO(t *testing.T) {
	internal.Setup()
	sm := storage.GetTmp("gateway")
	g, err := New(internal.Test_email, internal.Test_password, WithStorage(sm))
	if err != nil {
		t.Fatal(err)
	}
	g.SetEmptyTableOfContent()
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
		checksum := internal.GetHash(data)
		name := fmt.Sprintf("testfile%v", i)
		fName := internal.GetTmp() + "/" + name
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
		dName := internal.GetTmp() + "/" + fmt.Sprintf("d_testfile%v", i)
		err = g.Download(name, dName)
		if err != nil {
			t.Fatalf("failed download %s to local filename %s -> %v", name, dName, err)
		}
		dChecksum, err := internal.GetFileHash(dName)
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
		name := fmt.Sprintf("testfile%v", i)
		f, err := g.Exist(name)
		if err != nil {
			t.Fatal(err)
		}
		if f == nil {
			t.Errorf("expected one found entry for %v", name)
		}
	}
	internal.Clean("gateway")
}
