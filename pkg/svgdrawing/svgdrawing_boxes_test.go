package svgdrawing_test

import (
	"os"

	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
)

func TestSimpleSvg(t *testing.T) {
	tests := []struct {
		inputFile  string
		outputFile string
	}{
		{
			inputFile:  "../../resources/examples_boxes/simple_box.yaml",
			outputFile: "../../temp/TestSimpleSvg_box.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_box_nested.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_box_nested2.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested2.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_box_nested3.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested3.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_box_nested4.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested4.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_box_nested5.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested5.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_diamond.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/horizontal_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_vertical.yaml",
			outputFile: "../../temp/TestSimpleSvg_vcomplex.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_complex.yaml",
			outputFile: "../../temp/TestSimpleSvg_ccomplex.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/horizontal_nested_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond_nestedx.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/horizontal_nested_diamond2.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond_nestedx2.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/long_horizontal_vertical.yaml",
			outputFile: "../../temp/long_horizontal_vertical.svg",
		},
	}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		output, err := os.Create(test.outputFile)
		require.Nil(t, err)
		svgdrawing := svgdrawing.NewDrawing(output)
		svgdrawing.Start(doc.Title, doc.Height, doc.Width)
		svgdrawing.InitImages(doc.Images)
		svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
		doc.DrawBoxes(svgdrawing)
		svgdrawing.Done()
		output.Close()
		_, err = os.Stat(test.outputFile)
		require.Nil(t, err)
	}
}

func TestSvgWithConnections(t *testing.T) {
	type testData struct {
		inputFile  string
		outputFile string
		checkFunc  func(t *testing.T, doc *boxes.BoxesDocument)
	}
	runTests := func(tests []testData) {
		textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

		for _, test := range tests {
			b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
			require.Nil(t, err)
			doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
			require.Nil(t, err)
			doc.ConnectBoxes()
			doc.IncludeComments(textDimensionCalulator)
			doc.IncludeOverlays(textDimensionCalulator)
			output, err := os.Create(test.outputFile)
			require.Nil(t, err)
			svgdrawing := svgdrawing.NewDrawing(output)
			svgdrawing.Start(doc.Title, doc.Height, doc.Width)
			svgdrawing.InitImages(doc.Images)
			doc.DrawBoxes(svgdrawing)
			// DEBUG - Start
			//svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
			//doc.DrawRoads(svgdrawing)
			//doc.DrawStartPositions(svgdrawing)
			//doc.DrawConnectionNodes(svgdrawing)
			// DEBUG - End
			doc.DrawConnections(svgdrawing)
			doc.DrawComments(svgdrawing, textDimensionCalulator)

			svgdrawing.Done()
			output.Close()
			_, err = os.Stat(test.outputFile)
			require.Nil(t, err)
			test.checkFunc(t, doc)
		}
	}

	tests := []testData{
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_01.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_01.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_02.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_02.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_03.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_03.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_04.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_04.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_05.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_05.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/long_horizontal_01.yaml",
			outputFile: "../../temp/long_horizontal_01.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/long_horizontal_02.yaml",
			outputFile: "../../temp/long_horizontal_02.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/long_vertical_01.yaml",
			outputFile: "../../temp/long_vertical_01.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/long_vertical_02.yaml",
			outputFile: "../../temp/long_vertical_02.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/horizontal_nested_diamond2_connected.yaml",
			outputFile: "../../temp/horizontal_nested_diamond2_connected.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc)
			},
		},
		{
			inputFile:  "../../ui/data/boxes_random.yaml",
			outputFile: "../../temp/boxes_random.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_pics.yaml",
			outputFile: "../../temp/complex_horizontal_connected_pics.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/simple_diamond_connected.yaml",
			outputFile: "../../temp/simple_diamond_connected.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/boxes_random_truncated.yaml",
			outputFile: "../../temp/boxes_random_truncated.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_vertical.yaml",
			outputFile: "../../temp/TestSimpleSvg_vcomplex_connected.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_vertical2.yaml",
			outputFile: "../../temp/TestSimpleSvg_vcomplex_connected2.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/boxes_random_truncated2.yaml",
			outputFile: "../../temp/boxes_random_truncated2.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/boxes_random_truncated3.yaml",
			outputFile: "../../temp/boxes_random_truncated3.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_pics2.yaml",
			outputFile: "../../temp/complex_horizontal_connected_pics2.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)
				require.Equal(t, "r5_1", doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Id)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Caption)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Text1)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Text2)

				require.Equal(t, "r5_2", doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Id)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Caption)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Text1)
				require.Empty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Text2)
			},
		},
		{
			inputFile:  "../../resources/examples_boxes/complex_horizontal_connected_pics3.yaml",
			outputFile: "../../temp/complex_horizontal_connected_pics3.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument) {
				require.NotNil(t, doc.Connections)

				require.Equal(t, "r5_1", doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Id)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Caption)
				require.Empty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Text1)
				require.Empty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0].Text2)

				require.Equal(t, "r5_2", doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Id)
				require.NotEmpty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Caption)
				require.Empty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Text1)
				require.Empty(t, doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1].Text2)
			},
		},
	}
	runTests(tests)
}

