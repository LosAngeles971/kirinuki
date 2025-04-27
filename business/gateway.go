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
	"github.com/LosAngeles971/kirinuki/business/multistorage"
	"github.com/LosAngeles971/kirinuki/business/toc"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type Gateway struct {
	email           string
	password        string
	toc             *toc.TableOfContent
	storage         *multistorage.MultiStorage
}

type GatewayOption func(*Gateway)

func WithStorage(m *multistorage.MultiStorage) GatewayOption {
	return func(s *Gateway) {
		s.storage = m
	}
}

func New(email string, password string, opts ...GatewayOption) (*Gateway, error) {
	g := &Gateway{
		email:           email,
		password:        password,
		toc:             nil,
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.storage == nil {
		var err error
		g.storage, err = multistorage.New()
		if err != nil {
			return g, err
		}
		g.storage.AddLocal("tmp", os.TempDir())
	}
	return g, nil
}

func (g *Gateway) SetEmptyTableOfContent() error {
	var err error
	g.toc, err = toc.New(g.storage)
	if err != nil {
		g.toc = nil
		return err
	}
	return nil
}

func (g *Gateway) isOpen() bool {
	return g.toc != nil
}

func (g *Gateway) Login() error {
	if g.isOpen() {
		return nil
	}
	err := g.SetEmptyTableOfContent()
	if err != nil {
		return err
	}
	err = g.toc.Load(g.email, g.password)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) Logout() error {
	if !g.isOpen() {
		log.Errorf("session %s is already closed", g.email)
		return nil
	}
	err := g.toc.Store(g.email, g.password)
	if err != nil {
		return err
	}
	g.toc = nil
	return nil
}

func (g *Gateway) Exist(name string) (*kirinuki.File, error) {
	if !g.isOpen() {
		return nil, fmt.Errorf("session %s not open", g.email)
	}
	if f, ok := g.toc.Get(name); ok {
		return f, nil
	} else {
		return nil, fmt.Errorf("file %s not present", name)
	}
}

func (g *Gateway) Find(pattern string) ([]*kirinuki.File, error) {
	if !g.isOpen() {
		return nil, fmt.Errorf("session %s not open", g.email)
	}
	return g.toc.Find(pattern), nil
}

func (g *Gateway) Size() (int, error) {
	if !g.isOpen() {
		return 0, fmt.Errorf("session %s is not open", g.email)
	}
	return g.toc.Size(), nil
}

func (g *Gateway) Upload(filename string, name string, overwrite bool) error {
	if !g.isOpen() {
		return fmt.Errorf("session %s is not open", g.email)
	}
	if g.toc.Exist(name) && !overwrite {
		return fmt.Errorf("file %s already exists", name)
	}
	f := mosaic.NewFile(name, kirinuki.WithRandomkey())
	err := f.Upload(filename, g.storage)
	if err != nil {
		return err
	}
	if !g.toc.Add(f) {
		return fmt.Errorf("failed to add %s to TOC", name)
	}
	return nil
}

func (g *Gateway) Download(name string, filename string) error {
	if !g.isOpen() {
		return fmt.Errorf("session %s is not open", g.email)
	}
	f, ok := g.toc.Get(name)
	if !ok {
		return fmt.Errorf("file %s not present", name)
	}
	return f.Download(filename, g.storage)
}

func (g *Gateway) Info() error {
	if !g.isOpen() {
		return fmt.Errorf("session %s is not open", g.email)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append([]string{"Last update", time.Unix(g.toc.Lastupdate, 0).String()})
	table.Append([]string{"Size", fmt.Sprint(g.toc.Size())})
	table.Render()
	return nil
}

func (g *Gateway) PrintChunk(c *mosaic.Chunk) {
	t1 := tablewriter.NewWriter(os.Stdout)
	t1.Append([]string{"Index", fmt.Sprint(c.Index)})
	t1.Append([]string{"Name", c.Name})
	t1.Append([]string{"Size", fmt.Sprint(c.Real_size)})
	t1.Append([]string{"Checksum", c.Checksum})
	for _, t := range c.TargetNames {
		t1.Append([]string{"Target", t})
	}
	t1.Render()
}

func (g *Gateway) Stat(name string) error {
	if !g.isOpen() {
		return fmt.Errorf("session %s is not open", g.email)
	}
	k, ok := g.toc.Get(name)
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
