package cmd

import (
	"io/ioutil"

	"github.com/LosAngeles971/kirinuki/business"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload Kirinuki file",
	Long: `upload Kirinuki file.
Usage:
	kirinuki upload --email <email> --name <name> --filename <filename>`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("failed to load %s due to %v", filename, err)
		}
		g, err := business.New(email, askPassword(), scratch, getStorageMap())
		if err != nil {
			log.Fatalf("failed to create Gateway due to %v", err)
		}
		log.Infof("uploading %s ...", name)
		err = g.Upload(name, data, overwrite)
		if err != nil {
			log.Fatalf("failed to upload %s due to %v", name, err)
		}
		log.Infof("uploaded %s with %s", name, filename)
	},
}

func init() {
	uploadCmd.PersistentFlags().StringVar(&name, "name", "", "name of file")
	uploadCmd.PersistentFlags().StringVar(&filename, "filename", "", "output filename")
	uploadCmd.PersistentFlags().BoolVar(&overwrite, "overwrite", false, "overwrite file if exists")
	rootCmd.AddCommand(uploadCmd)
}
