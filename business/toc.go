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
	"encoding/json"
	"regexp"
	"time"
)

// TOC handles the table of content of the Kirinuki files
type TOC struct {
    Lastupdate	int64				`json:"lastupdate"`
    Kfiles		[]*Kirinuki 		`json:"kfiles"`
}

type TOCOption func(*TOC) error

// TOCWithData is used to load an existent table of content
func TOCWithData(data []byte) TOCOption {
	return func(t *TOC) error {
		return json.Unmarshal(data, &t)
	}
}

func newTOC(opts ...TOCOption) (*TOC, error) {
	t := &TOC{
		Lastupdate: time.Now().UnixNano(),
		Kfiles: []*Kirinuki{},
	}
	for _, opt := range opts {
		err := opt(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Exist returns true if the Kirinuki file with the given name exists
func (t TOC) exist(name string) bool {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return true
		}
	}
	return false
}

// Get returns the Kirinuki file with the given name and true if the file exists
// This method returns the Kirinuki file without the data
func (t TOC) get(name string) (*Kirinuki, bool) {
	for _, k := range t.Kfiles {
		if name == k.Name {
			return k, true
		}
	}
	return nil, false
}

func (t *TOC) add(k *Kirinuki) bool {
	if t.exist(k.Name) {
		return false
	}
	t.Kfiles = append(t.Kfiles, k)
	return true
}

func (t TOC) find(pattern string) []Kirinuki {
	rr := []Kirinuki{}
	for _, k := range t.Kfiles {
		match, _ := regexp.MatchString(pattern, k.Name)
		if match {
			rr = append(rr, *k)
		}
	}
	return rr
}