func TestSplitTxt(t *testing.T) {
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	tests := []struct {
		name           string
		inputText      string
		fontDef        *types.FontDef
		expectedWidth  int
		expectedHeight int
		expectedLines  []types.TextAndDimensions
	}{
		{
			name:      "Single line sans-serif",
			inputText: "Hello World",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  62,
			expectedHeight: 12,
			expectedLines: []types.TextAndDimensions{
				{Text: "Hello World", Width: 62, Height: 12},
			},
		},
		{
			name:      "Multi line sans-serif",
			inputText: "Hello World Hello World",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  96,
			expectedHeight: 26,
			expectedLines: []types.TextAndDimensions{
				{Text: "Hello World Hello", Width: 96, Height: 14},
				{Text: "World", Width: 28, Height: 12},
			},
		},
		{
			name:      "Multi-line sans-serif",
			inputText: "This is a long text that should wrap into multiple lines",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 51,
			},
			expectedWidth:  51,
			expectedHeight: 96,
			expectedLines: []types.TextAndDimensions{
				{Text: "This is a", Width: 51, Height: 14},
				{Text: "long text", Width: 51, Height: 14},
				{Text: "that", Width: 22, Height: 14},
				{Text: "should", Width: 34, Height: 14},
				{Text: "wrap into", Width: 51, Height: 14},
				{Text: "multiple", Width: 45, Height: 14},
				{Text: "lines", Width: 28, Height: 12},
			},
		},
		{
			name:      "Single line monospace",
			inputText: "Monospace",
			fontDef: &types.FontDef{
				Font:              "monospace",
				Size:              10,
				MaxLenBeforeBreak: 80,
			},
			expectedWidth:  59,
			expectedHeight: 10,
			expectedLines: []types.TextAndDimensions{
				{Text: "Monospace", Width: 59, Height: 10},
			},
		},
		{
			name:      "Multi-line serif",
			inputText: "Serif font with multiple lines",
			fontDef: &types.FontDef{
				Font:              "serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  81,
			expectedHeight: 26,
			expectedLines: []types.TextAndDimensions{
				{Text: "Serif font with", Width: 81, Height: 14},
				{Text: "multiple lines", Width: 75, Height: 12},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, lines := textDimensionCalculator.SplitTxt(tt.inputText, tt.fontDef)
			require.Equal(t, tt.expectedWidth, width)
			require.Equal(t, tt.expectedHeight, height)
			for i, line := range lines {
				require.Equal(t, tt.expectedLines[i].Text, line.Text, i)
				require.Equal(t, tt.expectedLines[i].Width, line.Width, i)
				require.Equal(t, tt.expectedLines[i].Height, line.Height, i)
			}
		})
	}
}
