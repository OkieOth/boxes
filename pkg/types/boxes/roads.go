package boxes

import (
	"sort"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

// traces the document and finds all reasonable "roads" to connect
func (doc *BoxesDocument) InitRoads() {
	doc.initRoadsImpl(&doc.Boxes, DefRoadType_All, true, true)
	// sort roads
	sort.Slice(doc.VerticalRoads, func(i, j int) bool {
		return doc.VerticalRoads[i].StartX < doc.VerticalRoads[j].StartX
	})
	sort.Slice(doc.HorizontalRoads, func(i, j int) bool {
		return doc.HorizontalRoads[i].StartY < doc.HorizontalRoads[j].StartY
	})
}

func (doc *BoxesDocument) addVerticalRoad(line ConnectionLine) {
	l := line
	if line.StartX > line.EndX || line.StartY > line.EndY {
		l = newConnectionLine(line.EndX, line.EndY, line.StartX, line.StartY)
	}
	for i := range len(doc.VerticalRoads) {
		cl := &doc.VerticalRoads[i]
		if (cl.StartX == l.StartX) &&
			(cl.EndX == l.EndX) &&
			(cl.StartY == l.StartY) &&
			(cl.EndY == l.EndY) {
			// already exists
			return
		}
		if cl.StartX == l.StartX {
			if (l.StartY < cl.StartY) && (l.EndY >= cl.StartY) {
				// extend the existing line to the top
				cl.StartY = l.StartY
				return
			} else if (l.EndY > cl.EndY) && (l.StartY <= cl.EndY) {
				// extend the existing line to the bottom
				cl.EndY = l.EndY
				return
			}
		}
	}
	doc.VerticalRoads = append(doc.VerticalRoads, l)
}

func (doc *BoxesDocument) addHorizontalRoad(line ConnectionLine) {
	l := line
	if line.StartX > line.EndX || line.StartY > line.EndY {
		l = newConnectionLine(line.EndX, line.EndY, line.StartX, line.StartY)
	}
	for i := range len(doc.HorizontalRoads) {
		cl := &doc.HorizontalRoads[i]
		if (cl.StartX == l.StartX) &&
			(cl.EndX == l.EndX) &&
			(cl.StartY == l.StartY) &&
			(cl.EndY == l.EndY) {
			// already exists
			return
		}
		if cl.StartY == l.StartY {
			if (l.StartX < cl.StartX) && (l.EndX >= cl.StartX) {
				// extend the existing line to the left
				cl.StartX = l.StartX
				return
			} else if (l.EndX > cl.EndX) && (l.StartX <= cl.EndX) {
				// extend the existing line to the right
				cl.EndX = l.EndX
				return
			}
		}
	}
	doc.HorizontalRoads = append(doc.HorizontalRoads, l)
}

func (doc *BoxesDocument) pointHasCollision(x, y int, startElem *LayoutElement, isForHorizontalLine bool) bool {
	switch doc.checkColl(x, y, &doc.Boxes, startElem, isForHorizontalLine) {
	case CollisionType_WithElem:
		return true
	case CollisionType_WithSurroundings:
		return true
	default:
		return false
	}
}

func (doc *BoxesDocument) roadUp(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY <= 0 {
		line.EndY = 0
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, false) {
	case CollisionType_WithElem:
		line.EndY += types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndY -= types.RasterSize
		doc.roadUp(line, startElem)
		return
	default:
		line.EndY -= types.RasterSize
		doc.roadUp(line, startElem)
	}
}

func (doc *BoxesDocument) roadDown(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY >= doc.Height {
		line.EndY = doc.Height
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, false) {
	case CollisionType_WithElem:
		line.EndY -= types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndY += types.RasterSize
		doc.roadDown(line, startElem)
		return
	default:
		line.EndY += types.RasterSize
		doc.roadDown(line, startElem)
	}
}

func (doc *BoxesDocument) roadLeft(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX <= 0 {
		line.EndX = 0
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, true) {
	case CollisionType_WithElem:
		line.EndX += types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndX -= types.RasterSize // bug???
		return
	default:
		line.EndX -= types.RasterSize
		doc.roadLeft(line, startElem)
	}
}

func (doc *BoxesDocument) roadRight(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX >= doc.Width {
		line.EndX = doc.Width
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, true) {
	case CollisionType_WithElem:
		line.EndX -= types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndX += types.RasterSize
		return
	default:
		line.EndX += types.RasterSize
		doc.roadRight(line, startElem)
	}
}

func (doc *BoxesDocument) elemHasParentWithTextImpl(elemToCheck, currentElem *LayoutElement, parentTxt bool) bool {
	if currentElem.Id != "" && currentElem.Id == elemToCheck.Id {
		// reached end of recursion
		return parentTxt
	}

	var textOverlaps bool
	if currentElem.XTextBox != nil && currentElem.WidthTextBox != nil {
		textOverlaps = elemToCheck.CenterX >= *currentElem.XTextBox && elemToCheck.CenterX <= (*currentElem.XTextBox+*currentElem.WidthTextBox)
	}

	myParentTxt := textOverlaps || parentTxt
	if doc.elemHasParentWithTextCont(elemToCheck, myParentTxt, currentElem.Horizontal) {
		return true
	}
	return doc.elemHasParentWithTextCont(elemToCheck, myParentTxt, currentElem.Vertical)
}

func (doc *BoxesDocument) elemHasParentWithTextCont(elem *LayoutElement, parentTxt bool, cont *LayoutElemContainer) bool {
	if cont == nil {
		return false
	}
	for _, e := range cont.Elems {
		if doc.elemHasParentWithTextImpl(elem, &e, parentTxt) {
			return true
		}
	}
	return false
}

func (doc *BoxesDocument) elemHasParentWithText(elem *LayoutElement) bool {
	return doc.elemHasParentWithTextImpl(elem, &doc.Boxes, false)
}

type DefRoadType int

const (
	DefRoadType_All DefRoadType = iota
	DefRoadType_Vertical
	DefRoadType_Horizontal
	DefRoadType_None
)

func (doc *BoxesDocument) initRoadsImpl(elem *LayoutElement, defRoadType DefRoadType, needsRight, needsBottom bool) {
	if elem.Caption == "" && elem.Text1 == "" && elem.Text2 == "" {
		defRoadType = DefRoadType_None
	}
	stepSize := 2 * types.RasterSize
	if elem.TopXToStart != nil { // elem has connections by itself
		// draw line from the top x start, till the first collision
		if !doc.pointHasCollision(*elem.TopXToStart, elem.Y+stepSize, elem, false) {
			upRoad := newConnectionLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y+stepSize)
			doc.roadUp(&upRoad, elem)
			doc.addVerticalRoad(upRoad)
		}
		// draw line from the bottom x start, till the first collision
		if !doc.pointHasCollision(*elem.BottomXToStart, elem.Y+elem.Height+stepSize, elem, false) {
			downRoad := newConnectionLine(*elem.BottomXToStart, elem.Y+elem.Height,
				*elem.BottomXToStart, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			doc.addVerticalRoad(downRoad)
		}

		// draw line from the left y start, till the first collision
		if !doc.pointHasCollision(elem.X-stepSize, *elem.LeftYToStart, elem, true) {
			leftRoad := newConnectionLine(elem.X, *elem.LeftYToStart, elem.X-stepSize, *elem.LeftYToStart)
			doc.roadLeft(&leftRoad, elem)
			doc.addHorizontalRoad(leftRoad)
		}
		// draw line from the right y start, till the first collision
		if !doc.pointHasCollision(elem.X+elem.Width+stepSize, *elem.RightYToStart, elem, true) {
			rightRoad := newConnectionLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+stepSize, *elem.RightYToStart)
			doc.roadRight(&rightRoad, elem)
			doc.addHorizontalRoad(rightRoad)
		}
		defRoadType = DefRoadType_All
	}
	if needsRight { // DEBUG, TODO: defRoadType == DefRoadType_All {
		// draw the line parallel to the right border, till the first collision, up and down
		if !doc.pointHasCollision(elem.X+elem.Width+stepSize, elem.CenterY, elem, false) {
			upRoad := newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y-stepSize)
			doc.roadUp(&upRoad, elem)
			downRoad := newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			l := newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
			doc.addVerticalRoad(l)
		}
	}
	if defRoadType == DefRoadType_Horizontal || defRoadType == DefRoadType_All {
		// draw the line parallel to the left border, till the first collision, up and down
		if !doc.pointHasCollision(elem.X-stepSize, elem.CenterY, elem, false) {
			upRoad := newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y-stepSize)
			doc.roadUp(&upRoad, elem)
			downRoad := newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			l := newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
			doc.addVerticalRoad(l)
		}
	}

	if defRoadType == DefRoadType_All {
		// draw the line parallel to the top border, till the first collision, left and right
		if !doc.pointHasCollision(elem.CenterX, elem.Y-stepSize, elem, true) {
			leftRoad := newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X-stepSize, elem.Y-stepSize)
			doc.roadLeft(&leftRoad, elem)
			rightRoad := newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X+elem.Width+stepSize, elem.Y-stepSize)
			doc.roadRight(&rightRoad, elem)
			l := newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
			doc.addHorizontalRoad(l)
		}
	}
	if needsBottom { // DEBUG, TODO: if defRoadType == DefRoadType_Vertical || defRoadType == DefRoadType_All {
		// draw the line parallel to the bottom border, till the first collision, left and right
		if !doc.pointHasCollision(elem.CenterX, elem.Y+elem.Height+stepSize, elem, true) {
			leftRoad := newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X-stepSize, elem.Y+elem.Height+stepSize)
			doc.roadLeft(&leftRoad, elem)
			rightRoad := newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
			doc.roadRight(&rightRoad, elem)
			l := newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
			doc.addHorizontalRoad(l)
		}
	}
	if elem.Vertical != nil {
		length := len(elem.Vertical.Elems)
		lastIndex := length - 1
		for i := range length {
			defRoadType := DefRoadType_All
			doc.initRoadsImpl(&elem.Vertical.Elems[i], defRoadType, true, i == lastIndex)
		}
	}
	if elem.Horizontal != nil {
		length := len(elem.Horizontal.Elems)
		lastIndex := length - 1
		for i := range len(elem.Horizontal.Elems) {
			defRoadType := DefRoadType_All
			doc.initRoadsImpl(&elem.Horizontal.Elems[i], defRoadType, i == lastIndex, true)
		}
	}
}
