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
package kirinuki

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
func SplitData(data []byte, n_chunks int) ([][]byte, error) {
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