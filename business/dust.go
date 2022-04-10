/*
 * Created on Fri Apr 08 2022
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
	"errors"
)

// getChunksNumberForKFile returns the number of chunks depending on the size of the source file
func getChunksNumberForKFile(file []byte) int {
	size := len(file)
	if size < 1000 {
		return 3
	}
	if size < 10000 {
		return 5
	}
	if size < 100000 {
		return 7
	}
	if size < 1000000 {
		return 9
	}
	return 11
}

// extractChunk extracts a subset from an array of bytes starting from the given index
func extractChunk(data []byte, startindex int, size int) []byte {
	chunk := []byte{}
	if size > 0 {
		for n := 0; n < size; n++ {
			index := startindex + n
			if index < len(data) {
				chunk = append(chunk, data[index])
			}
		}
		return chunk
	}
	index := startindex
	for index < len(data) {
		chunk = append(chunk, data[index])
		index++
	}
	return chunk
}

// splitFile splits a single array of bytes into an array of chunks
func splitFile(data []byte, n_chunks int) ([][]byte, error) {
	chunks := [][]byte{}
	if n_chunks < 1 {
		return chunks, errors.New("0 chunks not possible")
	}
	if n_chunks > (len(data) / 2) {
		return chunks, errors.New("chunks cannot be greater than half file's size")
	}
	chunk_size := len(data) / n_chunks
	start_index := 0
	for n := 0; n < n_chunks; n++ {
		if n == n_chunks-1 {
			chunks = append(chunks, extractChunk(data, start_index, -1))
		} else {
			chunks = append(chunks, extractChunk(data, start_index, chunk_size))
		}
		start_index += chunk_size
	}
	return chunks, nil
}
