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
	_ "embed"
	"io/ioutil"
	"os"
	"testing"
)

const (
	enigma_phrase = "Kirinuki is a secure password management software by LosAngeles971"
	enigma_email  = "losangeles971@gmail.com"
)

func TestEncryptionWithMainKey(t *testing.T) {
	e := New(WithMainkey(enigma_email, enigma_phrase))
	plaintext := []byte(enigma_phrase)
	encrypted, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := e.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}

func TestEncryptionWithRandomkey(t *testing.T) {
	e := New(WithRandomkey())
	plaintext := []byte(enigma_phrase)
	encrypted, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := e.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}

func TestEncryptionWithEncodedkey(t *testing.T) {
	e1 := New(WithRandomkey())
	key := e1.GetEncodedKey()
	plaintext := []byte(enigma_phrase)
	encrypted, err := e1.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	e2 := New(WithEncodedkey(key))
	decrypted, err := e2.Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	base := os.TempDir() + "/enigma"
	_ = os.Mkdir(base, os.ModePerm)
	sFile := base + "/plain.png"
	tFile := base + "/crypted.png"
	ttFile := base + "/decrypted.png"
	err := ioutil.WriteFile(sFile, hFile, 0755)
	if err != nil {
		t.Fatal(err)
	}
	h1, err := GetFileHash(sFile)
	if err != nil {
		t.Fatal(err)
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
	os.Remove(base)
}