/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Dust is in charge of splitting a file into multiple chunks and adding a pad to each chunk

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
package business

import (
	"errors"
)

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
