package boxes_test

import (
	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
	"github.com/stretchr/testify/require"
)

type DummyDimensionCalculator struct {
	width  int
	height int
}

func (d *DummyDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.width, d.height
}

func (d *DummyDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []types.TextAndDimensions) {
	return 0, 0, make([]types.TextAndDimensions, 0)
}

func (d *DummyDimensionCalculator) DimensionsWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int) {
	return d.width, d.height
}

func (d *DummyDimensionCalculator) SplitTxtWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int, lines []types.TextAndDimensions) {
	return 0, 0, make([]types.TextAndDimensions, 0)
}

func NewDummyDimensionCalculator(width, height int) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		width:  width,
		height: height,
	}
}

func checkLayoutElement(t *testing.T, le *boxes.LayoutElement, initX, initY int) {
	require.GreaterOrEqual(t, le.X, initX)
	require.GreaterOrEqual(t, le.Y, initY)
	require.Greater(t, le.CenterX, 0)
	require.Greater(t, le.CenterY, 0)
	require.Equal(t, le.CenterX, le.X+(le.Width/2))
	require.Equal(t, le.CenterY, le.Y+(le.Height/2))
	if le.Vertical != nil {
		require.GreaterOrEqual(t, le.Vertical.X, le.X)
		require.GreaterOrEqual(t, le.Vertical.Y, le.Y)
		for _, v := range le.Vertical.Elems {
			checkLayoutElement(t, &v, le.X, le.Y)
		}
	}
	if le.Horizontal != nil {
		require.GreaterOrEqual(t, le.Horizontal.X, le.X)
		require.GreaterOrEqual(t, le.Horizontal.Y, le.Y)
		for _, h := range le.Horizontal.Elems {
			checkLayoutElement(t, &h, le.X, le.Y)
		}
	}
}

func TestCenteredCoordinates(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		{
			inputFile: "../../../resources/examples_boxes/simple_diamond.yaml",
		},
		{
			inputFile: "../../../resources/examples_boxes/horizontal_diamond.yaml",
		},
		{
			inputFile: "../../../resources/examples_boxes/complex_vertical.yaml",
		},
		{
			inputFile: "../../../resources/examples_boxes/complex_horizontal.yaml",
		},
	}
	dc := NewDummyDimensionCalculator(100, 50)
	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, dc)
		require.Nil(t, err)
		checkLayoutElement(t, &doc.Boxes, 0, 0)
	}

}

func TestAreOnTheSameVerticalLevel(t *testing.T) {
	tests := []struct {
		inputFile string
		checkFunc func(t *testing.T, doc *boxes.BoxesDocument)
	}{
		{
			inputFile: "../../../resources/examples_boxes/simple_diamond.yaml",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc)
				e1 := doc.Boxes.Vertical.Elems[0]
				e2 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[0]
				e3 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[1]
				e4 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[2]
				e5 := doc.Boxes.Vertical.Elems[2]
				require.False(t, e1.AreOnTheSameVerticalLevel(&e2))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e4))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e5))

				require.True(t, e2.AreOnTheSameVerticalLevel(&e3))
				require.True(t, e2.AreOnTheSameVerticalLevel(&e4))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e4))
				require.True(t, e4.AreOnTheSameVerticalLevel(&e2))
				require.True(t, e4.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e2))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e4))
			},
		},
		{
			inputFile: "../../../resources/examples_boxes/horizontal_diamond.yaml",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc)
				e1 := doc.Boxes.Horizontal.Elems[0]
				e2 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0]
				e3 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1]
				e4 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[2]
				e5 := doc.Boxes.Horizontal.Elems[2]
				require.False(t, e1.AreOnTheSameVerticalLevel(&e2))
				require.True(t, e1.AreOnTheSameVerticalLevel(&e3))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e1))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e4))
				require.False(t, e4.AreOnTheSameVerticalLevel(&e1))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e5))
			},
		},
	}
	dc := NewDummyDimensionCalculator(100, 50)
	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, dc)
		require.Nil(t, err)
		test.checkFunc(t, doc)
	}

}
