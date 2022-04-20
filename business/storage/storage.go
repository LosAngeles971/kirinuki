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
	"fmt"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/sftp"
)

type Storage interface {
	Name() string
	Get(name string) ([]byte, error)
	Put(filename string, data []byte) error
}

// ConfigItem contains configurations for a specific storage target
type ConfigItem struct {
	Type string            `yaml:"type" json:"type"`
	Cfg  map[string]string `yaml:"config" json:"config"`
}

// newStorage creates a new Storage object for a given specific storage type
func newStorage(name string, ci ConfigItem) (Storage, error) {
	switch ci.Type {
	case local.Kind, s3.Kind:
		return newStowStorage(name, ci)
	case sftp.Kind:
		return NewSFTP(name, ci.Cfg), nil
	default:
		return nil, fmt.Errorf("unrecognized type of storage %s", ci.Type)
	}
}
