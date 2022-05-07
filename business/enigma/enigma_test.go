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
package enigma

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/LosAngeles971/kirinuki/internal"
)

func TestConfidentiality(t *testing.T) {
	internal.Setup()
	size := 50000
	data := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err.Error())
	}
	checksum := GetHash(data)
	sFile := internal.GetTmp() + "/plain.png"
	err := ioutil.WriteFile(sFile, data, 0755)
	if err != nil {
		t.Fatal(err)
	}
	tFile := internal.GetTmp() + "/crypted.png"
	ttFile := internal.GetTmp() + "/decrypted.png"
	h1, err := GetFileHash(sFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != checksum {
		t.Fatalf("file hash different from data hash %s - %s", checksum, h1)
	}
	e := New(WithRandomkey())
	err = e.EncryptFile(sFile, tFile)
	if err != nil {
		t.Fatalf("encryption failed -> %v", err)
	}
	err = e.DecryptFile(tFile, ttFile)
	if err != nil {
		t.Fatalf("decryption failed -> %v", err)
	}
	h2, err := GetFileHash(ttFile)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("mismatch %s %s", h1, h2)
	}
}