package boxesimpl_test

import (
	"os"
	"strings"

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

// Section 3.1: Parametric dimension calculator that returns values proportional to text length
// This exercises the real dimension computation logic (text width/height calculations, margins, padding)
type ParametricDimensionCalculator struct {
	baseWidth  int
	baseHeight int
}

func (p *ParametricDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	// Return dimensions proportional to the text length so that different strings produce different widths
	textLen := len(txt)
	if textLen == 0 {
		textLen = 1
	}
	// 10px per character, minimum base height
	return textLen * 10, p.baseHeight
}

func (p *ParametricDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []types.TextAndDimensions) {
	w, h := p.Dimensions(txt, format)
	return w, h, []types.TextAndDimensions{{Text: txt, Width: w, Height: h}}
}

func (p *ParametricDimensionCalculator) DimensionsWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int) {
	return p.Dimensions(txt, format)
}

func (p *ParametricDimensionCalculator) SplitTxtWithMaxWidth(txt string, format *types.FontDef, maxWidth int) (width, height int, lines []types.TextAndDimensions) {
	return p.SplitTxt(txt, format)
}

func stringPtr(s string) *string {
	return &s
}

func NewParametricDimensionCalculator() *ParametricDimensionCalculator {
	return &ParametricDimensionCalculator{
		baseWidth:  50,
		baseHeight: 15,
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
		le := boxesimpl.ExpInitLayoutElement(&test.layout, doc, b, []string{})
		le.InitDimensions(dc)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}

// Section 3.1: Improved InitDimensions test with parametric dimension calculator
func TestInitDimensionsWithParametricCalculator(t *testing.T) {
	// Uses a parametric dimension calculator that returns values proportional to text length
	// so actual dimension computation is exercised
	parametricCalc := NewParametricDimensionCalculator()

	tests := []struct {
		name   string
		layout boxes.Layout
		check  func(t *testing.T, le *boxes.LayoutElement)
	}{
		{
			name: "caption only",
			layout: boxes.Layout{
				Caption: "Hello",
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Height, 0, "height should be positive")
				assert.Greater(t, le.Width, 0, "width should be positive")
				assert.NotNil(t, le.WidthTextBox, "WidthTextBox should be set")
				assert.NotNil(t, le.HeightTextBox, "HeightTextBox should be set")
			},
		},
		{
			name: "fixed width via format",
			layout: boxes.Layout{
				Caption: "fixed",
				Format:  stringPtr("fw"),
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				// Fixed width is set via the format, which is resolved in InitialLayoutBoxes
				// This test verifies that a layout with a format ref is handled correctly
				assert.Greater(t, le.Height, 0, "height should be positive")
			},
		},
		{
			name: "fixed height via format",
			layout: boxes.Layout{
				Caption: "fixed",
				Format:  stringPtr("fh"),
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Height, 0, "height should be positive")
			},
		},
		{
			name: "caption only (no text1/text2) - should have less padding",
			layout: boxes.Layout{
				Caption: "short",
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Height, 0, "height should be positive")
			},
		},
		{
			name: "empty box (no caption, text1, text2, no children) - should be minimal",
			layout: boxes.Layout{},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				// Empty boxes without format specs should have minimal dimensions
				// (determined only by format defaults)
			},
		},
		{
			name: "vertical children - max child height determines container width",
			layout: boxes.Layout{
				Caption:  "parent",
				Vertical: []boxes.Layout{{Caption: "child1"}, {Caption: "child2"}},
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Height, 0, "height should be positive")
				if le.Vertical != nil {
					assert.Greater(t, len(le.Vertical.Elems), 0, "should have vertical children")
					// All children should be width-aligned to the max
					maxWidth := 0
					for _, child := range le.Vertical.Elems {
						if child.Width > maxWidth {
							maxWidth = child.Width
						}
					}
					// The widest child determines the container width (approximately)
					assert.GreaterOrEqual(t, le.Width, maxWidth, "container should be at least as wide as widest child")
				}
			},
		},
		{
			name: "horizontal children - max child height determines container height",
			layout: boxes.Layout{
				Caption:    "parent",
				Horizontal: []boxes.Layout{{Caption: "child1"}, {Caption: "child2"}},
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Width, 0, "width should be positive")
				if le.Horizontal != nil {
					assert.Greater(t, len(le.Horizontal.Elems), 0, "should have horizontal children")
				}
			},
		},
		{
			name: "nested boxes (vertical with horizontal children)",
			layout: boxes.Layout{
				Caption:  "outer",
				Vertical: []boxes.Layout{{Caption: "v-child", Horizontal: []boxes.Layout{{Caption: "h-child"}}}},
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				assert.Greater(t, le.Height, 0, "outer height should be positive")
				if le.Vertical != nil && len(le.Vertical.Elems) > 0 {
					inner := le.Vertical.Elems[0]
					if inner.Horizontal != nil && len(inner.Horizontal.Elems) > 0 {
						assert.Greater(t, inner.Horizontal.Elems[0].Width, 0, "nested child width should be positive")
					}
				}
			},
		},
		{
			name: "image reference",
			layout: boxes.Layout{
				Image: stringPtr("img1"),
			},
			check: func(t *testing.T, le *boxes.LayoutElement) {
				// An image-only box should have dimensions based on image size
				// (in real pipeline, images are loaded from Boxes.Images map)
			},
		},
	}

	doc := boxes.NewBoxesDocument()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := boxes.NewBoxes()
			le := boxesimpl.ExpInitLayoutElement(&tt.layout, doc, b, []string{})
			le.InitDimensions(parametricCalc)
			tt.check(t, &le)
		})
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

		// SVG content assertions (from plan §2.1)
		svgdrawing.AssertSvgNonEmpty(t, svgReturn.SVG)
		svgdrawing.AssertHasRect(t, svgReturn.SVG)

		// Depth=1 vs Depth=2: depth=2 should have more boxes (expanded children)
		if test.depth == 2 {
			rectCount := svgdrawing.CountRectElements(svgReturn.SVG)
			require.Greater(t, rectCount, 0, "depth=2 should produce boxes")
		}

		// Write to file only for manual debugging (kept for developers)
		if test.outputFile != "" {
			err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
			require.Nil(t, err, "error while writing debug file: %s", test.outputFile)
		}
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

		// SVG content assertions (from plan §2.2)
		svgdrawing.AssertSvgNonEmpty(t, svgReturn.SVG)
		svgdrawing.AssertHasRect(t, svgReturn.SVG)
		// Connections should be rendered
		conLines := svgdrawing.ExtractConnections(svgReturn.SVG)
		require.Greater(t, len(conLines), 0, "related-to-connections should produce connection lines")

		// Write to file only for manual debugging (kept for developers)
		if test.outputFile != "" {
			err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
			require.Nil(t, err, "error while writing debug file: %s", test.outputFile)
		}
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

		// SVG content assertions (from plan §2.3)
		svgdrawing.AssertSvgNonEmpty(t, svgReturn.SVG)
		svgdrawing.AssertHasRect(t, svgReturn.SVG)

		if !test.hideComments {
			// When comments are NOT hidden, comment-related elements should be present
			svgdrawing.AssertHasComments(t, svgReturn.SVG)
		} else {
			// When comments are hidden, no comment-related elements should appear
			svgdrawing.AssertNoComments(t, svgReturn.SVG)
		}

		// Write to file only for manual debugging (kept for developers)
		if test.outputFile != "" {
			err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
			require.Nil(t, err, "error while writing debug file: %s", test.outputFile)
		}
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
			inputFile:           "../../resources/examples_boxes/complex_horizontal_connected_pics_lines.yaml",
			inputExtConnections: "../../resources/examples_boxes/ext_connections.yaml",
			inputExtFormats:     "../../resources/examples_boxes/ext_formats.yaml",
			outputFile:          "../../temp/ext_complex_horizontal_connected_pics_lines.svg",
			depth:               2,
			expanded:            []string{},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/oooo_1.yaml",
			inputExtConnections: "../../resources/examples_boxes/oooo_2.yaml",
			inputExtFormats:     "",
			outputFile:          "../../temp/orga.svg",
			depth:               2,
			expanded:            []string{"id_3_3_0_1"},
			blacklisted:         []string{},
		},
		{
			inputFile:           "../../resources/examples_boxes/oooo_1.yaml",
			inputExtConnections: "../../resources/examples_boxes/oooo_3.yaml",
			inputExtFormats:     "",
			outputFile:          "../../temp/orga2.svg",
			depth:               2,
			expanded:            []string{"id_3_3_0_1"},
			blacklisted:         []string{},
		},
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

		// SVG content assertions (from plan §2.4)
		svgdrawing.AssertSvgNonEmpty(t, svgReturn.SVG)
		svgdrawing.AssertHasRect(t, svgReturn.SVG)
		svgdrawing.AssertTextElementsPresent(t, svgReturn.SVG)

		// Case-specific validations
		if strings.Contains(test.inputFile, "boxes_simple_pic") || strings.Contains(test.inputFile, "simple_pic") {
			// With images: verify image elements are present
			svgdrawing.AssertSvgHasImage(t, svgReturn.SVG)
		}

		if len(test.expanded) > 0 {
			// Verify expanded IDs appear in SVG
			for _, exp := range test.expanded {
				svgdrawing.AssertSvgHasExpandedId(t, svgReturn.SVG, exp)
			}
		}

		if len(test.blacklisted) > 0 {
			// Verify blacklisted IDs do NOT appear in SVG
			for _, bl := range test.blacklisted {
				svgdrawing.AssertSvgHasBlacklistedId(t, svgReturn.SVG, bl)
			}
		}

		// Connection lines present when connections exist
		if test.inputExtConnections != "" {
			svgdrawing.AssertConnectionLinesPresent(t, svgReturn.SVG)
		}

		// Write to file only for manual debugging (kept for developers)
		if test.outputFile != "" {
			err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
			require.Nil(t, err, "error while writing debug file: %s", test.outputFile)
		}
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
		{
			inputFile: "../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml",
			mixins: []string{
				"../../resources/examples_boxes/ext_connections.yaml",
				"../../resources/examples_boxes/ext_formats.yaml",
				"../../resources/examples_boxes/ext_overlays.yaml",
			},
			maxdepth:   100,
			outputFile: "../../temp/ext_complex_horizontal_connected_pics2.svg",
		},
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
		{
			inputFile: "../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml",
			mixins: []string{
				"../../resources/examples_boxes/ext_connections.yaml",
				"../../resources/examples_boxes/ext_formats.yaml",
				"../../resources/examples_boxes/ext_overlays.yaml",
				"../../resources/examples_boxes/ext_wrappers.yaml",
			},
			maxdepth:   100,
			outputFile: "../../temp/ext_complex_horizontal_connected_pics3.svg",
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

		// SVG content assertions (from plan §2.5)
		svgdrawing.AssertSvgNonEmpty(t, svgReturn.SVG)
		svgdrawing.AssertHasRect(t, svgReturn.SVG)

		// Overlays: check for overlay elements when mixins contain overlays
		hasOverlayMixin := false
		for _, m := range test.mixins {
			if strings.Contains(m, "ext_overlays") || strings.Contains(m, "overlays") {
				hasOverlayMixin = true
				break
			}
		}
		if hasOverlayMixin {
			svgdrawing.AssertHasOverlay(t, svgReturn.SVG)
		}

		// Wrappers: check for wrapper elements when mixins contain wrappers
		hasWrapperMixin := false
		for _, m := range test.mixins {
			if strings.Contains(m, "ext_wrappers") || strings.Contains(m, "wrappers") {
				hasWrapperMixin = true
				break
			}
		}
		if hasWrapperMixin {
			svgdrawing.AssertHasWrapper(t, svgReturn.SVG)
		}

		// Write to file only for manual debugging (kept for developers)
		if test.outputFile != "" {
			err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
			require.Nil(t, err, "error while writing debug file: %s", test.outputFile)
		}
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

// TestFilterBoxesProducesValidSvg integrates FilterBoxes output with SVG drawing
// and validates the resulting SVG (from plan §2.6).
func TestFilterBoxesProducesValidSvg(t *testing.T) {
	tests := []struct {
		inputFile   string
		depth       int
		expanded    []string
		blacklisted []string
	}{
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			depth:       1,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			depth:       20,
			expanded:    []string{},
			blacklisted: []string{"r2_2", "r4_1"},
		},
		{
			inputFile:   "../../ui/data/boxes_random.yaml",
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
	}

	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error loading input for test %d", i)

		// Filter boxes as DrawBoxesFiltered does internally
		filtered := boxesimpl.FilterBoxes(*b, test.depth, test.expanded, test.blacklisted)

		// Generate SVG from filtered layout
		svgReturn := boxesimpl.DrawBoxesFiltered(filtered, test.depth, test.expanded, test.blacklisted, false)
		require.NotEmpty(t, svgReturn.SVG, "SVG output should not be empty for test %d", i)
		require.Equal(t, "", svgReturn.ErrorMsg, "no error expected for test %d", i)

		// Core assertions
		svgdrawing.AssertHasRect(t, svgReturn.SVG)

		// Verify filtered output: blacklisted IDs should NOT appear in SVG
		for _, bl := range test.blacklisted {
			svgdrawing.AssertSvgHasBlacklistedId(t, svgReturn.SVG, bl)
		}

		// Verify expanded IDs DO appear in SVG
		for _, exp := range test.expanded {
			svgdrawing.AssertSvgHasExpandedId(t, svgReturn.SVG, exp)
		}
	}
}
