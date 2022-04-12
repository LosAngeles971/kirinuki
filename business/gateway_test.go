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
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestEndurance(t *testing.T) {
	sm, err := storage.NewStorageMap(storage.WithTemp())
	if err != nil {
		t.Fatal(err)
	}
	g, err := New(test_email, test_password, sm)
	if err != nil {
		t.Fatal(err)
	}
	err = g.CreateTableOfContent()
	if err != nil {
		t.Fatal(err)
	}
	tot := 50
	ee := newEnigma()
	for i :=0; i < tot; i++ {
		size := 50000
		data := make([]byte, size)
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			panic(err.Error())
		}
		checksum := ee.hash(data)
		name := fmt.Sprintf("test_file%v", i)
		err = g.Login()
		if err != nil {
			t.Fatal(err)
		}
		err = g.Upload(name, data, true)
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
		back, err := g.Download(name)
		if err != nil {
			t.Errorf("failed download %s due to %v", name, err)
		}
		back_checksum := ee.hash(back)
		if checksum != back_checksum {
			t.Fatalf("rebuild failed, expected hash [%v] not [%v]", checksum, back_checksum)
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
	if n != tot {
		t.Errorf("expected %v Kirinuki files not %v", tot, n)
	}
	for i :=0; i < tot; i++ {
		name := fmt.Sprintf("test_file%v", i)
		ok, err := g.Exist(name)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("expected one found entry for %v", name)
		}
	}
}