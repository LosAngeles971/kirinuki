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
	"os"
	"time"

	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/olekukonko/tablewriter"
)

type Gateway struct {
	session *Session
}

func New(email string, password string, m *storage.StorageMap) (Gateway, error) {
	s, err := NewSession(email, password, WithStorage(m))
	if err != nil {
		return Gateway{}, err
	}
	return Gateway{
		session: s,
	}, nil
}

func (g Gateway) CreateTableOfContent() error {
	return g.session.createTableOfContent()
}

func (g Gateway) Login() error {
	return g.session.login()
}

func (g Gateway) Logout() error {
	return g.session.logout()
}

func (g Gateway) Get(name string) (*Kirinuki, error) {
	if !g.session.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return nil, err
	}
	k, ok := toc.get(name)
	if !ok {
		return nil, fmt.Errorf("file %s is not present", name)
	}
	return k, nil
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

func (g Gateway) Exist(name string) (bool, error) {
	if !g.session.isOpen() {
		return false, fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return false, err
	}
	return toc.exist(name), nil
}

func (g Gateway) Size() (int, error) {
	if !g.session.isOpen() {
		return 0, fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return 0, err
	}
	return len(toc.Kfiles), nil
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
	k := NewKirinuki(name, g.session.getChunks(data), WithRandomkey())
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

func (g Gateway) Info() error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Last update", time.Unix(toc.Lastupdate, 0).String()})
	table.Append([]string{"Size", fmt.Sprint(len(toc.Kfiles))})
	table.Render()
	return nil
}

func (g Gateway) PrintChunk(c *chunk) {
	t1 := tablewriter.NewWriter(os.Stdout)
	t1.Append([]string{"Index", fmt.Sprint(c.Index)})
	t1.Append([]string{"Name", c.Name})
	t1.Append([]string{"Size", fmt.Sprint(c.Real_size)})
	for _, t := range c.Targets {
		t1.Append([]string{"Target", t})
	}
	t1.Render()
}

func (g Gateway) Stat(name string) error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	toc, err := g.session.getTOC()
	if err != nil {
		return err
	}
	k, ok := toc.get(name)
	if !ok {
		return fmt.Errorf("file %s not present", name)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Name", name})
	table.Append([]string{"Date", time.Unix(k.Date, 0).String()})
	table.Append([]string{"Encryption", fmt.Sprint(k.Encryption)})
	table.Append([]string{"Replicas", fmt.Sprint(k.Replicas)})
	table.Append([]string{"Checksum", fmt.Sprint(k.Checksum)})
	table.Append([]string{"Chunks", fmt.Sprint(len(k.Chunks))})
	table.Render()
	for _, c := range k.Chunks {
		g.PrintChunk(c)
	}
	return nil
}