package business

import (
	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

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