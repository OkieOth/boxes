package sub_test

import (
	"testing"

	"github.com/okieoth/boxes/cmd/sub"
	"github.com/stretchr/testify/require"
)

func TestRandomizeBoxes(t *testing.T) {
	input := "../../resources/examples_boxes/complex_complex.yaml"
	output := "../../temp/complex_complex_randomized.yaml"
	err := sub.RandomizeBoxes(input, output)
	require.Nil(t, err)
	require.FileExists(t, output)
}
