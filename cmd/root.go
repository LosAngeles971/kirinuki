package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var name string
var filename string
var email string
var scratch bool
var overwrite bool
var debug bool
var storageMap string

var rootCmd = &cobra.Command{
	Use:   "kirinuki",
	Short: "Kirinuki",
	Long:  `Kirinuki.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&email, "email", "", "email")
	rootCmd.PersistentFlags().BoolVar(&scratch, "scratch", false, "scratch your table of content")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug")
	rootCmd.PersistentFlags().StringVar(&storageMap, "storage", "", "storage file config")
	rootCmd.MarkPersistentFlagRequired("email")
}
