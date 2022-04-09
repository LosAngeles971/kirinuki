/*
 * Created on Fri Apr 08 2022
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
