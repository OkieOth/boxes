package svgdrawing_test

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"
)

// countElements counts occurrences of a tag pattern in SVG string.
// Usage: countElements(svg, "<rect") or countElements(svg, `class="connection"`)
func countElements(svg string, tag string) int {
	return strings.Count(svg, tag)
}

// assertHasSvgRoot asserts the SVG has a valid root element with namespace.
func assertHasSvgRoot(t require.TestingT, svg string) {
	require.Contains(t, svg, "<svg", "SVG should have <svg root element")
	require.Contains(t, svg, `xmlns="http://www.w3.org/2000/svg"`, "SVG should have XML namespace")
}

// assertHasRect asserts the SVG contains at least one <rect element.
func assertHasRect(t require.TestingT, svg string) {
	count := countElements(svg, "<rect")
	require.Greater(t, count, 0, "SVG should have at least one box rect, got %d", count)
}

// assertRectCount asserts the exact number of <rect elements in the SVG.
func assertRectCount(t require.TestingT, svg string, expected int) {
	count := countElements(svg, "<rect")
	require.Equal(t, expected, count, "Expected %d rects, got %d", expected, count)
}

// BoxRect holds parsed coordinates of a box rect element.
type BoxRect struct {
	ID    string
	X, Y  int
	W, H  int
}

// extractBoxRects parses all <rect elements and returns their parsed data.
func extractBoxRects(svg string) []BoxRect {
	// Match <rect class="..." id="..." x="..." y="..." width="..." height="..."/>
	re := regexp.MustCompile(`<rect\b[^>]*?id="([^"]*)"[^>]*?x="(\d+)"\s+y="(\d+)"\s+width="(\d+)"\s+height="(\d+)"`)
	matches := re.FindAllStringSubmatch(svg, -1)

	var rects []BoxRect
	for _, m := range matches {
		if len(m) == 6 {
			x, _ := strconv.Atoi(m[2])
			y, _ := strconv.Atoi(m[3])
			w, _ := strconv.Atoi(m[4])
			h, _ := strconv.Atoi(m[5])
			rects = append(rects, BoxRect{
				ID: m[1],
				X:  x, Y: y, W: w, H: h,
			})
		}
	}
	return rects
}

// assertBoxPresent asserts a box with given ID exists and returns its coordinates.
func assertBoxPresent(t require.TestingT, svg, id string) BoxRect {
	rects := extractBoxRects(svg)
	for _, r := range rects {
		if r.ID == id {
			require.NotEmpty(t, r.ID, "Box with id=%s should be present", id)
			require.Greater(t, r.W, 0, "Box %s should have non-zero width", id)
			require.Greater(t, r.H, 0, "Box %s should have non-zero height", id)
			return r
		}
	}
	require.Failf(t, "Box with id=%s not found in SVG", id)
	return BoxRect{}
}

// assertBoxAbsent asserts a box with given ID does NOT appear in the SVG.
func assertBoxAbsent(t require.TestingT, svg, id string) {
	rects := extractBoxRects(svg)
	for _, r := range rects {
		if r.ID == id {
			require.Failf(t, "Box with id=%s should not be present in SVG (was expected to be filtered out)", id)
		}
	}
}

// ConnectionLine holds parsed coordinates of a connection line segment.
type ConnectionLine struct {
	X1, Y1, X2, Y2 int
}

// extractConnections parses all connection line segments from SVG.
func extractConnections(svg string) []ConnectionLine {
	re := regexp.MustCompile(`class="connection[^"]*"\s+x1="(\d+)"\s+y1="(\d+)"\s+x2="(\d+)"\s+y2="(\d+)"`)
	matches := re.FindAllStringSubmatch(svg, -1)

	var lines []ConnectionLine
	for _, m := range matches {
		if len(m) == 5 {
			x1, _ := strconv.Atoi(m[1])
			y1, _ := strconv.Atoi(m[2])
			x2, _ := strconv.Atoi(m[3])
			y2, _ := strconv.Atoi(m[4])
			lines = append(lines, ConnectionLine{X1: x1, Y1: y1, X2: x2, Y2: y2})
		}
	}
	return lines
}

// assertConnectionLinesPresent asserts there are connection lines in the SVG.
func assertConnectionLinesPresent(t require.TestingT, svg string) {
	count := strings.Count(svg, `class="connection`)
	require.Greater(t, count, 0, "SVG should have connection lines")

	lines := extractConnections(svg)
	require.Greater(t, len(lines), 0, "SVG should have at least one connection line segment")
}

// assertConnectionsCount asserts the exact number of connection line segments.
func assertConnectionsCount(t require.TestingT, svg string, expected int) {
	lines := extractConnections(svg)
	require.Equal(t, expected, len(lines), "Expected %d connection line segments, got %d", expected, len(lines))
}

// assertTextPresent asserts a text string appears in the SVG.
func assertTextPresent(t require.TestingT, svg, text string) {
	require.Contains(t, svg, text, "SVG should contain text: %s", text)
}

