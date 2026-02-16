package boxes

// This nasty function tries to resolve collisions of connections lines.
// The approach to solve the issue is, to move overlapping parts and strech
// the elements and the picture in cases.
func (doc *BoxesDocument) adjustForOverlappingConnections() {
	doc.adjustForOverlappingHorizontalLines()
	doc.adjustForOverlappingVerticalLines()
	doc.adjustStartAndEndLines()
}

func (doc *BoxesDocument) fixVerticalStartAndEndOfHorizontalLine(xStart, xEnd, y, yBegin, connectionIndex int) {
	count := len(doc.VerticalLines)
	startFound := false
	endFound := false
	for i := range count {
		line := &doc.VerticalLines[i]
		if line.ConnectionIndex == connectionIndex {
			if line.StartX == xStart || line.StartX == xEnd {
				// the vertical line cuts the triggering line
				if line.StartY == y {
					// vertical line goes down from the triggering line ... shorten the v-line
					line.StartY += doc.LineDist
					// this was removed for a bug ... could be wrong
					//line.EndY += doc.LineDist
					startFound = true
					if startFound && endFound {
						return
					}
					continue
				}
				if line.EndY == y {
					// vertical line goes down to the triggering line ... increase the v-line
					line.EndY += doc.LineDist
					endFound = true
					if startFound && endFound {
						return
					}
					continue
				}
			}
		}
	}
}

func (doc *BoxesDocument) fixHorizontalStartAndEndOfVerticalLine(offset, yStart, yEnd, x, xBegin, connectionIndex int) {
	count := len(doc.HorizontalLines)
	startFound := false
	endFound := false
	for i := range count {
		line := &doc.HorizontalLines[i]
		if line.ConnectionIndex == connectionIndex {
			if line.StartY == yStart || line.StartY == yEnd {
				// the horizontal line cuts the triggering line
				if line.StartX == x {
					// horizontal line goes right from the triggering line
					line.StartX += offset
					line.EndX += offset
					startFound = true
					if startFound && endFound {
						return
					}
					continue
				}
				if line.EndX == x {
					// horizontal line goes right to the triggering line
					line.EndX += offset
					endFound = true
					if startFound && endFound {
						return
					}
					continue
				}
			}
		}
	}
}

func (doc *BoxesDocument) adjustForOverlappingHorizontalLines() {
	xOffset := 0
	horizontalCount := len(doc.HorizontalLines)
	//endLineAdjustments := make(map[string]EndLineAdjustments, 0)
	for i := range horizontalCount {
		if i == 0 {
			continue
		}
		last := &doc.HorizontalLines[i-1]
		current := &doc.HorizontalLines[i]
		if current.StartY == last.StartY &&
			OverlapsHorizontally(current.StartX, current.EndX, last.StartX, last.EndX) {
			// adjust in case start and end lines ... but it's important before moving the h-line
			doc.fixVerticalStartAndEndOfHorizontalLine(current.StartX, current.EndX, current.StartY, current.StartY, current.ConnectionIndex)
			yStart := current.StartY
			current.StartY += doc.LineDist
			current.EndY = current.StartY
			for j := i + 1; j < horizontalCount; j++ {
				cur2 := &doc.HorizontalLines[j]
				// adjust in case start and end lines ... but it's important before moving the h-line
				doc.fixVerticalStartAndEndOfHorizontalLine(cur2.StartX, cur2.EndX, cur2.StartY, yStart, cur2.ConnectionIndex)
				cur2.StartY += doc.LineDist
				cur2.EndY = cur2.StartY
			}
			doc.StretchAndMoveVertical(yStart, doc.LineDist)
			xOffset += doc.LineDist
		}
	}
	doc.Width += xOffset
}

