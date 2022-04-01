package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/spf13/cobra"

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

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "storage manager",
	Long: `storage manager.
Usage:
	kirinuki storage <subcommand> [options]`,
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

var storageInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init storage file",
	Long: `init storage file.
Usage:
	kirinuki storage init --storage <filename>`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := storage.Init(storageMap)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(storageMap, data, 0755)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageInitCmd)
	rootCmd.AddCommand(storageCmd)
}