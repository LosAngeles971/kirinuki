package cmd

import (
	"io/ioutil"

	"github.com/LosAngeles971/kirinuki/business"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download Kirinuki file",
	Long: `download Kirinuki file.
Usage:
	kirinuki download --email <email> --name <name> --filename <filename>`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := business.New(email, askPassword(), scratch, getStorageMap())
		if err != nil {
			log.Fatalf("failed to create Gateway due to %v", err)
		}
		log.Infof("downloading %s ...", name)
		data, err := g.Download(name)
		if err != nil {
			log.Fatalf("failed to download %s due to %v", name, err)
		}
		log.Infof("saving %s to local file %s ...", name, filename)
		err = ioutil.WriteFile(filename, data, 0755)
		if err != nil {
			log.Fatalf("failed to save %s to %s due to %v", name, filename, err)
		}
		log.Infof("saved %s to local file %s", name, filename)
	},
}

func init() {
	downloadCmd.PersistentFlags().StringVar(&name, "name", "", "name of file")
	downloadCmd.PersistentFlags().StringVar(&filename, "filename", "", "output filename")
	rootCmd.AddCommand(downloadCmd)
}
