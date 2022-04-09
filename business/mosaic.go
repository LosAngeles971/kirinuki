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
	"errors"
	"fmt"
	"log"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

const (
	STATE_COMPLETED = 1
	STATE_FAILED    = 2
	STATE_IDLE      = 0

	OPERATION_UPLOAD   = 0
	OPERATION_DOWNLOAD = 1
)

// mosaic is in charge of providing upload/download functionalities for Kirinuki files
type mosaic struct {
	k           *Kirinuki
	workers      map[string][]*worker
	ss           []storage.Storage
}

func newMosaic(k *Kirinuki, ss []storage.Storage) (*mosaic, error) {
	if len(ss) == 0 {
		return nil, errors.New("empty storage array")
	}
	m := &mosaic{
		k: k,
		ss: ss,
	}
	m.workers = map[string][]*worker{}
	for _, c := range m.k.Chunks {
		m.workers[c.Name] = []*worker{}
		for _, t := range m.ss {
			w := worker{}
			w.init(c, t)
			m.workers[c.Name] = append(m.workers[c.Name], &w)
		}
	}
	return m, nil
}

// Check if the chunk has been processed (download or upload)
func (m *mosaic) availability(ch *chunk) int {
	c := 0
	for _, w := range m.workers[ch.Name] {
		if w.state == STATE_COMPLETED {
			c++
		}
	}
	return c
}

// Get one idle worker for the chunk
func (m *mosaic) getIdleWorker(ch *chunk) (*worker, bool) {
	for _, w := range m.workers[ch.Name] {
		if w.state == STATE_IDLE {
			return w, true
		}
	}
	return &worker{}, false
}

func (m *mosaic) getQueue() (map[string]*worker, error) {
	queue := map[string]*worker{}
	for _, c := range m.k.Chunks {
		if m.availability(c) < 1 {
			w, ok := m.getIdleWorker(c)
			if ok {
				queue[c.Name] = w
			} else {
				return queue, errors.New("no more idle workers for chunk: " + c.Name)
			}
		}
	}
	return queue, nil
}

func (m *mosaic) check() bool {
	for _, c := range m.k.Chunks {
		if m.availability(c) < 1 {
			return false
		}
	}
	return true
}

func (m *mosaic) run(operation int) error {
	log.Println("uploading kirinuki file...")
	max_cycles := len(m.k.Chunks)*len(m.ss) + 1
	cycle := 0
	running := true
	for running {
		queue, err := m.getQueue()
		if err != nil {
			return err
		}
		log.Printf("queue of %v workers for cycle %v", len(queue), cycle)
		if len(queue) == 0 {
			if !m.check() {
				return errors.New("operation failed")
			} else {
				return nil
			}
		}
		for _, w := range queue {
			if operation == OPERATION_DOWNLOAD {
				w.download()
			} else {
				w.upload()
			}
		}
		cycle++
		if cycle == max_cycles {
			return errors.New("reached the number of max cycles for the operation: " + fmt.Sprint(cycle))
		}
	}
	if !m.check() {
		return errors.New("operation failed")
	} else {
		return nil
	}
}

func getKirinuki(k *Kirinuki, ss []storage.Storage) ([]byte, error) {
	m, err := newMosaic(k, ss)
	if err != nil {
		return nil, err
	}
	err = m.run(OPERATION_DOWNLOAD)
	if err != nil {
		return nil, err
	}
	return m.k.Build()
}

func putKiriuki(k *Kirinuki, ss []storage.Storage) error {
	m, err := newMosaic(k, ss)
	if err != nil {
		return err
	}
	m.k.setCrushMap(m.ss)
	return m.run(OPERATION_UPLOAD)
}
