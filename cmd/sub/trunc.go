package sub

import (
	"fmt"
	"os"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/spf13/cobra"
)

var TruncInputCmd = &cobra.Command{
	Use:   "trunc",
	Short: "Truncate in input files",
	Long:  `Randomize the texts in input files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("call this command with one of the available sub commands")
	},
}

var TruncBoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Truncate boxes from the input files",
	Long:  `Truncate boxes from the input files, mostly for debug reasons`,
	Run: func(cmd *cobra.Command, args []string) {
		err := TruncBoxes(From, Output, depth, expandedIds, blacklistedIds)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func TruncBoxes(input, output string, depth int, expanded, blacklisted []string) error {
	boxes, err := boxesimpl.LoadBoxesFromFile(input)
	if err != nil {
		return fmt.Errorf("Error while loading input: %v", err)
	}
	filtered := boxesimpl.FilterBoxes(*boxes, depth, expanded, blacklisted)
	return serializeToYaml(&filtered, output)
}

var expandedIds []string
var blacklistedIds []string
var depth int

func init() {
	initDefaultFlags(TruncInputCmd)
	TruncInputCmd.PersistentFlags().StringSliceVarP(&expandedIds, "expand", "e", []string{}, "IDs of boxes to expand, can be multible times used")
	TruncInputCmd.PersistentFlags().StringSliceVarP(&blacklistedIds, "blacklisted", "b", []string{}, "IDs of blacklisted boxes, can be multible times used")
	TruncInputCmd.PersistentFlags().IntVarP(&depth, "depth", "d", 2, "Default depth to truncate, default is 2")

	TruncInputCmd.AddCommand(TruncBoxesCmd)
}
