package sub

import (
	"fmt"
	"os"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/spf13/cobra"
)

var mixinFiles []string
var boxesDebug bool

var BoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Draws boxes based on given layout",
	Long:  `Draws boxes and their connections and layouts them.`,
	Run: func(cmd *cobra.Command, args []string) {
		layout, err := boxesimpl.LoadBoxesFromFile(From)
		if err != nil {
			fmt.Println("Error while loading input:", err)
			os.Exit(1)
		}

		mixins := make([]boxes.BoxesFileMixings, 0, len(mixinFiles))
		for _, f := range mixinFiles {
			m, err := types.LoadInputFromFile[boxes.BoxesFileMixings](f)
			if err != nil {
				fmt.Printf("Error while loading mixin file %q: %v\n", f, err)
				os.Exit(1)
			}
			mixins = append(mixins, *m)
		}

		ret := boxesimpl.DrawBoxesFilteredExt(*layout, mixins, nil, depth, expandedIds, blacklistedIds, boxesDebug)
		if ret.ErrorMsg != "" {
			fmt.Println("Error while drawing boxes:", ret.ErrorMsg)
			os.Exit(1)
		}

		if Output == "" {
			fmt.Println(ret.SVG)
		} else {
			if err := os.WriteFile(Output, []byte(ret.SVG), 0644); err != nil {
				fmt.Println("Error while writing output:", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	initDefaultFlags(BoxesCmd)
	BoxesCmd.Flags().StringSliceVarP(&mixinFiles, "mixin", "m", []string{}, "Paths to mixin files, can be used multiple times")
	BoxesCmd.Flags().StringSliceVarP(&expandedIds, "expand", "e", []string{}, "IDs of boxes to expand, can be used multiple times")
	BoxesCmd.Flags().StringSliceVarP(&blacklistedIds, "blacklisted", "b", []string{}, "IDs of blacklisted boxes, can be used multiple times")
	BoxesCmd.Flags().IntVarP(&depth, "depth", "d", 2, "Default depth to show, default is 2")
	BoxesCmd.Flags().BoolVar(&boxesDebug, "debug", false, "Enable debug rendering (roads, start positions, connection nodes)")
}
