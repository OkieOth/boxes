package sub

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/spf13/cobra"
)

var RandomizeCmd = &cobra.Command{
	Use:   "randomize",
	Short: "Randomize the texts in input files",
	Long:  `Randomize the texts in input files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("call this command with one of the available sub commands")
	},
}

var RandomizeBoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Randomize the texts in boxes input files",
	Long:  `Randomize the texts in boxes input files`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RandomizeBoxes(From, Output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	initDefaultFlags(RandomizeCmd)
	RandomizeCmd.AddCommand(RandomizeBoxesCmd)
}

func RandomizeBoxes(input, output string) error {
	boxes, err := boxesimpl.LoadBoxesFromFile(input)
	if err != nil {
		return fmt.Errorf("Error while loading input: %v", err)
	}
	faker := gofakeit.NewFaker(source.NewJSF(11), true)
	randomizeBoxesLayout(&boxes.Boxes, faker)
	return serializeToYaml(boxes, output)
}

func randomizeText(txt string, faker *gofakeit.Faker) string {
	if txt == "" {
		return txt
	}
	return faker.Sentence()
}

func randomizeBoxesLayoutContainer(cont []boxes.Layout, faker *gofakeit.Faker) []boxes.Layout {
	if len(cont) == 0 {
		return []boxes.Layout{}
	}
	ret := make([]boxes.Layout, 0)
	for _, c := range cont {
		randomizeBoxesLayout(&c, faker)
		ret = append(ret, c)
	}
	return ret
}

func randomizeBoxesLayout(l *boxes.Layout, faker *gofakeit.Faker) {
	l.Caption = randomizeText(l.Caption, faker)
	l.Text1 = randomizeText(l.Text1, faker)
	l.Text2 = randomizeText(l.Text2, faker)
	l.Horizontal = randomizeBoxesLayoutContainer(l.Horizontal, faker)
	l.Vertical = randomizeBoxesLayoutContainer(l.Vertical, faker)
}
