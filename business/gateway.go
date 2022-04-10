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
	"fmt"

	"github.com/LosAngeles971/kirinuki/business/storage"
)

type Gateway struct {
	session *Session
}

func New(email string, password string, scratch bool, m *storage.StorageMap) (Gateway, error) {
	s, err := NewSession(email, password, scratch, WithStorage(m))
	if err != nil {
		return Gateway{}, err
	}
	return Gateway{
		session: s,
	}, nil
}

func (g Gateway) Login() error {
	return g.session.login()
}

func (g Gateway) Logout() error {
	return g.session.logout()
}

func (g Gateway) Find(pattern string) ([]Kirinuki, error) {
	if !g.session.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return nil, err
	}
	return toc.find(pattern), nil
}

func (g Gateway) Upload(name string, data []byte, overwrite bool) error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return err
	}
	if toc.exist(name) && !overwrite {
		return fmt.Errorf("file with name %s already exists", name)
	}
	// overwrite in any case
	k := NewKirinuki(name)
	err = k.addData(data)
	if err != nil {
		return err
	}
	err = putKiriuki(k, g.session.storage.Array())
	if err != nil {
		return err
	}
	ok := toc.add(k)
	if !ok {
		return fmt.Errorf("failed to add %s to TOC", name)
	}
	return nil
}

func (g Gateway) Download(name string) ([]byte, error) {
	if !g.session.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return nil, err
	}
	k, ok := toc.get(name)
	if !ok {
		return nil, fmt.Errorf("file %s not present", name)
	}
	return getKirinuki(k, g.session.storage.Array())
}