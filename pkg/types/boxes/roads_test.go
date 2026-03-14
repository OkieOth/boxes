package boxes_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortAscending(t *testing.T) {
	testData := []struct {
		input    []int
		expected []int
	}{
		{[]int{3, 1, 2}, []int{1, 2, 3}},
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{[]int{}, []int{}},
	}

	for _, d := range testData {
		types.SortAscending(d.input)
		assert.Equal(t, d.expected, d.input)
	}
}

func TestSortDescending(t *testing.T) {
	testData := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{5, 4, 3, 2, 1}, []int{5, 4, 3, 2, 1}},
		{[]int{}, []int{}},
	}

	for _, d := range testData {
		types.SortDescending(d.input)
		assert.Equal(t, d.expected, d.input)
	}
}

func TestRoads(t *testing.T) {
	testData := []struct {
		inputFile  string
		outputFile string
		checkFunc  func(t *testing.T, doc *boxes.BoxesDocument, i int)
	}{
		{
			inputFile:  "../../../resources/examples_boxes/complex_horizontal_connected_pics2.yaml",
			outputFile: "../../../temp/complex_horizontal_connected_pics2_roads.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument, i int) {
				require.NotNil(t, doc, "test:", i)
			},
		},
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
	for i, test := range testData {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		doc.InitStartPositions()
		doc.InitRoads()
		test.checkFunc(t, doc, i)
		output, err := os.Create(test.outputFile)
		require.Nil(t, err)
		svgdrawing := svgdrawing.NewDrawing(output)
		svgdrawing.Start(doc.Title, doc.Height, doc.Width)
		svgdrawing.InitImages(doc.Images)
		doc.DrawBoxes(svgdrawing)
		svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
		doc.DrawRoads(svgdrawing)
		svgdrawing.Done()
		output.Close()
		_, err = os.Stat(test.outputFile)
		require.Nil(t, err)

		checkForRoadOverlapsHorizontally(t, doc, i)
		checkForRoadOverlapsVertically(t, doc, i)
	}
}

func shouldHandle(elem *boxes.LayoutElement) bool {
	if elem.Caption == "" && elem.Text1 == "" && elem.Text2 == "" && elem.Id == "" {
		return false
	}
	return true
}

func checkForRoadOverlapsHorizontally(t *testing.T, doc *boxes.BoxesDocument, testIndex int) {
	for lineIndex, line := range doc.HorizontalRoads {
		checkOverlapsHorizontalForLayout(t, &doc.Boxes, line, testIndex, lineIndex)
	}
}

func checkOverlapsHorizontalForLayout(t *testing.T, layout *boxes.LayoutElement, line boxes.ConnectionLine, testIndex, lineIndex int) {
	if shouldHandle(layout) && (line.StartY == layout.Y || line.StartY == (layout.Y+layout.Height)) {
		// TODO
		// require.False(t, boxes.OverlapsHorizontally(line.StartX, line.EndX, layout.X, layout.X+layout.Width),
		// 	fmt.Sprintf("horizontal check(yLine=%d, yLayout=%d, yLayout2=%d), layoutId=%s: line.StartX=%d, line.EndX=%d, layout.X=%d, layout.X+width=%d, testIndex=%d, lineIndex=%d",
		// 		line.StartY, layout.Y, layout.Y+layout.Height, layout.Id, line.StartX, line.EndX, layout.X, layout.X+layout.Width, testIndex, lineIndex))
	}
	checkOverlapsHorizontalForLayoutCont(t, layout.Horizontal, line, testIndex, lineIndex)
	checkOverlapsHorizontalForLayoutCont(t, layout.Vertical, line, testIndex, lineIndex)
}

func checkOverlapsHorizontalForLayoutCont(t *testing.T, cont *boxes.LayoutElemContainer, line boxes.ConnectionLine, testIndex, lineIndex int) {
	if cont != nil {
		for _, layout := range cont.Elems {
			checkOverlapsHorizontalForLayout(t, &layout, line, testIndex, lineIndex)
		}
	}
}

func checkForRoadOverlapsVertically(t *testing.T, doc *boxes.BoxesDocument, testIndex int) {
	for lineIndex, line := range doc.VerticalRoads {
		checkOverlapsVerticalForLayout(t, &doc.Boxes, line, testIndex, lineIndex)
	}
}

func checkOverlapsVerticalForLayout(t *testing.T, layout *boxes.LayoutElement, line boxes.ConnectionLine, testIndex, lineIndex int) {
	if shouldHandle(layout) && (line.StartX == layout.X || line.StartX == (layout.X+layout.Width)) {
		require.False(t, boxes.OverlapsVertically(line.StartY, line.EndY, layout.Y, layout.Y+layout.Height),
			fmt.Sprintf("vertical check(xLine=%d, xLayout=%d, xLayout2=%d), layoutId=%s: line.StartY=%d, line.EndY=%d, layout.Y=%d, layout.Y+height=%d, testIndex=%d, lineIndex=%d",
				line.StartX, layout.X, layout.X+layout.Width, layout.Id, line.StartY, line.EndY, layout.Y, layout.Y+layout.Height, testIndex, lineIndex))
	}
	checkOverlapsVerticalForLayoutCont(t, layout.Horizontal, line, testIndex, lineIndex)
	checkOverlapsVerticalForLayoutCont(t, layout.Vertical, line, testIndex, lineIndex)
}

func checkOverlapsVerticalForLayoutCont(t *testing.T, cont *boxes.LayoutElemContainer, line boxes.ConnectionLine, testIndex, lineIndex int) {
	if cont != nil {
		for _, layout := range cont.Elems {
			checkOverlapsVerticalForLayout(t, &layout, line, testIndex, lineIndex)
		}
	}
}
