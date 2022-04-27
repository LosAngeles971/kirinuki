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
package toc

import (
	_ "embed"
	"os"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/sirupsen/logrus"
)

func TestTOC(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	base := os.TempDir() + "/toc"
	_ = os.Mkdir(base, os.ModePerm)
	
	toc, err :=	New()
	if err != nil {
		t.Fatal(err)
	}
	k := kirinuki.NewKirinuki("test", kirinuki.WithRandomkey())
	ok := toc.Add(k)
	if !ok {
		t.Fatal("failed add")
	}
	tocFile := base + "/" + mosaic.GetFilename(24)
	err = toc.Save(tocFile)
	if err != nil {
		t.Fatal(err)
	}
	toc2, err := New(WithFilename(tocFile))
	if err != nil {
		t.Fatal(err)
	}
	ok = toc2.Exist("test")
	if !ok {
		t.Fatalf("missed %s", "test")
	}
	os.RemoveAll(base)
}