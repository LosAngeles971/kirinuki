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
	"testing"
)

type split_test struct {
	input  []byte
	output [][]byte
	chunks int
}

var split_test1 []split_test = []split_test{
	{
		input: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		output: [][]byte{
			{0, 1, 2},
			{3, 4, 5},
			{6, 7, 8, 9},
		},
		chunks: 3,
	},
}

func TestSplit(t *testing.T) {
	for _, test := range split_test1 {
		chunks, err := splitFile(test.input, test.chunks)
		if err != nil {
			t.Fatal(err)
		}
		if len(chunks) != len(test.output) {
			t.Fatalf("expected %v chunks not %v", len(test.output), len(chunks))
		}
		for i := range test.output {
			if len(test.output[i]) != len(chunks[i]) {
				t.Fatalf("expected size of %v not %v for chunk %v", len(test.output[i]), len(chunks[i]), i)
			}
			for j := range test.output[i] {
				if test.output[i][j] != chunks[i][j] {
					t.Fatalf("expected value %v not %v for chunk %v at index %v", test.output[i][j], chunks[i][j], i, j)
				}
			}
		}
	}
}