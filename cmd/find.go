package cmd

import (
	//"it/losangeles971/kirinuki/kirinuki"
	"log"
	"github.com/spf13/cobra"
)

var contain string

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find files into TOC.",
	Long: `Find files into TOC.
Usage:
	kirinuki find --contain`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Finding files into TOC...")
	},
}

func init() {
	findCmd.Flags().StringVar(&contain, "contain", "", "Substring")
	rootCmd.AddCommand(findCmd)
}


