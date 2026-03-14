package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/okieoth/boxes/cmd/sub"
)

var rootCmd = &cobra.Command{
	Use:   "draw",
	Short: "Tool to draw SVGs from YAML input",
	Long:  `Creates SVGs with different chart types from YAML definitions. See sub-commands for provided diagrams.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please call this with one of the provided sub-commands")
	},
}

func init() {
	rootCmd.AddCommand(sub.BoxesCmd)
	rootCmd.AddCommand(sub.RandomizeCmd)
	rootCmd.AddCommand(sub.TruncInputCmd)
	rootCmd.AddCommand(sub.VersionCmd)
}

func main() {
	rootCmd.Execute()
}
