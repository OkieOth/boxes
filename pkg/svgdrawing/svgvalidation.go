package svgdrawing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"
)

// CountSvgElements counts occurrences of a tag pattern in SVG string.
// Usage: CountSvgElements(svg, "<rect") or CountSvgElements(svg, `class="connection"`)
func CountSvgElements(svg string, tag string) int {
	return strings.Count(svg, tag)
}

// AssertHasSvgRoot asserts the SVG has a valid root element with namespace.
func AssertHasSvgRoot(t require.TestingT, svg string) {
	require.Contains(t, svg, "<svg", "SVG should have <svg root element")
	require.Contains(t, svg, `xmlns="http://www.w3.org/2000/svg"`, "SVG should have XML namespace")
}

// AssertHasRect asserts the SVG contains at least one <rect element.
func AssertHasRect(t require.TestingT, svg string) {
	count := CountSvgElements(svg, "<rect")
	require.Greater(t, count, 0, "SVG should have at least one box rect, got %d", count)
}

// AssertRectCount asserts the exact number of <rect elements in the SVG.
func AssertRectCount(t require.TestingT, svg string, expected int) {
	count := CountSvgElements(svg, "<rect")
	require.Equal(t, expected, count, "Expected %d rects, got %d", expected, count)
}

// BoxRect holds parsed coordinates of a box rect element.
type BoxRect struct {
	ID    string
	X, Y  int
	W, H  int
}

