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
package storage

import (
	"os"
	"testing"
)

// TestLocal verifies put/get files using local storage target
func TestLocal(t *testing.T) {
	ll, err := newStorage("test", ConfigItem{
		Type: "local",
		Cfg: map[string]string{
			"path": os.TempDir(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	input := "test of filesystem storage"
	err = ll.Put("test1", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	output, err := ll.Get("test1")
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != input {
		t.Fatalf("expected [%s] not [%s]", input, string(output))
	}
}