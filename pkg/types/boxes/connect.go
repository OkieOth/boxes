package boxes

import (
	"fmt"
	"slices"

	"github.com/okieoth/boxes/pkg/types"
)

type PointToTest struct {
	X            int
	Y            int
	HasCollision bool
}

type checkForCollFunc2 func(x, y int, currentElem, startElem *LayoutElement, isForHorizontalLine bool) CollisionType

func (doc *BoxesDocument) checkForCollInContainer2(cont *LayoutElemContainer, f checkForCollFunc2, x, y int, startElem *LayoutElement, isForHorizontalLine bool) CollisionType {
	if cont != nil {
		for i := range cont.Elems {
			e := &cont.Elems[i]
			if e == startElem {
				continue
			}
			if ct := f(x, y, e, startElem, isForHorizontalLine); ct != CollisionType_NoCollision {
				return ct
			}
		}
	}
	return CollisionType_NoCollision
}

type CollisionType int

const (
	CollisionType_NoCollision CollisionType = iota
	CollisionType_WithElem
	CollisionType_WithSurroundings
)

func (doc *BoxesDocument) isParentInContainer(container *LayoutElemContainer, possibleParent, elemToCheckFor *LayoutElement) bool {
	if container != nil {
		for _, subElem := range container.Elems {
			if subElem.Id == elemToCheckFor.Id {
				return true
			}
			if doc.isParent(&subElem, elemToCheckFor) {
				return true
			}
		}
	}
	return false
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func (doc *BoxesDocument) isParent(possibleParent, elemToCheckFor *LayoutElement) bool {
	if doc.isParentInContainer(possibleParent.Vertical, possibleParent, elemToCheckFor) {
		return true
	}
	return doc.isParentInContainer(possibleParent.Horizontal, possibleParent, elemToCheckFor)
}

func (doc *BoxesDocument) checkColl(x, y int, currentElem, startElem *LayoutElement, isForHorizontalLine bool) CollisionType {
	if (currentElem != startElem) && (doc.ShouldHandle(currentElem)) {
		currentElemIsParentToStart := false
		if startElem != nil {
			currentElemIsParentToStart = doc.isParent(currentElem, startElem)
		}
		if !currentElemIsParentToStart {
			// end - not needed any longer because only roads need that this function
			if currentElem.XTextBox != nil && currentElem.DontBlockConPaths != nil && *currentElem.DontBlockConPaths {
				curMinX := *currentElem.XTextBox
				curMaxX := *currentElem.XTextBox + *currentElem.WidthTextBox
				curMinY := *currentElem.YTextBox
				curMaxY := *currentElem.YTextBox + *currentElem.HeightTextBox
				if (x <= curMaxX) && (x >= curMinX) &&
					(y >= curMinY) && (y <= curMaxY) {
					return CollisionType_WithElem
				}

			} else {
				curMinX := currentElem.X
				curMaxX := currentElem.X + currentElem.Width
				curMinY := currentElem.Y
				curMaxY := currentElem.Y + currentElem.Height
				if (x < curMaxX) && (x > curMinX) &&
					(y > curMinY) && (y < curMaxY) {
					return CollisionType_WithElem
				}
				if (x <= (curMaxX + types.RasterSize)) && (x >= (curMinX - types.RasterSize)) &&
					(y >= (curMinY - types.RasterSize)) && (y <= (curMaxY + types.RasterSize)) {
					return CollisionType_WithSurroundings
				}
			}
		}
		// not needed any longer because only roads need that this function
		if currentElemIsParentToStart {
			// check if there is a collision with the text
			if currentElem.WidthTextBox != nil {
				if (x <= (*currentElem.XTextBox + *currentElem.WidthTextBox + 5)) && (x >= (*currentElem.XTextBox - 5)) &&
					(y >= *currentElem.YTextBox) && (y <= (*currentElem.YTextBox + *currentElem.HeightTextBox)) {
					return CollisionType_WithElem
				}
				if isForHorizontalLine {
					curMinY := currentElem.Y
					curMaxY := currentElem.Y + currentElem.Height
					if (y == curMinY) || (y == curMaxY) {
						return CollisionType_WithElem
					}
				} else {
					curMinX := currentElem.X
					curMaxX := currentElem.X + currentElem.Width
					if (x == curMaxX) || (x == curMinX) {
						return CollisionType_WithElem
					}
				}
			}
		}
		// end - not needed any longer because only roads need that this function
	}
	if ct := doc.checkForCollInContainer2(currentElem.Vertical, doc.checkColl, x, y, startElem, isForHorizontalLine); ct != CollisionType_NoCollision {
		return ct
	}
	return doc.checkForCollInContainer2(currentElem.Horizontal, doc.checkColl, x, y, startElem, isForHorizontalLine)
}

func newConnectionLine(x1, y1, x2, y2 int) ConnectionLine {
	return ConnectionLine{
		StartX: x1,
		StartY: y1,
		EndX:   x2,
		EndY:   y2,
	}
}

func (doc *BoxesDocument) connectContImpl(layoutCont []LayoutElement) {
	for i := range layoutCont {
		doc.connectImpl(&layoutCont[i])
	}
}

// The function guarantees that the lines always to top/down or left to right.
// That simplifies the code for separating overlapping lines
func (doc *BoxesDocument) createConnection(x1, y1, x2, y2 int) ConnectionLine {
	inverseDirection := false
	if x1 == x2 {
		// vertical line
		if y1 > y2 {
			// bottom/up line
			inverseDirection = true
			y1, y2 = y2, y1
		}
	} else {
		if x1 > x2 {
			// right to left
			inverseDirection = true
			x1, x2 = x2, x1
		}
	}
	line := ConnectionLine{
		StartX:           x1,
		StartY:           y1,
		EndX:             x2,
		EndY:             y2,
		InverseDirection: inverseDirection,
	}
	return line
}

func (doc *BoxesDocument) extendConnection(line ConnectionLine, x, y int) ConnectionLine {
	if line.StartX == line.EndX {
		// vertical line
		if line.StartY > y {
			line.StartY = y
		} else {
			line.EndY = y
		}
	} else {
		if line.StartX > x {
			// right to left
			line.StartX = x
		} else {
			line.EndX = x
		}
	}
	return line
}

// aggregates line parts that have the same direction
func (doc *BoxesDocument) reduceConnectionLines(connElem *ConnectionElem) {
	reducedParts := make([]ConnectionLine, 0)
	var lastE *ConnectionLine
	for _, e := range connElem.Parts {
		if lastE == nil {
			lastE = &e
			lastE.IsStart = true
		} else {
			if lastE.StartX == lastE.EndX {
				// lastE is vertical line
				if e.StartX == e.EndX {
					// ... still vertical line
					// expects that e.StartX == lastE.StartX
					if e.StartY <= lastE.StartY {
						lastE.StartY = e.StartY
					} else {
						lastE.EndY = e.EndY
					}
				} else {
					// change from vertical to horizontal line
					lastE.LineIndex = len(reducedParts)
					reducedParts = append(reducedParts, *lastE)
					lastE = &e
				}
			} else {
				// lastE is horizontal line
				if e.StartY == e.EndY {
					// ... still horizontal line
					// expects that e.StartY == lastE.StartY
					if e.StartX <= lastE.StartX {
						lastE.StartX = e.StartX
					} else {
						lastE.EndX = e.EndX
					}
				} else {
					// change to vertical line
					lastE.LineIndex = len(reducedParts)
					reducedParts = append(reducedParts, *lastE)
					lastE = &e
				}
			}
		}
	}
	lastPart := connElem.Parts[len(connElem.Parts)-1]
	if lastPart.DestLayoutId != nil {
		lastE.DestLayoutId = lastPart.DestLayoutId
	}
	if lastPart.SrcLayoutId != nil {
		lastE.SrcLayoutId = lastPart.SrcLayoutId
	}
	lastE.LineIndex = len(reducedParts)
	lastE.IsEnd = true
	reducedParts = append(reducedParts, *lastE)
	connElem.Parts = reducedParts
}

func (doc *BoxesDocument) createAConnectionPath(path []ConnectionNode, c *LayoutElemConnection, srcId, destId string) {
	if len(path) < 2 {
		return
	}
	pathLen := len(path)
	pathToDraw := path[1 : pathLen-1]
	connElem := NewConnectionElem()
	connElem.From = path[0].NodeId
	connElem.To = path[pathLen-1].NodeId
	connElem.HiddenComments = c.HiddenComments
	connElem.ConnectionIndex = len(doc.Connections)
	if c.Format != nil {
		connElem.Format = c.Format
	}
	connElem.Comment = c.Comment

	var lastX, lastY int
	pathElemCount := len(pathToDraw) - 1
	for i, p := range pathToDraw {
		if i > 0 {
			var line ConnectionLine
			line = doc.createConnection(lastX, lastY, p.X, p.Y)
			if i == 1 {
				line.SrcLayoutId = &srcId
			}
			if i == pathElemCount {
				line.DestLayoutId = &destId
			}
			line.ConnectionIndex = connElem.ConnectionIndex
			connElem.Parts = append(connElem.Parts, line)
		}
		lastX = p.X
		lastY = p.Y
	}
	// aggregates line parts that have the same direction
	doc.reduceConnectionLines(connElem)
	doc.Connections = append(doc.Connections, *connElem)
}

func (doc *BoxesDocument) findBoxInContWithId(cont *LayoutElemContainer, id string) *LayoutElement {
	if cont == nil {
		return nil
	}
	for i := range cont.Elems {
		found := doc.findBoxWithId(&cont.Elems[i], id)
		if found != nil {
			return found
		}
	}
	return nil
}

func (doc *BoxesDocument) findBoxWithId(box *LayoutElement, id string) *LayoutElement {
	if box.Id == id {
		return box
	}
	found := doc.findBoxInContWithId(box.Vertical, id)
	if found != nil {
		return found
	}
	return doc.findBoxInContWithId(box.Horizontal, id)
}

func (doc *BoxesDocument) FindBoxWithId(id string) *LayoutElement {
	return doc.findBoxWithId(&doc.Boxes, id)
}

func (doc *BoxesDocument) connectImpl(layout *LayoutElement) {
	if len(layout.Connections) > 0 {
		for i := range layout.Connections {
			c := layout.Connections[i]
			srcId, destId := layout.Id, c.DestId
			if path, _, ok := doc.DijkstraPath(srcId, destId); ok {
				doc.createAConnectionPath(path, &c, srcId, destId)
			} else {
				fmt.Printf("Couldn't calculate path: src=%s, dest=>%s<\n", srcId, destId)
			}
		}
	}
	if layout.Horizontal != nil {
		doc.connectContImpl(layout.Horizontal.Elems)
	}
	if layout.Vertical != nil {
		doc.connectContImpl(layout.Vertical.Elems)
	}
}

func (doc *BoxesDocument) separateConnectionLines() {
	doc.HorizontalLines = make([]ConnectionLine, 0)
	doc.VerticalLines = make([]ConnectionLine, 0)
	for _, c := range doc.Connections {
		for _, p := range c.Parts {
			p.Format = c.Format
			if p.StartX == p.EndX {
				// vertical
				doc.VerticalLines = append(doc.VerticalLines, p)
			} else {
				// assumed horizontal
				doc.HorizontalLines = append(doc.HorizontalLines, p)
			}
		}
	}
	slices.SortFunc(doc.HorizontalLines, func(l1, l2 ConnectionLine) int {
		yDiff := l1.StartY - l2.StartY
		if yDiff != 0 {
			return yDiff
		} else {
			return l1.StartX - l2.StartX
		}
	})
	slices.SortFunc(doc.VerticalLines, func(l1, l2 ConnectionLine) int {
		xDiff := l1.StartX - l2.StartX
		if xDiff != 0 {
			return xDiff
		} else {
			return l1.StartY - l2.StartY
		}
	})
}

func (doc *BoxesDocument) ConnectBoxes() {
	doc.InitStartPositions()
	doc.InitRoads()
	doc.Roads2ConnectionNodes()
	doc.connectImpl(&doc.Boxes)
	doc.separateConnectionLines()
	doc.adjustForOverlappingConnections()
}

func (doc *BoxesDocument) horizontalRoads2ConnectionNodes() {
	for _, h := range doc.HorizontalRoads {
		// horizontal line
		newH := true
		minHX, maxHX := minMax(h.StartX, h.EndX)
		for j, v := range doc.VerticalRoads {
			// verticalLine
			minVY, maxVY := minMax(v.StartY, v.EndY)
			if (minVY <= h.StartY) && (maxVY >= h.StartY) && // the vertical road covers the y-range of the horizontal line
				(minHX <= v.StartX) && (v.StartX <= maxHX) {
				newNode := CreateConnectionNode(v.StartX, h.StartY)
				newNodeIndex := len(doc.ConnectionNodes)
				id := fmt.Sprintf("__n_%d", newNodeIndex)
				newNode.NodeId = &id
				if j > 0 && !newH {
					weight := v.StartX - doc.ConnectionNodes[len(doc.ConnectionNodes)-1].X
					newEdge := CreateConnectionEdge(v.StartX, h.StartY, weight)
					newEdge.DestNodeId = &id
					doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges = append(doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges, newEdge)
					newEdge2 := CreateConnectionEdge(doc.ConnectionNodes[len(doc.ConnectionNodes)-1].X, doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Y, weight)
					newEdge2.DestNodeId = doc.ConnectionNodes[len(doc.ConnectionNodes)-1].NodeId
					newNode.Edges = append(newNode.Edges, newEdge2)
				}
				doc.ConnectionNodes = append(doc.ConnectionNodes, *newNode)
				newH = false
			}
			if v.StartX > maxHX {
				// the other vertical lines are more right then the current horizontal line
				break
			}
		}
	}
}

func (doc *BoxesDocument) verticalEdges() {
	// adding vertical edges
	for _, v := range doc.VerticalRoads {
		// vertical line
		newV := true
		minVY, maxVY := minMax(v.StartY, v.EndY)
		var lastY int
		for j, h := range doc.HorizontalRoads {
			// verticalLine
			minHX, maxHX := minMax(h.StartX, h.EndX)
			if (minHX <= v.StartX) && (maxHX >= v.StartX) && // the vertical road covers the y-range of the horizontal line
				(minVY <= h.StartY) && (h.StartY <= maxVY) {
				nodeX, nodeY := v.StartX, h.StartY
				if j > 0 && !newV {
					nodeIndex := doc.getMatchingNodeIndexes(nodeX, nodeY)
					lastNodeIndex := doc.getMatchingNodeIndexes(nodeX, lastY)
					if nodeIndex > -1 && lastNodeIndex > -1 {
						weight := nodeY - lastY
						newEdgeUp := CreateConnectionEdge(v.StartX, lastY, weight)
						newEdgeUp.DestNodeId = doc.ConnectionNodes[lastNodeIndex].NodeId
						doc.ConnectionNodes[nodeIndex].Edges = append(doc.ConnectionNodes[nodeIndex].Edges, newEdgeUp)

						newEdgeDown := CreateConnectionEdge(nodeX, nodeY, weight)
						newEdgeDown.DestNodeId = doc.ConnectionNodes[nodeIndex].NodeId
						doc.ConnectionNodes[lastNodeIndex].Edges = append(doc.ConnectionNodes[lastNodeIndex].Edges, newEdgeDown)
					}
				}
				j++
				lastY = h.StartY
				newV = false
			}
			if h.StartY > maxVY {
				// the other horizontal lines are more below the length of the current vertical line
				break
			}
		}
	}
}

func (doc *BoxesDocument) initEdgesForBoxConnections() {
	for i := range doc.ConnectionNodes {
		if doc.ConnectionNodes[i].BoxId != nil {
			// node on a potential start point of a connected box found
			// ... should currently have only one outbound edge!!!
			nodeIndex := doc.getMatchingNodeIndexes(doc.ConnectionNodes[i].Edges[0].X, doc.ConnectionNodes[i].Edges[0].Y)
			if nodeIndex > -1 {
				weight := doc.ConnectionNodes[i].Edges[0].Weight
				doc.ConnectionNodes[i].Edges[0].DestNodeId = doc.ConnectionNodes[nodeIndex].NodeId
				newEdge := CreateConnectionEdge(doc.ConnectionNodes[i].X, doc.ConnectionNodes[i].Y, weight)
				newEdge.DestNodeId = doc.ConnectionNodes[i].NodeId
				doc.ConnectionNodes[nodeIndex].Edges = append(doc.ConnectionNodes[nodeIndex].Edges, newEdge)
			}
		}
	}
}

func (doc *BoxesDocument) Roads2ConnectionNodes() {
	doc.horizontalRoads2ConnectionNodes()
	doc.verticalEdges()
	doc.initEdgesForBoxConnections()
}

func (doc *BoxesDocument) getMatchingNodeIndexes(x, y int) int {
	for i, n := range doc.ConnectionNodes {
		if n.X == x && n.Y == y {
			return i
		}
	}
	return -1
}

func (doc *BoxesDocument) InitStartPositions() {
	doc.initStartPositionsImpl(&doc.Boxes)
}

func CreateConnectionNode(x, y int) *ConnectionNode {
	node := NewConnectionNode()
	node.X = x
	node.Y = y
	return node
}

func CreateConnectionEdge(x, y, weight int) ConnectionEdge {
	return ConnectionEdge{
		X:      x,
		Y:      y,
		Weight: weight,
	}
}

func CreateConnectionEdge2(x, y int, id string, index int) ConnectionEdge {
	return ConnectionEdge{
		X:          x,
		Y:          y,
		DestNodeId: &id,
	}
}

func (doc *BoxesDocument) newConnectionNodeFromStartPos(boxId string, x, y int, edges []ConnectionEdge) {
	n := CreateConnectionNode(x, y)
	index := len(doc.ConnectionNodes)
	nodeId := fmt.Sprintf("__n_%d", index)
	n.NodeId = &nodeId
	n.BoxId = &boxId
	for _, e := range edges {
		n.Edges = append(n.Edges, e)
	}
	doc.ConnectionNodes = append(doc.ConnectionNodes, *n)
}

func (doc *BoxesDocument) initStartPositionsImpl(elem *LayoutElement) {
	if doc.ShouldHandle(elem) {
		if slices.Contains(doc.ConnectedElems, elem.Id) {
			elem.BottomXToStart = &elem.CenterX
			elem.TopXToStart = &elem.CenterX
			elem.LeftYToStart = &elem.CenterY
			elem.RightYToStart = &elem.CenterY
			weight := 2 * types.RasterSize
			noTop := elem.ConnRestrictions != nil && elem.ConnRestrictions.NoTop
			noBottom := elem.ConnRestrictions != nil && elem.ConnRestrictions.NoBottom
			noLeft := elem.ConnRestrictions != nil && elem.ConnRestrictions.NoLeft
			noRight := elem.ConnRestrictions != nil && elem.ConnRestrictions.NoRight
			// add topX
			if !noTop {
				doc.newConnectionNodeFromStartPos(elem.Id, *elem.TopXToStart, elem.Y,
					[]ConnectionEdge{
						CreateConnectionEdge(*elem.TopXToStart, elem.Y-weight, weight),
					})
			}
			// add bottomX
			if !noBottom {
				doc.newConnectionNodeFromStartPos(elem.Id, *elem.BottomXToStart, elem.Y+elem.Height,
					[]ConnectionEdge{
						CreateConnectionEdge(*elem.BottomXToStart, elem.Y+elem.Height+weight, weight),
					})
			}
			// add leftY
			if !noLeft {
				doc.newConnectionNodeFromStartPos(elem.Id, elem.X, *elem.LeftYToStart,
					[]ConnectionEdge{
						CreateConnectionEdge(elem.X-weight, *elem.LeftYToStart, weight),
					})
			} else {
				fmt.Println("DEBUG - noLeft set", elem.Caption)
			}
			// add rightY
			if !noRight {
				doc.newConnectionNodeFromStartPos(elem.Id, elem.X+elem.Width, *elem.RightYToStart,
					[]ConnectionEdge{
						CreateConnectionEdge(elem.X+elem.Width+weight, *elem.LeftYToStart, weight),
					})
			}
		}
	}
	if elem.Vertical != nil {
		for i := range elem.Vertical.Elems {
			doc.initStartPositionsImpl(&elem.Vertical.Elems[i])
		}
	}
	if elem.Horizontal != nil {
		for i := range elem.Horizontal.Elems {
			doc.initStartPositionsImpl(&elem.Horizontal.Elems[i])
		}
	}
}
