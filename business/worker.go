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
	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

//worker is in charge of upload/downloading a single chunk to/from a target storage
type worker struct {
	state  int
	chunk  *chunk
	target storage.Storage
}

func (w *worker) init(ch *chunk, t storage.Storage) {
	w.state = STATE_IDLE
	w.chunk = ch
	w.target = t
}

func (w *worker) upload() {
	log.Printf("uploading chunk %s ...", w.chunk.Name)
	err := w.target.Put(w.chunk.Name, w.chunk.data)
	if err == nil {
		w.state = STATE_COMPLETED
	} else {
		w.state = STATE_FAILED
		log.Fatal(err)
	}
}

func (w *worker) download() {
	log.Printf("downloading chunk %s ...", w.chunk.Name)
	data, err := w.target.Get(w.chunk.Name)
	if err == nil {
		w.chunk.data = data
		w.state = STATE_COMPLETED
	} else {
		w.state = STATE_FAILED
		log.Fatal(err)
	}
}