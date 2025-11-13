package cmd

import (
	"fmt"

	"github.com/mcdxwell/wotd/pkg/wotd"
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:     "random",
	Aliases: []string{"rand"},
	Short:   "random wotd",
	Long:    "Fetch random wotd from Merriam-Webster",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		word := wotd.RandomWord()
		fmt.Println(word)
		if giveLink {
			fmt.Println(wotd.Link(word))
		}
	},
}

func init() {
	randomCmd.Flags().BoolVarP(&giveLink, "link", "l", false, "--link")
	rootCmd.AddCommand(randomCmd)
}
