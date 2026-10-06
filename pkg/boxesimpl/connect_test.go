package boxesimpl_test

import (
	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/stretchr/testify/require"
)

func TestConnections(t *testing.T) {
	tests := []struct {
		inputFile   string
		checkFunc   func(d *boxes.BoxesDocument)
		depth       int
		expanded    []string
		blacklisted []string
	}{
		{
			inputFile: "../../ui/data/boxes_random.yaml",
			checkFunc: func(d *boxes.BoxesDocument) {
				require.NotNil(t, d)
			},
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test, test:", i)
		filtered := boxesimpl.FilterBoxes(*b, test.depth, test.expanded, test.blacklisted)
		textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
		doc, err := boxesimpl.InitialLayoutBoxes(&filtered, textDimensionCalulator)
		require.Nil(t, err, "error while initial layout, test:", i)
		doc.ConnectBoxes()
		test.checkFunc(doc)
	}
}

// TestConnectionsRenderSvgCorrectly validates that connected documents produce
// valid SVG with connection lines (from plan §3.1).
func TestConnectionsRenderSvgCorrectly(t *testing.T) {
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../ui/data/boxes_random.yaml")
	require.Nil(t, err)

	filtered := boxesimpl.FilterBoxes(*b, 2, []string{}, []string{})
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(&filtered, textDimensionCalulator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// Generate SVG and validate connection lines are present
	svgReturn := boxesimpl.DrawBoxesFiltered(*b, 2, []string{}, []string{}, false)

	// Verify connection lines exist in the SVG
	conLines := svgdrawing.ExtractConnections(svgReturn.SVG)
	require.Greater(t, len(conLines), 0, "connected doc should produce connection lines")

	// Also verify the document structure (original check)
	require.NotNil(t, doc)
}
