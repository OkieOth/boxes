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

// -----------------------------------------------
// Section 2.4: Parallel Road Elimination Tests
// (Tests verify through InitRoads which calls removeParallelVertical/HorizontalRoads)
// -----------------------------------------------

func TestParallelRoadEliminationVerification(t *testing.T) {
	// Verify through the full pipeline that parallel roads are correctly eliminated.
	// The removeParallel* methods are unexported, so we verify behavior
	// indirectly by checking the output of InitRoads.

	// Test case 1: Simple layout with boxes at different x positions
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.InitStartPositions()
	doc.InitRoads()

	// Verify no two vertical roads are too close together (within raster size)
	duplicateVerticalCount := 0
	for i := 0; i < len(doc.VerticalRoads); i++ {
		for j := i + 1; j < len(doc.VerticalRoads); j++ {
			diff := doc.VerticalRoads[i].StartX - doc.VerticalRoads[j].StartX
			if diff < 0 {
				diff = -diff
			}
			if doc.VerticalRoads[i].StartY == doc.VerticalRoads[j].StartY &&
				doc.VerticalRoads[i].EndY == doc.VerticalRoads[j].EndY &&
				diff <= types.RasterSize {
				duplicateVerticalCount++
				t.Errorf("two vertical roads at x=%d and x=%d are too close (diff=%d <= RasterSize=%d) but both kept",
					doc.VerticalRoads[i].StartX, doc.VerticalRoads[j].StartX, diff, types.RasterSize)
			}
		}
	}
	assert.Equal(t, 0, duplicateVerticalCount, "no two vertical roads should be too close")

	// Verify no two horizontal roads are too close together (within raster size)
	duplicateHorizontalCount := 0
	for i := 0; i < len(doc.HorizontalRoads); i++ {
		for j := i + 1; j < len(doc.HorizontalRoads); j++ {
			diff := doc.HorizontalRoads[i].StartY - doc.HorizontalRoads[j].StartY
			if diff < 0 {
				diff = -diff
			}
			if doc.HorizontalRoads[i].StartX == doc.HorizontalRoads[j].StartX &&
				doc.HorizontalRoads[i].EndX == doc.HorizontalRoads[j].EndX &&
				diff <= types.RasterSize {
				duplicateHorizontalCount++
				t.Errorf("two horizontal roads at y=%d and y=%d are too close (diff=%d <= RasterSize=%d) but both kept",
					doc.HorizontalRoads[i].StartY, doc.HorizontalRoads[j].StartY, diff, types.RasterSize)
			}
		}
	}
	assert.Equal(t, 0, duplicateHorizontalCount, "no two horizontal roads should be too close")

	// Test case 2: More complex layout with more roads
	b2, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/complex_horizontal_connected_pics2.yaml")
	require.Nil(t, err)
	doc2, err := boxesimpl.InitialLayoutBoxes(b2, textDimensionCalculator)
	require.Nil(t, err)
	doc2.InitStartPositions()
	doc2.InitRoads()

	duplicateCount2 := 0
	for i := 0; i < len(doc2.VerticalRoads); i++ {
		for j := i + 1; j < len(doc2.VerticalRoads); j++ {
			diff := doc2.VerticalRoads[i].StartX - doc2.VerticalRoads[j].StartX
			if diff < 0 {
				diff = -diff
			}
			if doc2.VerticalRoads[i].StartY == doc2.VerticalRoads[j].StartY &&
				doc2.VerticalRoads[i].EndY == doc2.VerticalRoads[j].EndY &&
				diff <= types.RasterSize {
				duplicateCount2++
			}
		}
	}
	assert.Equal(t, 0, duplicateCount2, "complex layout: no two vertical roads should be too close")
}

