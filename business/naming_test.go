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
	"strings"
	"testing"
)

func TestNames(t *testing.T) {
	n := newNaming()
	name := n.getNameForChunk()
	if len(name) != n.chunk_name_size {
		t.Fatalf("expected length %v not %v", n.chunk_name_size, len(name))
	}
	if strings.Contains(name, "/") {
		t.Fatal("name for chunk cannot contain the slash")
	}
	s, err := NewSession(test_email, test_password)
	if err != nil {
		t.Fatalf("session failed %v", err)
	}
	for i := 0; i < s.chunksForTOC; i++ {
		name := n.getNameForTOCChunk(s, i)
		if strings.Contains(name, "/") {
			t.Fatal("name for chunk cannot contain the slash")
		}
	}
}