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
	"testing"
)

//go:embed test.yml
var yCfgFile string

//go:embed test.json
var jCfgFile string

// TestLoad just verifies the capability of loading storage map from a confifuration file
func TestLoad(t *testing.T) {
	m1, err := NewStorageMap(WithYAMLData([]byte(yCfgFile)))
	if err != nil {
		t.Fatal(err)
	}
	m2, err := NewStorageMap(WithJSONData([]byte(jCfgFile)))
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []*StorageMap{m1, m2,} {
		_, err =  m.Get("local")
		if err != nil {
			t.Fatal(err)
		}
		if m.Size() != 1 {
			t.Fatalf("not expected size of %v", m.Size())
		}
	}
}

func TestTempStorage(t *testing.T) {
	sm, err := NewStorageMap(WithTemp())
	if err != nil {
		t.Fatal(err)
	}
	if sm.Size() != 1 {
		t.Fatalf("storage array should have the %s storage", TEMP_STORAGE)
	}
	ss := sm.Array()
	if len(ss) != 1 {
		t.Fatalf("wrong array size %v", len(ss))
	}
}