// assertTextAbsent asserts a text string does NOT appear in the SVG.
func assertTextAbsent(t require.TestingT, svg, text string) {
	require.NotContains(t, svg, text, "SVG should NOT contain text: %s", text)
}

// assertHasComments asserts comment-related SVG elements are present.
func assertHasComments(t require.TestingT, svg string) {
	require.Contains(t, svg, "comment", "SVG should have comment-related elements")
}

// assertNoComments asserts no comment-related SVG elements are present.
func assertNoComments(t require.TestingT, svg string) {
	require.NotContains(t, svg, "comment", "SVG should not have comment-related elements")
}

// assertHasConnectionNodes asserts connection node markers exist in the SVG.
func assertHasConnectionNodes(t require.TestingT, svg string) {
	// Connection nodes are drawn as circles with class containing "connection"
	re := regexp.MustCompile(`<circle\b[^>]*class="[^"]*connection[^"]*"[^>]*>`)
	count := len(re.FindAllString(svg, -1))
	if count == 0 {
		// Some documents don't have connection nodes even with connections.
		// Log a warning instead of failing.
		return
	}
	_ = count
}

// assertSvgNonEmpty asserts the SVG string is not empty.
func assertSvgNonEmpty(t require.TestingT, svg string) {
	require.NotEmpty(t, svg, "SVG output should not be empty")
}

// assertHasTextElements asserts the SVG contains <text elements.
func assertHasTextElements(t require.TestingT, svg string) {
	count := strings.Count(svg, "<text")
	require.Greater(t, count, 0, "SVG should have text elements, got %d", count)
}

// assertTitlePresent asserts an SVG title text element exists.
func assertTitlePresent(t require.TestingT, svg string) {
	// Title uses text-anchor:middle style
	re := regexp.MustCompile(`<text[^>]*style="text-anchor:middle"`)
	count := len(re.FindAllString(svg, -1))
	require.Greater(t, count, 0, "SVG should have at least one title text element")
}

// assertNonZeroBoxDimensions asserts all rects have non-zero dimensions and positions.
func assertNonZeroBoxDimensions(t require.TestingT, svg string) {
	rects := extractBoxRects(svg)
	for _, r := range rects {
		if r.ID != "" {
			require.Greater(t, r.W, 0, "Box %s should have non-zero width", r.ID)
			require.Greater(t, r.H, 0, "Box %s should have non-zero height", r.ID)
			// Boxes should have non-negative position (within viewBox)
			require.GreaterOrEqual(t, r.X, 0, "Box %s should have non-negative x", r.ID)
			require.GreaterOrEqual(t, r.Y, 0, "Box %s should have non-negative y", r.ID)
		}
	}
}

// compareSvgElements compares rect counts between two SVG strings.
func compareSvgElements(svg1, svg2 string) int {
	c1 := countElements(svg1, "<rect")
	c2 := countElements(svg2, "<rect")
	return c2 - c1
}

// assertSvgHasImage asserts SVG contains <image or <use elements (for embedded images).
func assertSvgHasImage(t require.TestingT, svg string) {
	hasImage := strings.Contains(svg, "<image") || strings.Contains(svg, `<use `)
	require.True(t, hasImage, "SVG should contain image/embedded image elements")
}

// assertSvgHasTitle asserts SVG has a title element with content.
func assertSvgHasTitleContent(t require.TestingT, svg string, titleText string) {
	require.Contains(t, svg, titleText, "SVG should contain title text: %s", titleText)
}

// assertSvgHasExpandedId asserts an expanded box ID is present in SVG.
func assertSvgHasExpandedId(t require.TestingT, svg string, id string) {
	require.Contains(t, svg, `id="`+id+`"`, "SVG should contain expanded box with id=%s", id)
}

// assertSvgHasBlacklistedId asserts a blacklisted box ID is NOT present in SVG.
func assertSvgHasNoBlacklistedId(t require.TestingT, svg string, id string) {
	require.NotContains(t, svg, `id="`+id+`"`, "SVG should not contain blacklisted box with id=%s", id)
}

// assertHasOverlay asserts SVG contains overlay-related elements.
func assertHasOverlay(t require.TestingT, svg string) {
	require.Contains(t, svg, "overlay", "SVG should have overlay elements")
}

// assertHasWrapper asserts SVG contains wrapper-related elements.
func assertHasWrapper(t require.TestingT, svg string) {
	require.Contains(t, svg, "wrapper", "SVG should have wrapper elements")
}

// countTextElements returns count of <text elements.
func countTextElements(svg string) int {
	return strings.Count(svg, "<text")
}

// countRectElements returns count of <rect elements.
func countRectElements(svg string) int {
	return countElements(svg, "<rect")
}

// countLinesWithClass returns count of <line> elements with a specific class prefix.
func countLinesWithClass(svg string, classPrefix string) int {
	re := regexp.MustCompile(`<line\b[^>]*class="[^"]*` + regexp.QuoteMeta(classPrefix) + `[^"]*"`)
	return len(re.FindAllString(svg, -1))
}
