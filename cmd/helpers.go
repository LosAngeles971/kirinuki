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
package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/LosAngeles971/kirinuki/business"
	"github.com/LosAngeles971/kirinuki/business/storage"

	log "github.com/sirupsen/logrus"
)

func getStorageMap() *storage.StorageMap {
	log.Debugf("getting storage from config file %s", storageMap)
	data, err := ioutil.ReadFile(storageMap)
	if err != nil {
		log.Fatalf("cannot load storage map from file %s", storageMap)
	}
	var sm *storage.StorageMap
	var err2 error
	if strings.HasSuffix(storageMap, ".yml") || strings.HasSuffix(storageMap, ".yaml") {
		sm, err2 = storage.NewStorageMap(storage.WithYAMLData(data))
	} else {
		sm, err2 = storage.NewStorageMap(storage.WithJSONData(data))
	}
	if err2 != nil {
		log.Fatalf("storage map file %s is corrupted or invalid, err = %v", storageMap, err)
	}
	return sm
}

func getGateway(email string, password string) *business.Gateway {
	g, err := business.New(email, askPassword(), business.WithStorage(getStorageMap()))
	if err != nil {
		log.Fatalf("gateway creation failed [%v]", err)
		panic(err)
	}
	return g
}