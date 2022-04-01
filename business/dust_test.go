/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

testing file splitting features

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
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