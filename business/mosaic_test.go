/*
 * Created on Sat Apr 09 2022
 * Author @LosAngeles971
 *
 * The MIT License (MIT)
 * Copyright (c) 2022 @LosAngeles971
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software
 * and associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
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
