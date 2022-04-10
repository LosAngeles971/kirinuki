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

type k_data_test struct {
	name     string
	data     []byte
	checksum string
}

//go:embed test_file1.png
var test_file1 []byte

var k_data_tests []k_data_test = []k_data_test{
	{
		name: "test1",
		data: test_file1,
	},
}

// TestMosaic verifies upload and download of Kirinuki files
func TestMosaicWithoutEncryption(t *testing.T) {
	sm, err := storage.NewStorageMap(storage.WithTemp())
	if err != nil {
		t.Fatal(err)
	}
	ee := newEnigma()
	for _, tt := range k_data_tests {
		tt.checksum = ee.hash(tt.data) 
		k1 := NewKirinuki(tt.name)
		if k1.Encryption {
			t.Fatalf("encryption should be off [%v]", k1.Encryption)
		}
		k2 := NewKirinuki(tt.name, WithRandomkey())
		if !k2.Encryption {
			t.Fatalf("encryption should be on [%v]", k2.Encryption)
		}
		err := k1.addData(tt.data)
		if err != nil {
			t.Fatal(err)
		}
		err = k2.addData(tt.data)
		if err != nil {
			t.Fatal(err)
		}
		err = putKiriuki(k1, sm.Array())
		if err != nil {
			t.Fatal(err)
		}
		err = putKiriuki(k2, sm.Array())
		if err != nil {
			t.Fatal(err)
		}
		d1, err := getKirinuki(k1, sm.Array())
		if err != nil {
			t.Fatal(err)
		}
		d2, err := getKirinuki(k2, sm.Array())
		if err != nil {
			t.Fatal(err)
		}
		if len(tt.data) != len(d1) {
			t.Fatalf("[no encryption] wrong size, expected %v not %v", len(tt.data), len(d1))
		}
		if len(tt.data) != len(d2) {
			t.Fatalf("[encryption] wrong size, expected %v not %v", len(tt.data), len(d2))
		}
		ck1 := ee.hash(d1)
		ck2 := ee.hash(d2)
		if tt.checksum != ck1 {
			t.Fatalf("[no encryption] rebuild failed, expected hash [%v] not [%v]", tt.checksum, ck1)
		}
		if tt.checksum != ck2 {
			t.Fatalf("[encryption] rebuild failed, expected hash [%v] not [%v]", tt.checksum, ck2)
		}
	}
}