func (doc *BoxesDocument) adjustForOverlappingVerticalLines() {
	yOffset := 0
	verticalCount := len(doc.VerticalLines)
	//endLineAdjustments := make(map[string]EndLineAdjustments, 0)
	for i := range verticalCount {
		if i == 0 {
			continue
		}
		last := &doc.VerticalLines[i-1]
		current := &doc.VerticalLines[i]
		if current.StartX == last.StartX &&
			OverlapsVertically(current.StartY, current.EndY, last.StartY, last.EndY) {
			doc.fixHorizontalStartAndEndOfVerticalLine(doc.LineDist, current.StartY, current.EndY, current.StartX, current.StartX, current.ConnectionIndex)
			xStart := current.StartX
			current.StartX += doc.LineDist
			current.EndX = current.StartX
			for j := i + 1; j < verticalCount; j++ {
				cur2 := &doc.VerticalLines[j]
				doc.fixHorizontalStartAndEndOfVerticalLine(doc.LineDist, cur2.StartY, cur2.EndY, cur2.StartX, xStart, cur2.ConnectionIndex)
				cur2.StartX += doc.LineDist
				cur2.EndX = cur2.StartX
			}
			// strech all horizontal lines, that are in range and not moved before
			doc.StretchAndMoveHorizontal(xStart, doc.LineDist)
			yOffset += doc.LineDist
		}
	}
	doc.Height += yOffset
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func absInt2(x int) (int, bool) {
	if x < 0 {
		return -x, true
	}
	return x, false
}

// that's the past variant before I separated both cases
// remains only for "documentation" purposes
// TODO remove
func (doc *BoxesDocument) adjustEndLine(line *ConnectionLine) {
	destId := line.DestLayoutId
	if destId == nil {
		destId = line.SrcLayoutId
	}
	if destId == nil {
		return
	}
	box := doc.FindBoxWithId(*destId)
	if box == nil {
		return
	}
	if line.StartY == line.EndY {
		// horizontal line
		diffXStart := absInt(box.X - line.StartX)
		diffXEnd := absInt(box.X - line.EndX)
		if diffXStart > diffXEnd {
			// line form left
			line.EndX = box.X
		} else {
			// line to right
			line.StartX = box.X + box.Width
		}

	} else {
		// vertical line
		// the following line is fragile :-/
		diffYStart := absInt(box.Y - line.EndY)
		diffYEnd := absInt((box.Y + box.Height) - line.EndY)
		if diffYStart < diffYEnd {
			// line down
			line.EndY = box.Y
		} else {
			// line up
			line.StartY = box.Y + box.Height
		}
	}
}

func (doc *BoxesDocument) adjustHorizontalEndLineNow(line *ConnectionLine, boxId string, dest bool) {
	box := doc.FindBoxWithId(boxId)
	if box == nil {
		return
	}

	if line.InverseDirection {
		if dest {
			line.StartX = box.X + box.Width
		} else {
			line.EndX = box.X
		}
	} else {
		if dest {
			line.EndX = box.X
		} else {
			line.StartX = box.X + box.Width
		}
	}
}

func (doc *BoxesDocument) adjustHorizontalEndLine(line *ConnectionLine) {
	if line.DestLayoutId != nil {
		doc.adjustHorizontalEndLineNow(line, *line.DestLayoutId, true)
	}
	if line.SrcLayoutId != nil {
		doc.adjustHorizontalEndLineNow(line, *line.SrcLayoutId, false)
	}
}

func (doc *BoxesDocument) adjustVerticalEndLineNow(line *ConnectionLine, boxId string, dest bool) {
	box := doc.FindBoxWithId(boxId)
	if box == nil {
		return
	}

	if line.InverseDirection {
		if dest {
			line.StartY = box.Y + box.Height
		} else {
			line.EndY = box.Y
		}
	} else {
		if dest {
			line.EndY = box.Y
		} else {
			line.StartY = box.Y + box.Height
		}
	}

	// diffYStart := absInt(box.Y - line.EndY)
	// diffYEnd := absInt((box.Y + box.Height) - line.EndY)
	// if diffYStart < diffYEnd {
	// 	// line down
	// 	line.EndY = box.Y
	// } else {
	// 	// line up
	// 	line.StartY = box.Y + box.Height
	// }
}

func (doc *BoxesDocument) adjustVerticalEndLine(line *ConnectionLine) {
	if line.DestLayoutId != nil {
		doc.adjustVerticalEndLineNow(line, *line.DestLayoutId, true)
	}
	if line.SrcLayoutId != nil {
		doc.adjustVerticalEndLineNow(line, *line.SrcLayoutId, false)
	}
}

func (doc *BoxesDocument) adjustHorizontalEndLines() {
	for i := range len(doc.HorizontalLines) {
		doc.adjustHorizontalEndLine(&doc.HorizontalLines[i])
	}
}

func (doc *BoxesDocument) adjustVerticalEndLines() {
	for i := range len(doc.VerticalLines) {
		doc.adjustVerticalEndLine(&doc.VerticalLines[i])
	}
}

func (doc *BoxesDocument) adjustStartAndEndLines() {
	// Critical - in case of not proper line endings on boxes look here
	doc.adjustHorizontalEndLines()
	doc.adjustVerticalEndLines()
}

type EndLineAdjustments struct {
	lines  []*ConnectionLine
	maxLen int
}

func NewEndLineAdjustments() EndLineAdjustments {
	return EndLineAdjustments{
		lines:  make([]*ConnectionLine, 0),
		maxLen: 0,
	}
}

func OverlapsVertically(yTop1, yBottom1, yTop2, yBottom2 int) bool {
	return ((yTop1 >= yTop2 && yTop1 <= yBottom2) || // yTop1 is in range of v-line 2
		(yBottom1 >= yTop2 && yBottom1 <= yBottom2)) || // ...yBottom1 is in range of v-line2
		((yTop2 >= yTop1 && yTop2 <= yBottom1) ||
			(yBottom2 >= yTop1 && yBottom2 <= yBottom1)) ||
		(yTop1 == yTop2) || (yTop1 == yBottom2) ||
		(yBottom1 == yTop2) || (yBottom1 == yBottom2)
}

func OverlapsHorizontally(xLeft1, xRight1, xLeft2, xRight2 int) bool {
	return ((xLeft1 >= xLeft2 && xLeft1 <= xRight2) || // yTop1 is in range of v-line 2
		(xRight1 >= xLeft2 && xRight1 <= xRight2)) || // ...yBottom1 is in range of v-line2
		((xLeft2 >= xLeft1 && xLeft2 <= xRight1) ||
			(xRight2 >= xLeft1 && xRight2 <= xRight1)) ||
		(xLeft1 == xLeft2) || (xLeft1 == xRight2) ||
		(xRight1 == xLeft2) || (xRight1 == xRight2)
}
