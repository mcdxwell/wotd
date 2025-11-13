package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wotd",
	Short: "wotd is a tool for Merriam-Webster's word of the day",
	Long:  "wotd is a tool for Merriam-Webster's word of the day",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing wotd '%s'\n", err)
		os.Exit(1)
	}
}
