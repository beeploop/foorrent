package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "foorrent",
	Short: "A simple CLI BitTorrent client",
	Long:  `A simple CLI BitTorrent client`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//
}
