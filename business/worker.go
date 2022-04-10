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