package toc

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

import (
	_ "embed"
	"testing"

	"github.com/LosAngeles971/kirinuki/business/kirinuki"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/stretchr/testify/require"
)

func TestTOC(t *testing.T) {
	tsm := storage.NewTestLocalMultistorage("toc")
	toc, err :=	New(tsm.GetMultiStorage())
	require.Nil(t, err)
	k := kirinuki.NewFile("test", kirinuki.WithRandomkey())
	ok := toc.Add(k)
	require.True(t, ok)
	err = toc.Store(storage.Test_email, storage.Test_password)
	require.Nil(t, err)
	toc.Files = nil
	err = toc.Load(storage.Test_email, storage.Test_password)
	require.Nil(t, err)
	require.True(t, toc.Exist("test"))
	tsm.Clean()
}