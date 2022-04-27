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

	log "github.com/sirupsen/logrus"
	"github.com/LosAngeles971/kirinuki/business/enigma"
	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/mosaic"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/LosAngeles971/kirinuki/business/toc"
	"github.com/olekukonko/tablewriter"
)

type Gateway struct {
	email           string
	password        string
	chunksForTOC    int
	chunk_name_size int
	toc             *toc.TableOfContent
	storage         *storage.StorageMap
	tempDir         string
}

type GatewayOption func(*Gateway)

func WithStorage(m *storage.StorageMap) GatewayOption {
	return func(s *Gateway) {
		s.storage = m
	}
}

func WithTemp(temp string) GatewayOption {
	return func(s *Gateway) {
		s.tempDir = temp
	}
}

func New(email string, password string, opts ...GatewayOption) (*Gateway, error) {
	g := &Gateway{
		email:           email,
		chunksForTOC:    3,
		password:        password,
		chunk_name_size: 48,
		tempDir:         os.TempDir(),
		toc:             nil,
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.storage == nil {
		var err error
		g.storage, err = storage.NewStorageMap()
		if err != nil {
			return g, err
		}
	}
	return g, nil
}

func (g *Gateway) loadTableOfContent() (*toc.TableOfContent, error) {
	chunks := toc.GetChunks(g.email, g.password, g.storage.Array(), g.tempDir)
	if len(chunks) < 1 {
		return nil, fmt.Errorf("no chunks %v for toc", len(chunks))
	}
	ecTocFile := g.tempDir + "/" + mosaic.GetFilename(24)
	mm := mosaic.New(mosaic.WithStorage(g.storage.Array()))
	err := mm.Download(chunks, ecTocFile)
	if err != nil {
		return nil, err
	}
	ee := enigma.New(enigma.WithMainkey(g.email, g.password))
	tocFile := g.tempDir + "/" +mosaic.GetFilename(24)
	err = ee.DecryptFile(ecTocFile, tocFile)
	if err != nil {
		return nil, err
	}
	toc, err := toc.New(toc.WithFilename(tocFile))
	if err != nil {
		return nil, err
	}
	return toc, nil
}

func (g *Gateway) CreateTableOfContent() error {
	var err error
	g.toc, err = toc.New()
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
	toc, err := g.loadTableOfContent()
	if err != nil {
		return err
	}
	g.toc = toc
	return nil
}

func (g *Gateway) Logout() error {
	if !g.isOpen() {
		log.Errorf("session %s is already closed", g.email)
		return nil
	}
	tocFile := g.tempDir + "/" +mosaic.GetFilename(24)
	err := g.toc.Save(tocFile)
	if err != nil {
		return err
	}
	ecTocFile := g.tempDir + "/" +mosaic.GetFilename(24)
	ee := enigma.New(enigma.WithMainkey(g.email, g.password))
	err = ee.EncryptFile(tocFile, ecTocFile)
	if err != nil {
		return err
	}
	mm := mosaic.New(mosaic.WithStorage(g.storage.Array()), mosaic.WithTempDir(g.tempDir))
	chunks := toc.GetChunks(g.email, g.password, g.storage.Array(), g.tempDir)
	if len(chunks) < 1 {
		return fmt.Errorf("no chunks %v for toc", len(chunks))
	}
	err = mm.UploadWithChunks(chunks, ecTocFile)
	if err != nil {
		return err
	}
	g.toc = nil
	return nil
}

func (g *Gateway) Get(name string) (*kirinuki.Kirinuki, error) {
	if !g.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.email)
	}
	k, ok := g.toc.Get(name)
	if !ok {
		return nil, fmt.Errorf("file %s is not present", name)
	}
	return k, nil
}

func (g *Gateway) Find(pattern string) ([]kirinuki.Kirinuki, error) {
	if !g.isOpen() {
		return nil, fmt.Errorf("session %s is not open", g.email)
	}
	return g.toc.Find(pattern), nil
}

func (g *Gateway) Exist(name string) (bool, error) {
	if !g.isOpen() {
		return false, fmt.Errorf("session %s is not open", g.email)
	}
	return g.toc.Exist(name), nil
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
	k := kirinuki.NewKirinuki(name, kirinuki.WithRandomkey())
	err := k.Upload(filename, g.storage.Array())
	if err != nil {
		return err
	}
	ok := g.toc.Add(k)
	if !ok {
		return fmt.Errorf("failed to add %s to TOC", name)
	}
	return nil
}

func (g *Gateway) Download(name string, filename string) error {
	if !g.isOpen() {
		return fmt.Errorf("session %s is not open", g.email)
	}
	k, ok := g.toc.Get(name)
	if !ok {
		return fmt.Errorf("file %s not present", name)
	}
	return k.Download(filename, g.storage.Array())
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
	for _, t := range c.Targets {
		t1.Append([]string{"Target", t.Name()})
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