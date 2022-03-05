/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Mosaic is in charge of uploading/downloading functions of Kirinuki's files

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
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
