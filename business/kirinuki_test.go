package business

import (
	"testing"
)

type k_data_test struct {
	name string
	data []byte
}
var k_data_tests []k_data_test = []k_data_test{
	{
		name: "test1",
		data: []byte{0,1,2,3,4,5,6,7,8,9}, 
	},
	{
		name: "test2",
		data: []byte("0,1,2,3,4,5,6,7,8,9,0,a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z"),
	},
}

func TestKirinuki(t *testing.T) {
	for _, tt := range k_data_tests {
		k, err := NewKirinuki(WithKirinukiData(tt.name, tt.data))
		if err != nil {
			t.Fatal(err)
		}
		if len(k.Chunks) != getChunksNumberForKFile(tt.data) {
			t.Fatalf("unxpected number of chunk %v  expected %v", len(k.Chunks), getChunksNumberForKFile(tt.data))
		}
	}
}