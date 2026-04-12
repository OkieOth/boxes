package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

<<<<<<< HEAD
const Version = "1.4.1"
=======
const Version = "1.4.0"
>>>>>>> main

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the program",
	Long:  "Shows the version of the program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
