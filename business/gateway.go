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

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/business/toc"
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
	var err error
	g.session.toc, err = toc.New()
	if err != nil {
		g.session.toc = nil
		return err
	}
	return nil
}

func (g Gateway) Login() error {
	return g.session.login()
}

func (g Gateway) Logout() error {
	return g.session.logout()
}

func (g Gateway) Get(name string) (*kirinuki.Kirinuki, error) {
	if !g.session.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.session.email)
	}
	k, ok := g.session.toc.Get(name)
	if !ok {
		return nil, fmt.Errorf("file %s is not present", name)
	}
	return k, nil
}

func (g Gateway) Find(pattern string) ([]kirinuki.Kirinuki, error) {
	if !g.session.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.session.email)
	}
	return g.session.toc.Find(pattern), nil
}

func (g Gateway) Exist(name string) (bool, error) {
	if !g.session.isOpen() {
		return false, fmt.Errorf("session %s is not open", g.session.email)
	}
	return g.session.toc.Exist(name), nil
}

func (g Gateway) Size() (int, error) {
	if !g.session.isOpen() {
		return 0, fmt.Errorf("session %s is not open", g.session.email)
	}
	return g.session.toc.Size(), nil
}

func (g Gateway) Upload(filename string, name string, overwrite bool) error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	if g.session.toc.Exist(name) && !overwrite {
		return fmt.Errorf("file %s already exists", name)
	}
	k := kirinuki.NewKirinuki(name, kirinuki.WithRandomkey())
	err := k.Upload(filename, g.session.storage.Array())
	if err != nil {
		return err
	}
	ok := g.session.toc.Add(k)
	if !ok {
		return fmt.Errorf("failed to add %s to TOC", name)
	}
	return nil
}

func (g Gateway) Download(name string, filename string) error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	k, ok := g.session.toc.Get(name)
	if !ok {
		return fmt.Errorf("file %s not present", name)
	}
	return k.Download(filename, g.session.storage.Array())
}

func (g Gateway) Info() error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Last update", time.Unix(g.session.toc.Lastupdate, 0).String()})
	table.Append([]string{"Size", fmt.Sprint(g.session.toc.Size())})
	table.Render()
	return nil
}

func (g Gateway) PrintChunk(c *mosaic.Chunk) {
	t1 := tablewriter.NewWriter(os.Stdout)
	t1.Append([]string{"Index", fmt.Sprint(c.Index)})
	t1.Append([]string{"Name", c.Name})
	t1.Append([]string{"Size", fmt.Sprint(c.Real_size)})
	t1.Append([]string{"Checksum", c.Checksum})
	for _, t := range c.Targets {
		t1.Append([]string{"Target", t.Name()})
	}
	t1.Render()
}

func (g Gateway) Stat(name string) error {
	if !g.session.isOpen() {
		return fmt.Errorf("session %s is not open", g.session.email)
	}
	k, ok := g.session.toc.Get(name)
	if !ok {
		return fmt.Errorf("file %s not present", name)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Name", name})
	table.Append([]string{"Date", time.Unix(k.Date, 0).String()})
	table.Append([]string{"Checksum", fmt.Sprint(k.Checksum)})
	table.Append([]string{"Chunks", fmt.Sprint(len(k.Chunks))})
	table.Render()
	for _, c := range k.Chunks {
		g.PrintChunk(c)
	}
	return nil
}