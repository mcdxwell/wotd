package cmd

import (
	"fmt"

	"github.com/mcdxwell/wotd/pkg/wotd"
	"github.com/spf13/cobra"
)

var giveLink bool
var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{"td", "now"},
	Short:   "today's wotd",
	Long:    "Fetch today's wotd from Merriam-Webster",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		word := wotd.Wotd()
		fmt.Println(word)
		if giveLink {
			fmt.Println(wotd.Link(word))
		}
	},
}

func init() {
	todayCmd.Flags().BoolVarP(&giveLink, "link", "l", false, "--link")
	rootCmd.AddCommand(todayCmd)
}
