package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/LosAngeles971/kirinuki/business"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find files into the storage using regex matching",
	Long: `find files into the storage regex matching.
Usage:
	kirinuki find --name <regex pattern>`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := business.New(email, askPassword(), scratch, getStorageMap())
		if err != nil {
			log.Fatalf("failed to create Gateway due to %v", err)
		}
		rr, err := g.Find(name)
		if err != nil {
			log.Fatal(err)
		}
		for _, k := range rr {
			log.Infof("kirinuki %s - date %v ", k.Name, k.Date)
		}
	},
}

func init() {
	findCmd.PersistentFlags().StringVar(&name, "name", "", "pattern for name finding")
	findCmd.MarkPersistentFlagRequired("name")
	rootCmd.AddCommand(findCmd)
}


