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
	"testing"
)

const (
	enigma_phrase = "Kirinuki is a secure password management software by LosAngeles971"
	enigma_email  = "losangeles971@gmail.com"
)

func TestHash(t *testing.T) {
	e := newEnigma()
	h := e.hash([]byte(enigma_phrase))
	t.Log(h)
}

func TestEncryptionWithMainKey(t *testing.T) {
	e := newEnigma(withMainkey(enigma_email, enigma_phrase))
	plaintext := []byte(enigma_phrase)
	encrypted, err := e.encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := e.decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}

func TestEncryptionWithRandomkey(t *testing.T) {
	e := newEnigma(withRandomkey())
	plaintext := []byte(enigma_phrase)
	encrypted, err := e.encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := e.decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}

func TestEncryptionWithEncodedkey(t *testing.T) {
	e1 := newEnigma(withRandomkey())
	key := e1.getEncodedKey()
	plaintext := []byte(enigma_phrase)
	encrypted, err := e1.encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	e2 := newEnigma(withEncodedkey(key))
	decrypted, err := e2.decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	text2 := string(decrypted)
	if string(plaintext) != text2 {
		t.FailNow()
	}
}
