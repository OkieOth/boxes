package sub

import (
	"fmt"
	"os"

	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var From string
var Output string

func initDefaultFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&From, "from", "f", "", "Path to the input file (required)")
	cmd.PersistentFlags().StringVarP(&Output, "output", "o", "", "Output to create, default is stdout")
	cmd.MarkPersistentFlagRequired("from")
	cmd.MarkPersistentFlagFilename("from")
	cmd.MarkPersistentFlagFilename("output")
}

func serializeToYaml(b *boxes.Boxes, output string) error {
	bytes, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Errorf("Error while serialize to yaml: %v", err)
	}
	if output == "" {
		fmt.Println(string(bytes))
	} else {
		outputFile, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("Error while creating output file: %v", err)
		}
		defer outputFile.Close()
		_, err = outputFile.Write(bytes)
		if err != nil {
			return fmt.Errorf("Error while writing output file: %v", err)
		}
	}
	return nil
}
