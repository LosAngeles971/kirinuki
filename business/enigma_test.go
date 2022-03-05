/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Testing hashing and encryption features

This testing session DOES NOT need external data, neither it interacts with the system


+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"testing"
)

const (
	enigma_phrase = "Kirinuki is a secure password management software by LosAngeles971"
	enigma_email  = "losangeles971@gmail.com"
)

func TestHash(t *testing.T) {
	e := NewEnigma()
	h := e.hash([]byte(enigma_phrase))
	t.Log(h)
}

func TestEncryption(t *testing.T) {
	e := NewEnigma(WithMainkey(enigma_email, enigma_phrase))
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