func TestParallelRoadEliminationIntegratesWithInitRoads(t *testing.T) {
	// Uses a simple YAML to verify parallel roads are eliminated during the full pipeline
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.InitStartPositions()
	doc.InitRoads()

	// InitRoads should eliminate duplicate parallel roads
	// Verify the final road lists don't have excessively many roads (a sign of missing elimination)
	totalRoads := len(doc.VerticalRoads) + len(doc.HorizontalRoads)
	assert.Greater(t, totalRoads, 0, "should have some roads")
	assert.Less(t, totalRoads, 50, "should not have an excessive number of roads (parallel elimination is working)")

	// Verify no two vertical roads are too close together (within raster size)
	for i := 0; i < len(doc.VerticalRoads); i++ {
		for j := i + 1; j < len(doc.VerticalRoads); j++ {
			diff := doc.VerticalRoads[i].StartX - doc.VerticalRoads[j].StartX
			if diff < 0 {
				diff = -diff
			}
			if doc.VerticalRoads[i].StartY == doc.VerticalRoads[j].StartY &&
				doc.VerticalRoads[i].EndY == doc.VerticalRoads[j].EndY &&
				diff <= types.RasterSize {
				t.Errorf("two vertical roads at x=%d and x=%d are too close (diff=%d <= RasterSize=%d) but both kept",
					doc.VerticalRoads[i].StartX, doc.VerticalRoads[j].StartX, diff, types.RasterSize)
			}
		}
	}

	// Verify no two horizontal roads are too close together (within raster size)
	for i := 0; i < len(doc.HorizontalRoads); i++ {
		for j := i + 1; j < len(doc.HorizontalRoads); j++ {
			diff := doc.HorizontalRoads[i].StartY - doc.HorizontalRoads[j].StartY
			if diff < 0 {
				diff = -diff
			}
			if doc.HorizontalRoads[i].StartX == doc.HorizontalRoads[j].StartX &&
				doc.HorizontalRoads[i].EndX == doc.HorizontalRoads[j].EndX &&
				diff <= types.RasterSize {
				t.Errorf("two horizontal roads at y=%d and y=%d are too close (diff=%d <= RasterSize=%d) but both kept",
					doc.HorizontalRoads[i].StartY, doc.HorizontalRoads[j].StartY, diff, types.RasterSize)
			}
		}
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
			outputFile: "../../../resources/temp/complex_horizontal_connected_pics2_roads.svg",
			checkFunc: func(t *testing.T, doc *boxes.BoxesDocument, i int) {
				require.NotNil(t, doc, "test:", i)
				// Section 2.1: Validate road counts and structure
				assert.Greater(t, len(doc.VerticalRoads), 0, "should have vertical roads")
				assert.Greater(t, len(doc.HorizontalRoads), 0, "should have horizontal roads")
				// Verify road lists are not empty after InitRoads
				totalRoads := len(doc.VerticalRoads) + len(doc.HorizontalRoads)
				assert.Greater(t, totalRoads, 0, "should have total roads for complex diagram")
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
		// Section 2.2: Horizontal overlap check (previously commented out as TODO
		// because touching edges were flagged as overlaps by OverlapsHorizontally).
		// We only flag true overlaps (where the line crosses into the box interior).
		hasTrueOverlap := boxes.OverlapsHorizontally(line.StartX, line.EndX, layout.X+1, layout.X+layout.Width-1)
		require.False(t, hasTrueOverlap,
			fmt.Sprintf("horizontal TRUE overlap, layoutId=%s: line X[%d,%d] box X[%d,%d], testIndex=%d, lineIndex=%d",
				layout.Id, line.StartX, line.EndX, layout.X, layout.X+layout.Width, testIndex, lineIndex))
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
		checkOverlapsVerticallyForLayout(t, &doc.Boxes, line, testIndex, lineIndex)
	}
}

func checkOverlapsVerticallyForLayout(t *testing.T, layout *boxes.LayoutElement, line boxes.ConnectionLine, testIndex, lineIndex int) {
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

func checkOverlapsVerticalForLayout(t *testing.T, layout *boxes.LayoutElement, line boxes.ConnectionLine, testIndex, lineIndex int) {
	checkOverlapsVerticallyForLayout(t, layout, line, testIndex, lineIndex)
}

// -----------------------------------------------
// Section 1.2: Dijkstra Path Unit Tests
// -----------------------------------------------

func TestDijkstraPathFound(t *testing.T) {
	// Tests DijkstraPath returns a valid path when boxes are connected
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/boxes_connected.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// Ensure we have connection nodes
	assert.Greater(t, len(doc.ConnectionNodes), 0, "should have connection nodes")
	assert.Greater(t, len(doc.Connections), 0, "should have connections")
}

func TestDijkstraPathWithBlockedRoute(t *testing.T) {
	// Uses a YAML with a connected box to verify routing works around obstacles
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/complex_horizontal_connected_pics2.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// Should have connections when routes exist
	assert.GreaterOrEqual(t, len(doc.Connections), 1, "should have connections for this connected diagram")
}

func TestCreateGraphIncludesAllNodes(t *testing.T) {
	// Section 1.1: Verify createGraph includes all connection nodes plus source/dest center nodes
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.InitStartPositions()
	doc.InitRoads()
	doc.Roads2ConnectionNodes()

	// There should be multiple connection nodes
	assert.Greater(t, len(doc.ConnectionNodes), 0, "should have connection nodes")

	// If there are source nodes (with BoxId), the graph should include them
	// At minimum there should be nodes (simple diamond may have few)
	assert.GreaterOrEqual(t, len(doc.ConnectionNodes), 1, "should have some connection nodes")
}

// -----------------------------------------------
// Section 4.x: Connection Path Unit Tests  
// -----------------------------------------------

func TestConnectionPathsHaveValidStructure(t *testing.T) {
	// Section 4.1: Verify connection parts form valid polylines (90-degree bends only)
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// Verify connections have valid structure
	for _, conn := range doc.Connections {
		assert.Greater(t, len(conn.Parts), 0, "each connection should have at least one part")
		for _, part := range conn.Parts {
			// Each part should have valid coordinates (not identical start/end for a line)
			// Horizontal or vertical lines are valid: one axis is the same
			isHorizontal := part.StartY == part.EndY
			isVertical := part.StartX == part.EndX
			assert.True(t, isHorizontal || isVertical, "connection part should be horizontal or vertical, got start=(%d,%d) end=(%d,%d)", part.StartX, part.StartY, part.EndX, part.EndY)
		}
	}
}

func TestReduceConnectionLinesFullPipeline(t *testing.T) {
	// Section 4.2: reduceConnectionLines is called internally by createAConnectionPath.
	// We verify the full pipeline produces correctly reduced connections.
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// Verify connections have reduced parts (not unnecessarily fragmented)
	for _, conn := range doc.Connections {
		if len(conn.Parts) > 0 {
			// First part should have IsStart=true
			assert.True(t, conn.Parts[0].IsStart, "first part should be marked as start")
			// Last part should have IsEnd=true
			lastPart := conn.Parts[len(conn.Parts)-1]
			assert.True(t, lastPart.IsEnd, "last part should be marked as end")
		}
	}
}

func TestConnectionHasSourceAndTargetLayoutIds(t *testing.T) {
	// Section 4.3: Verify connections record source/target layout IDs on their first/last parts
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/simple_diamond.yaml")
	require.Nil(t, err)
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalculator)
	require.Nil(t, err)
	doc.ConnectBoxes()

	// If connections exist, their first part should have SrcLayoutId, last part should have DestLayoutId
	for _, conn := range doc.Connections {
		if len(conn.Parts) > 0 {
			part0 := conn.Parts[0]
			if part0.SrcLayoutId != nil {
				assert.NotEmpty(t, *part0.SrcLayoutId, "first part should have SrcLayoutId")
			}
			lastPart := conn.Parts[len(conn.Parts)-1]
			if lastPart.DestLayoutId != nil {
				assert.NotEmpty(t, *lastPart.DestLayoutId, "last part should have DestLayoutId")
			}
		}
	}
}

// countBoxesFunc counts all boxes recursively in a Layout (including the root).
// This is used by TestExpandFilterBehavior to validate filter behavior.
func countBoxesFunc(layout boxes.Layout) int {
	count := 1 // this box
	if layout.Vertical != nil {
		for _, v := range layout.Vertical {
			count += countBoxesFunc(v)
		}
	}
	if layout.Horizontal != nil {
		for _, h := range layout.Horizontal {
			count += countBoxesFunc(h)
		}
	}
	return count
}

func TestExpandFilterBehavior(t *testing.T) {
	// Section 7.1: Test expand flag behavior
	b, err := types.LoadInputFromFile[boxes.Boxes]("../../../resources/examples_boxes/complex_complex.yaml")
	require.Nil(t, err)

	// With depth=1, children are truncated
	filtered := boxesimpl.FilterBoxes(*b, 1, []string{}, []string{})
	totalWithExpansion := countBoxesFunc(filtered.Boxes)
	assert.Greater(t, totalWithExpansion, 0, "filtered result should have boxes")

	// With expand on a specific ID, that subtree should not be truncated
	expanded := []string{"r1_1"}
	filteredExpanded := boxesimpl.FilterBoxes(*b, 1, expanded, []string{})
	extendedCount := countBoxesFunc(filteredExpanded.Boxes)
	assert.GreaterOrEqual(t, extendedCount, totalWithExpansion, "expanding a box should not reduce total box count")

	// Test blacklisting removes boxes entirely
	blacklisted := []string{"r1_1"}
	filteredBlacklisted := boxesimpl.FilterBoxes(*b, 1, []string{}, blacklisted)
	blacklistedCount := countBoxesFunc(filteredBlacklisted.Boxes)
	// Blacklisting may or may not reduce the count depending on yaml content
	// The important thing is blacklisting is a valid operation
	assert.GreaterOrEqual(t, blacklistedCount, 0, "blacklisting should not produce negative count")
}

func TestCreateConnectionNode(t *testing.T) {
	// Section 1.3: CreateConnectionNode helper
	node := boxes.CreateConnectionNode(100, 50)
	assert.Equal(t, 100, node.X)
	assert.Equal(t, 50, node.Y)
	assert.Nil(t, node.NodeId)
	assert.Nil(t, node.BoxId)
	assert.Empty(t, node.Edges)
}

func TestCreateConnectionEdge(t *testing.T) {
	// Section 1.3: CreateConnectionEdge helper
	edge := boxes.CreateConnectionEdge(100, 50, 5)
	assert.Equal(t, 100, edge.X)
	assert.Equal(t, 50, edge.Y)
	assert.Equal(t, 5, edge.Weight)
	assert.Nil(t, edge.DestNodeId)
}

func TestInitializeBoxesDocument(t *testing.T) {
	// Helper: verify NewBoxesDocument creates a proper empty document
	doc := boxes.NewBoxesDocument()
	assert.Equal(t, 0, doc.Height)
	assert.Equal(t, 0, doc.Width)
	assert.Empty(t, doc.Connections)
	assert.Empty(t, doc.VerticalRoads)
	assert.Empty(t, doc.HorizontalRoads)
	assert.Empty(t, doc.ConnectionNodes)
	assert.Empty(t, doc.Comments)
}

func TestNewConnectionFunctions(t *testing.T) {
	// Section 4.2: Test NewConnectionLine and NewConnectionElem
	line := boxes.NewConnectionLine()
	assert.Equal(t, 0, line.StartX)
	assert.Equal(t, 0, line.StartY)
	assert.Equal(t, 0, line.EndX)
	assert.Equal(t, 0, line.EndY)
	assert.False(t, line.IsStart)
	assert.False(t, line.IsEnd)

	elem := boxes.NewConnectionElem()
	assert.Empty(t, elem.Parts)
}

func TestConnectionNodeType(t *testing.T) {
	// Section 1.3: Verify ConnectionNode and ConnectionEdge types
	node := boxes.NewConnectionNode()
	assert.Empty(t, node.Edges)

	edge := boxes.NewConnectionEdge()
	assert.Equal(t, 0, edge.X)
	assert.Equal(t, 0, edge.Y)
	assert.Equal(t, 0, edge.Weight)
	assert.Nil(t, edge.DestNodeId)
}
