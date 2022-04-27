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
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/business/toc"
	"github.com/sirupsen/logrus"
)

const (
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
	s, err := NewSession(test_email, test_password, WithStorage(sm), WithTemp(base))
	if err != nil {
		t.Fatalf("failed to create session [%v]", err)
	}
	s.toc, err = toc.New()
	if err != nil {
		t.Fatal(err)
	}
	k := kirinuki.NewKirinuki("test", kirinuki.WithRandomkey())
	ok := s.toc.Add(k)
	if !ok {
		t.Fatal("missed add")
	}
	err = s.logout()
	if err != nil {
		t.Fatal(err)
	}
	err = s.login()
	if err != nil {
		t.Fatalf("failed login [%v]", err)
	}
	if !s.isOpen() {
		t.Fatal("inconsistent state")
	}
	os.RemoveAll(base)
}