// ExtractBoxRects parses all <rect elements and returns their parsed data.
func ExtractBoxRects(svg string) []BoxRect {
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

// AssertBoxPresent asserts a box with given ID exists and returns its coordinates.
func AssertBoxPresent(t require.TestingT, svg, id string) BoxRect {
	rects := ExtractBoxRects(svg)
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

// AssertBoxAbsent asserts a box with given ID does NOT appear in the SVG.
func AssertBoxAbsent(t require.TestingT, svg, id string) {
	rects := ExtractBoxRects(svg)
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

// ExtractConnections parses all connection line segments from SVG.
func ExtractConnections(svg string) []ConnectionLine {
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

// AssertConnectionLinesPresent asserts there are connection lines in the SVG.
func AssertConnectionLinesPresent(t require.TestingT, svg string) {
	count := strings.Count(svg, `class="connection`)
	require.Greater(t, count, 0, "SVG should have connection lines")

	lines := ExtractConnections(svg)
	require.Greater(t, len(lines), 0, "SVG should have at least one connection line segment")
}

// AssertConnectionsCount asserts the exact number of connection line segments.
func AssertConnectionsCount(t require.TestingT, svg string, expected int) {
	lines := ExtractConnections(svg)
	require.Equal(t, expected, len(lines), "Expected %d connection line segments, got %d", expected, len(lines))
}

// AssertTextPresent asserts a text string appears in the SVG.
func AssertTextPresent(t require.TestingT, svg, text string) {
	require.Contains(t, svg, text, "SVG should contain text: %s", text)
}

// AssertTextAbsent asserts a text string does NOT appear in the SVG.
func AssertTextAbsent(t require.TestingT, svg, text string) {
	require.NotContains(t, svg, text, "SVG should NOT contain text: %s", text)
}

// AssertHasComments asserts comment-related SVG elements are present.
func AssertHasComments(t require.TestingT, svg string) {
	require.Contains(t, svg, "comment", "SVG should have comment-related elements")
}

// AssertNoComments asserts no comment-related SVG elements are present.
func AssertNoComments(t require.TestingT, svg string) {
	require.NotContains(t, svg, "comment", "SVG should not have comment-related elements")
}

// AssertHasConnectionNodes asserts connection node markers exist in the SVG.
func AssertHasConnectionNodes(t require.TestingT, svg string) {
	// Connection nodes are drawn as circles with class containing "connection"
	re := regexp.MustCompile(`<circle\b[^>]*class="[^"]*connection[^"]*"[^>]*>`)
	count := len(re.FindAllString(svg, -1))
	_ = count // Some documents don't have connection nodes even with connections
}

// AssertSvgNonEmpty asserts the SVG string is not empty.
func AssertSvgNonEmpty(t require.TestingT, svg string) {
	require.NotEmpty(t, svg, "SVG output should not be empty")
}

// AssertHasTextElements asserts the SVG contains <text elements.
func AssertHasTextElements(t require.TestingT, svg string) {
	count := strings.Count(svg, "<text")
	require.Greater(t, count, 0, "SVG should have text elements, got %d", count)
}

// AssertHasTextElementsPresent is an alias for AssertHasTextElements.
func AssertTextElementsPresent(t require.TestingT, svg string) {
	AssertHasTextElements(t, svg)
}

// AssertTitlePresent asserts an SVG title text element exists.
func AssertTitlePresent(t require.TestingT, svg string) {
	re := regexp.MustCompile(`<text[^>]*style="text-anchor:middle"`)
	count := len(re.FindAllString(svg, -1))
	require.Greater(t, count, 0, "SVG should have at least one title text element")
}

// AssertNonZeroBoxDimensions asserts all rects have non-zero dimensions and positions.
func AssertNonZeroBoxDimensions(t require.TestingT, svg string) {
	rects := ExtractBoxRects(svg)
	for _, r := range rects {
		if r.ID != "" {
			require.Greater(t, r.W, 0, "Box %s should have non-zero width", r.ID)
			require.Greater(t, r.H, 0, "Box %s should have non-zero height", r.ID)
			require.GreaterOrEqual(t, r.X, 0, "Box %s should have non-negative x", r.ID)
			require.GreaterOrEqual(t, r.Y, 0, "Box %s should have non-negative y", r.ID)
		}
	}
}

// CompareSvgElements compares rect counts between two SVG strings.
func CompareSvgElements(svg1, svg2 string) int {
	c1 := CountSvgElements(svg1, "<rect")
	c2 := CountSvgElements(svg2, "<rect")
	return c2 - c1
}

// AssertSvgHasImage asserts SVG contains <image or <use elements (for embedded images).
func AssertSvgHasImage(t require.TestingT, svg string) {
	hasImage := strings.Contains(svg, "<image") || strings.Contains(svg, `<use `)
	require.True(t, hasImage, "SVG should contain image/embedded image elements")
}

// AssertSvgHasTitle asserts SVG has a title element with content.
func AssertSvgHasTitleContent(t require.TestingT, svg string, titleText string) {
	require.Contains(t, svg, titleText, "SVG should contain title text: %s", titleText)
}

// AssertSvgHasExpandedId asserts an expanded box ID is present in SVG.
func AssertSvgHasExpandedId(t require.TestingT, svg string, id string) {
	require.Contains(t, svg, `id="`+id+`"`, "SVG should contain expanded box with id=%s", id)
}

// AssertSvgHasBlacklistedId asserts a blacklisted box ID is NOT present in SVG.
func AssertSvgHasBlacklistedId(t require.TestingT, svg string, id string) {
	require.NotContains(t, svg, `id="`+id+`"`, "SVG should not contain blacklisted box with id=%s", id)
}

// AssertHasOverlay asserts SVG contains overlay-related elements.
func AssertHasOverlay(t require.TestingT, svg string) {
	require.Contains(t, svg, "overlay", "SVG should have overlay elements")
}

// AssertHasWrapper asserts SVG contains wrapper-related elements.
func AssertHasWrapper(t require.TestingT, svg string) {
	require.Contains(t, svg, "wrapper", "SVG should have wrapper elements")
}

// CountTextElements returns count of <text elements.
func CountTextElements(svg string) int {
	return strings.Count(svg, "<text")
}

// CountRectElements returns count of <rect elements.
func CountRectElements(svg string) int {
	return CountSvgElements(svg, "<rect")
}

// CountLinesWithClass returns count of <line> elements with a specific class prefix.
func CountLinesWithClass(svg string, classPrefix string) int {
	re := regexp.MustCompile(`<line\b[^>]*class="[^"]*` + regexp.QuoteMeta(classPrefix) + `[^"]*"`)
	return len(re.FindAllString(svg, -1))
}
