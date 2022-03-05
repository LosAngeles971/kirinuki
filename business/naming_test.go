/*+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Testing naming features

This testing session DOES NOT need external data, nether to interact with the system

+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"testing"
)

func TestChunkName(t *testing.T) {
	n := newNaming()
	name := n.getNameForChunk()
	if len(name) != n.chunk_name_size {
		t.FailNow()
	}
}