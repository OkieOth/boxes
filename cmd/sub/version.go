package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.14.1"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the program",
	Long:  "Shows the version of the program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
