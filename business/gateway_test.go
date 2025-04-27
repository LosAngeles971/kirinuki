package business

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

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/multistorage"
	"github.com/LosAngeles971/kirinuki/business/helpers"
	"github.com/LosAngeles971/kirinuki/business/config"
	"github.com/stretchr/testify/require"
)

const (
	test_files    = 1
)

// TestSession tests the creation, alteration, storing and re-opening of a session
func TestSession(t *testing.T) {
	tsm := multistorage.NewTestLocalMultistorage("gateway")
	g, err := New(multistorage.Test_email, multistorage.Test_password, WithStorage(tsm.GetMultiStorage()))
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
	tsm.Clean()
}

func TestIO(t *testing.T) {
	tsm := multistorage.NewTestLocalMultistorage("gateway")
	g, err := New(multistorage.Test_email, multistorage.Test_password, WithStorage(tsm.GetMultiStorage()))
	require.Nil(t, err)
	g.SetEmptyTableOfContent()
	err = g.Logout()
	require.Nil(t, err)
	for i := 0; i < test_files; i++ {
		size := 50000
		data := make([]byte, size)
		_, err := io.ReadFull(rand.Reader, data)
		require.Nil(t, err)
		checksum := helpers.GetHash(data)
		name := fmt.Sprintf("testfile%v", i)
		fName := config.GetTmp() + "/" + name
		err = ioutil.WriteFile(fName, data, 0755)
		require.Nil(t, err)
		err = g.Login()
		require.Nil(t, err)
		err = g.Upload(fName, name, false)
		require.Nil(t, err)
		err = g.Logout()
		require.Nil(t, err)
		err = g.Login()
		require.Nil(t, err)
		dName := config.GetTmp() + "/" + fmt.Sprintf("d_testfile%v", i)
		err = g.Download(name, dName)
		require.Nil(t, err, fmt.Sprintf("failed download %s to local filename %s -> %v", name, dName, err))
		dChecksum, err := helpers.GetFileHash(dName)
		require.Nil(t, err)
		require.Equal(t, dChecksum, checksum, fmt.Sprintf("rebuild failed, expected hash [%v] not [%v]", checksum, dChecksum))
	}
	err = g.Login()
	require.Nil(t, err)
	n, err := g.Size()
	require.Nil(t, err)
	if n != test_files {
		t.Errorf("expected %v Kirinuki files not %v", test_files, n)
	}
	for i := 0; i < test_files; i++ {
		name := fmt.Sprintf("testfile%v", i)
		f, err := g.Exist(name)
		require.Nil(t, err)
		require.NotNil(t, f)
	}
	tsm.Clean()
}
