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
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

func TestTOC(t *testing.T) {
	sm, err := storage.NewStorageMap()
	if err != nil {
		t.Fatal(err)
	}
	sm.Add("test", storage.ConfigItem{
		Type: "filesystem",
		Cfg: map[string]string{
			"path": os.TempDir(),
		},
	})
	session, err := NewSession(test_email, test_password, true, WithStorage(sm))
	if err != nil {
		t.Fatalf("failed to create a session due to %v", err)
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
		k, err := NewKirinuki(WithKirinukiData(tt.name, tt.data))
		if err != nil {
			t.Fatal(err)
		}
		ok := toc.Add(k)
		if !ok {
			t.Fatal("File not added")
		}
		if !toc.Find(tt.name) {
			t.Fatalf("toc does not contain kirinuki %s", tt.name)
		}
	}
	err = session.logout()
	if err != nil {
		t.Fatalf("failed to logout from session due to %v", err)
	}
	err = session.login()
	if err != nil {
		t.Fatalf("failed to open a session due to %v", err)
	}
	toc2, err := session.getTOC()
	if err != nil {
		t.Fatalf("failed to get toc (2) from session due to %v", err)
	}
	for _, tt := range k_data_tests {
		if !toc2.Find(tt.name) {
			t.Fatalf("reloaded toc does not contain kirinuki %s", tt.name)
		}
	}
}