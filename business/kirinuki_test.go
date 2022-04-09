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

// test if the entire file is splitted into the expected number of chunks
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