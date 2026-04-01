package boxesimpl_test

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/boxes/pkg/boxesimpl"
	"github.com/okieoth/boxes/pkg/svgdrawing"
	"github.com/okieoth/boxes/pkg/types"
	"github.com/okieoth/boxes/pkg/types/boxes"
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

func TestDrawBoxesFromFile(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		{
			inputFile: "../../resources/examples_boxes/simple_box.yaml",
		},
		{
			inputFile: "../../resources/examples_boxes/simple_diamond.yaml",
		}}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		require.NotNil(t, doc)
	}

}

// func tested types.InitDimensions
func TestInitDimensions(t *testing.T) {
	tests := []struct {
		layout         boxes.Layout
		expectedHeight int
		expectedWidth  int
	}{
		{
			layout: boxes.Layout{
				Caption: "test1",
			},
			expectedHeight: 100,
			expectedWidth:  150,
		},
		{
			// extends "test1" with an additional text1
			layout: boxes.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 160,
			expectedWidth:  150,
		},
		{
			// extends "test2" with an additional text2
			layout: boxes.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 220,
			expectedWidth:  150,
		},
		{
			// basic vertical layout test
			layout: boxes.Layout{
				Caption: "test4",
				Text1:   "test4-text1",
				Text2:   "test4-text2",
				Vertical: []boxes.Layout{
					{
						Caption: "test4-V1",
					},
					{
						Caption: "test4-V2",
					},
					{
						Caption: "test3-V3",
					},
				},
			},
			expectedHeight: 560,
			expectedWidth:  150,
		},
		{
			// basic horizontal layout test
			layout: boxes.Layout{
				Caption: "test5",
				Text1:   "test5-text1",
				Text2:   "test5-text2",
				Horizontal: []boxes.Layout{
					{
						Caption: "test5-V1",
					},
					{
						Caption: "test5-V2",
					},
					{
						Caption: "test5-V3",
					},
				},
			},
			expectedHeight: 320,
			expectedWidth:  530,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50)
	doc := boxes.NewBoxesDocument()
	for _, test := range tests {
		b := boxes.NewBoxes()
		le := boxesimpl.ExpInitLayoutElement(&test.layout, doc, b)
		le.InitDimensions(dc)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}

func TestLoadBoxesFromFile(t *testing.T) {
	tests := []struct {
		testFile        string
		shouldLoad      bool
		fileToCompareTo string
	}{
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext.yaml",
			shouldLoad:      true,
			fileToCompareTo: "../../resources/examples_boxes/horizontal_nested_diamond.yaml",
		},
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext2.yaml",
			shouldLoad:      true,
			fileToCompareTo: "../../resources/examples_boxes/horizontal_nested_diamond.yaml",
		},
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext2_fail.yaml",
			shouldLoad:      false,
			fileToCompareTo: "",
		},
	}
	for _, test := range tests {
		b, err := boxesimpl.LoadBoxesFromFile(test.testFile)
		if test.shouldLoad {
			require.Nil(t, err)
			require.NotNil(t, b)
			// Compare the loaded file with the original file
			b2, err := types.LoadInputFromFile[boxes.Boxes](test.fileToCompareTo)
			require.Nil(t, err)
			assert.Equal(t, b, b2, "The loaded boxes should match the expected boxes from the file")
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestDrawBoxesForUi(t *testing.T) {
	tests := []struct {
		inputFile   string
		outputFile  string
		depth       int
		expanded    []string
		blacklisted []string
	}{
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			outputFile:  "../../temp/complex_complex_filtered_1.svg",
			depth:       1,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			outputFile:  "../../temp/complex_complex_filtered_2.svg",
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile:   "../../ui/data/boxes_random.yaml",
			outputFile:  "../../temp/boxes_random_2.svg",
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		svgReturn := boxesimpl.DrawBoxesFiltered(*b, test.depth, test.expanded, test.blacklisted, true)

		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestDrawBoxesRelatedToConnections(t *testing.T) {
	tests := []struct {
		inputFile  string
		outputFile string
	}{
		{
			inputFile:  "../../resources/examples_boxes/boxes_random.yaml",
			outputFile: "../../temp/boxes_random_shrinked_1.svg",
		},
		{
			inputFile:  "../../resources/examples_boxes/boxes_random.yaml",
			outputFile: "../../temp/boxes_random_shrinked_2.svg",
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		svgReturn := boxesimpl.DrawBoxesRelatedToConnections(*b, []boxes.BoxesFileMixings{}, nil, true)
		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestDrawBoxesForUiComments(t *testing.T) {
	tests := []struct {
		inputFile    string
		outputFile   string
		depth        int
		hideComments bool
	}{
		{
			inputFile:    "../../resources/examples_boxes/boxes_random.yaml",
			outputFile:   "../../temp/boxes_random_with_notes.svg",
			depth:        2,
			hideComments: false,
		},
		{
			inputFile:    "../../resources/examples_boxes/boxes_random.yaml",
			outputFile:   "../../temp/boxes_random_no_notes.svg",
			depth:        2,
			hideComments: true,
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		svgReturn := boxesimpl.DrawBoxesFilteredComments(*b, test.depth, []string{}, []string{}, test.hideComments, true)

		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestDrawBoxesForUiExt(t *testing.T) {
	tests := []struct {
		inputFile           string
		inputExtConnections string
		inputExtFormats     string
		outputFile          string
		depth               int
		expanded            []string
		blacklisted         []string
	}{
		{
			inputFile:           "../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml",
			inputExtConnections: "../../resources/examples_boxes/ext_connections.yaml",
			inputExtFormats:     "../../resources/examples_boxes/ext_formats.yaml",
			outputFile:          "../../temp/ext_complex_horizontal_connected_pics.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_01.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_01.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_01_2.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_01_2.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_01_3.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_01_3.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_02.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_02.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_02_2.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_02_2.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_02_3.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_02_3.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_03.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_03.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_03_2.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_03_2.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_03_3.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_03_3.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_00_3.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_00_3.svg",
			depth:               10,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_simple_pic_00_4.yaml",
			inputExtConnections: "",
			inputExtFormats:     "../../resources/examples_boxes/boxes_simple_pic_format.yaml",
			outputFile:          "../../temp/boxes_simple_pic_00_4.svg",
			depth:               10,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/boxes_random.yaml",
			inputExtConnections: "",
			inputExtFormats:     "",
			outputFile:          "../../temp/boxes_random.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)

		mixins := make([]boxes.BoxesFileMixings, 0)
		if test.inputExtConnections != "" {
			extConnections, err := types.LoadInputFromFile[boxes.BoxesFileMixings](test.inputExtConnections)
			require.Nil(t, err)
			require.NotNil(t, extConnections)
			mixins = append(mixins, *extConnections)
		}

		if test.inputExtFormats != "" {
			extFormats, err := types.LoadInputFromFile[boxes.BoxesFileMixings](test.inputExtFormats)
			require.Nil(t, err)
			require.NotNil(t, extFormats)
			mixins = append(mixins, *extFormats)
		}

		svgReturn := boxesimpl.DrawBoxesFilteredExt(*b, mixins, nil, test.depth, test.expanded, test.blacklisted, false)

		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestDrawBoxesWithOverlays(t *testing.T) {
	tests := []struct {
		inputFile   string
		expandedIds []string
		mixins      []string
		outputFile  string
		maxdepth    int
	}{
		// {
		// 	inputFile: "../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml",
		// 	mixins: []string{
		// 		"../../resources/examples_boxes/ext_connections.yaml",
		// 		"../../resources/examples_boxes/ext_formats.yaml",
		// 		"../../resources/examples_boxes/ext_overlays.yaml",
		// 	},
		// 	maxdepth:   100,
		// 	outputFile: "../../temp/ext_complex_horizontal_connected_pics2.svg",
		// },
		{
			inputFile:   "../../resources/examples_boxes/boxes_connected.yaml",
			outputFile:  "../../temp/boxes_connected.svg",
			expandedIds: []string{"id2"},
			mixins:      []string{},
			maxdepth:    1,
		},
		{
			inputFile:   "../../resources/examples_boxes/boxes_connected_2.yaml",
			outputFile:  "../../temp/boxes_connected_2.svg",
			expandedIds: []string{},
			mixins:      []string{},
			maxdepth:    1,
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)

		mixins := make([]boxes.BoxesFileMixings, 0)
		for i := range test.mixins {
			mixin, err := types.LoadInputFromFile[boxes.BoxesFileMixings](test.mixins[i])
			require.Nil(t, err)
			require.NotNil(t, mixin)
			mixins = append(mixins, *mixin)
		}

		svgReturn := boxesimpl.DrawBoxesFilteredExt(*b, mixins, nil, test.maxdepth, test.expandedIds, []string{}, false)

		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestFilterBoxes(t *testing.T) {
	tests := []struct {
		inputFile   string
		checkFunc   func(b *boxes.Boxes)
		depth       int
		expanded    []string
		blacklisted []string
	}{
		{
			inputFile: "../../resources/examples_boxes/complex_complex.yaml",
			checkFunc: func(b *boxes.Boxes) {
				for _, e := range b.Boxes.Horizontal {
					require.Equal(t, 0, len(e.Horizontal), "got unexpected horizontal childs (1-1)")
					require.Equal(t, 0, len(e.Vertical), "got unexpected vertical childs (1-1)")
				}
				for _, e := range b.Boxes.Vertical {
					require.Equal(t, 0, len(e.Horizontal), "got unexpected horizontal childs (1-2)")
					require.Equal(t, 0, len(e.Vertical), "got unexpected vertical childs (1-2)")
				}
			},
			depth:       1,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile: "../../resources/examples_boxes/complex_complex.yaml",
			checkFunc: func(b *boxes.Boxes) {
				found := false
				for _, e := range b.Boxes.Horizontal {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				for _, e := range b.Boxes.Vertical {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				require.True(t, found, "didn't find second level")
			},
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile: "../../resources/examples_boxes/complex_complex.yaml",
			checkFunc: func(b *boxes.Boxes) {
				found := false
				for _, e := range b.Boxes.Horizontal {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				for _, e := range b.Boxes.Vertical {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				require.True(t, found, "didn't find second level")
			},
			depth:       20,
			expanded:    []string{},
			blacklisted: []string{"r2_2", "r4_1"},
		},
		{
			inputFile: "../../ui/data/boxes_random.yaml",
			checkFunc: func(b *boxes.Boxes) {
				require.NotNil(t, b)
				require.Len(t, b.Boxes.Vertical[0].Horizontal[0].Vertical, 0)
				require.Len(t, b.Boxes.Vertical[0].Horizontal[0].Horizontal, 0)
				require.Len(t, b.Boxes.Vertical[0].Horizontal[0].Connections, 3)
			},
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		filtered := boxesimpl.FilterBoxes(*b, test.depth, test.expanded, test.blacklisted)
		test.checkFunc(&filtered)
	}
}
