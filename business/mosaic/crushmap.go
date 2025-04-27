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
 package mosaic

 import (
	 "github.com/LosAngeles971/kirinuki/business/config"
	 "github.com/LosAngeles971/kirinuki/business/helpers"
)

// assigning storage targets to every chunk of the file
func getCrushMap(targets []string) []*Chunk {
	chunks := []*Chunk{}
	// create a chunk for every available storage
	for i := 0; i < len(targets); i++ {
		name := helpers.GetFilename(nameSize)
		c := NewChunk(i, name, WithFilename(config.GetTmp()+"/"+name))
		c.TargetNames = targets
		chunks = append(chunks, c)
	}
	return chunks
}