/*
 * Created on Sat Apr 09 2022
 * Author @LosAngeles971
 *
 * The MIT License (MIT)
 * Copyright (c) 2022 @LosAngeles971
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software
 * and associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial
 * portions of the Software.
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