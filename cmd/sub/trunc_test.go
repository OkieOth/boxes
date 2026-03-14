package sub_test

import (
	"testing"

	"github.com/okieoth/boxes/cmd/sub"
	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/stretchr/testify/require"
)

func TestTruncBoxes(t *testing.T) {
	tests := []struct {
		inputFile      string
		outputFile     string
		depth          int
		expandedIds    []string
		blacklistedIds []string
		checkFunc      func(t *testing.T, outputFile string, testIndex int)
	}{
		// {
		// 	inputFile:      "../../ui/data/boxes_random.yaml",
		// 	outputFile:     "../../temp/boxes_random_truncated.yaml",
		// 	depth:          2,
		// 	expandedIds:    []string{"id_1_1"},
		// 	blacklistedIds: []string{},
		// },
		{
			inputFile:      "../../ui/data/boxes_random.yaml",
			outputFile:     "../../temp/boxes_random_truncated3.yaml",
			depth:          2,
			expandedIds:    []string{"id_3_6"},
			blacklistedIds: []string{},
			checkFunc: func(t *testing.T, outputFile string, testIndex int) {
				boxes, err := boxesimpl.LoadBoxesFromFile(outputFile)
				require.Nil(t, err, "error while loading created file, test:")
				require.Len(t, boxes.Boxes.Vertical[2].Horizontal[5].Vertical[0].Connections, 1)
			},
		},
	}
	for i, test := range tests {
		err := sub.TruncBoxes(test.inputFile, test.outputFile, test.depth, test.expandedIds, test.blacklistedIds)
		require.Nil(t, err, "error in TruncBoxes call, test:", i)
		require.FileExists(t, test.outputFile)
		test.checkFunc(t, test.outputFile, i)
	}
}
