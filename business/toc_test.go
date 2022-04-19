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
	_ "embed"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

//go:embed storage/minio.json
var minio []byte

func getStorage() (*storage.StorageMap, error) {
	return storage.NewStorageMap(storage.WithJSONData(minio))
}

func TestTOC(t *testing.T) {
	sm, err := getStorage()
	if err != nil {
		t.Fatal(err)
	}
	session, err := NewSession(test_email, test_password, WithStorage(sm))
	if err != nil {
		t.Fatalf("failed to create a session due to %v", err)
	}
	err = session.createTableOfContent()
	if err != nil {
		t.Fatalf("failed to create new table of content %v", err)
	}
	err = session.login()
	if err != nil {
		t.Fatalf("failed to open a session due to %v", err)
	}
	toc, err := session.getTOC()
	if err != nil {
		t.Fatalf("failed to get toc from session due to %v", err)
	}
	for _, tt := range k_data_tests {
		k := NewKirinuki(tt.name)
		err := k.addData(tt.data)
		if err != nil {
			t.Fatal(err)
		}
		ok := toc.add(k)
		if !ok {
			t.Fatal("File not added")
		}
		if !toc.exist(tt.name) {
			t.Fatalf("toc does not contain kirinuki %s", tt.name)
		}
	}
	err = session.logout()
	if err != nil {
		t.Fatalf("cannot logout [%v]", err)
	}
	err = session.login()
	if err != nil {
		t.Fatalf("cannot login  [%v]", err)
	}
	toc2, err := session.getTOC()
	if err != nil {
		t.Fatalf("failed to get toc (2) from session due to %v", err)
	}
	for _, tt := range k_data_tests {
		if !toc2.exist(tt.name) {
			t.Fatalf("reloaded toc does not contain kirinuki %s", tt.name)
		}
	}
}