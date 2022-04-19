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
package storage

import (
	_ "embed"
	"testing"
)

//go:embed sftp.json
var sftpFile []byte

//go:embed minio.json
var minioFile []byte

func doTest(sName string, sFile []byte, t *testing.T) {
	sm, err := NewStorageMap(WithJSONData(sFile))
	if err != nil {
		t.Fatal(err)
	}
	s, err := sm.Get(sName)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Put("testfile", sftpFile)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Get("testfile")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSFTP(t *testing.T) {
	doTest("sftp", sftpFile, t)
}

func TestMinio(t *testing.T) {
	doTest("minio", minioFile, t)
}