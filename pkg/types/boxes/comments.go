package boxes

import (
	"fmt"

	"github.com/okieoth/boxes/pkg/types"
)

func (doc *BoxesDocument) IncludeComments(c types.TextDimensionCalculator) error {
	doc.collectCommentsFromLayout(&doc.Boxes, c)
	doc.collectCommentsFromConnections(c)
	doc.CommentMarkerRadius = (doc.CommentMarkerRadius / 2) + 4
	// adjust the document height, based on the comments
	if len(doc.Comments) > 0 {
		neededMarkerSpace := doc.CommentMarkerRadius * 2
		currentY := doc.Height + types.GlobalPadding
		for i := range doc.Comments {
			if doc.Comments[i].Text == "" {
				continue
			}
			currentY += getMax(doc.Comments[i].TextHeight, neededMarkerSpace)
			currentY += types.GlobalPadding

		}
		doc.Height = currentY + doc.GlobalPadding
	}
	return nil
}

func (doc *BoxesDocument) checkForOverlappingCommentMarkers(x, y int) bool {
	overlapRange := doc.GlobalPadding * 2
	for i := range doc.Comments {
		c := doc.Comments[i]
		diffX := absInt(c.MarkerX - x)
		diffY := absInt(c.MarkerY - y)
		if diffX <= overlapRange && diffY <= overlapRange {
			return true
		}
	}
	return false
}

func (doc *BoxesDocument) collectCommentFromConnectionsImpl(c *ConnectionElem, checkForOverlap bool, dimensionsCalc types.TextDimensionCalculator) bool {
	for li := range doc.VerticalLines {
		// search for the start line of the connection in the vertical lines
		l := doc.VerticalLines[li]
		if l.ConnectionIndex == c.ConnectionIndex {
			x := l.StartX
			diff, changed := absInt2(l.EndY - l.StartY)
			if diff < 20 {
				continue
			}
			var y int
			if changed {
				y = l.EndY - (diff / 2)
			} else {
				y = l.StartY + (diff / 2)
			}
			if !checkForOverlap || !doc.checkForOverlappingCommentMarkers(x, y) {
				if c.HiddenComments {
					cc := doc.createHiddenComment(x, y, dimensionsCalc)
					doc.Comments = append(doc.Comments, cc)
				} else {
					label, customMarker := doc.newLabel(c.Comment.Label)
					cc := doc.newCommentContainer(c.Comment.Text, label, c.Comment.Format, x, y, false, dimensionsCalc, customMarker, &l.ConnectionIndex)
					cc.Step = c.Step
					doc.Comments = append(doc.Comments, cc)
				}
				return true
			}
		}
	}
	for li := range doc.HorizontalLines {
		// search for the start line of the connection in the horizontal lines
		l := doc.HorizontalLines[li]
		if l.ConnectionIndex == c.ConnectionIndex {
			diff, changed := absInt2(l.EndX - l.StartX)
			if diff < 20 {
				continue
			}
			var x int
			if changed {
				x = l.EndX - (diff / 2)
			} else {
				x = l.StartX + (diff / 2)
			}
			y := l.StartY
			if !checkForOverlap || !doc.checkForOverlappingCommentMarkers(x, y) {
				if c.HiddenComments {
					cc := doc.createHiddenComment(x, y, dimensionsCalc)
					doc.Comments = append(doc.Comments, cc)
				} else {
					label, customMarker := doc.newLabel(c.Comment.Label)
					cc := doc.newCommentContainer(c.Comment.Text, label, c.Comment.Format, x, y, false, dimensionsCalc, customMarker, &l.ConnectionIndex)
					cc.Step = c.Step
					doc.Comments = append(doc.Comments, cc)
				}
				return true
			}
		}
	}
	return false
}

func (doc *BoxesDocument) collectCommentsFromConnections(dimensionsCalc types.TextDimensionCalculator) {
	for i := range doc.Connections {
		c := doc.Connections[i]
		if c.Comment != nil {
			if !doc.collectCommentFromConnectionsImpl(&c, true, dimensionsCalc) {
				// in case no marker could set without overlap ... then it's set with overlap
				doc.collectCommentFromConnectionsImpl(&c, false, dimensionsCalc)
			}
		}
	}
}

func (doc *BoxesDocument) GetCommentFormat(format *string) CommentFormat {
	var boxFormat *BoxFormat
	if format != nil {
		if f, ok := doc.Formats[*format]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		if f, ok := doc.Formats["defaultComment"]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		if f, ok := doc.Formats["default"]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		return CommentFormat{
			FontText:   types.InitFontDef(nil, "serif", 8, false, false, 10),
			FontMarker: types.InitFontDef(nil, "monospace", 8, false, true, 10),
			Line:       *types.InitLineDef(nil),
			Fill:       *types.InitFillDef(nil),
		}
	} else {
		if boxFormat.Fill == nil {
			boxFormat.Fill = types.InitFillDef(nil)
		}
		if boxFormat.Line == nil {
			boxFormat.Line = types.InitLineDef(nil)
		}
		return CommentFormat{
			FontText:   boxFormat.FontComment,
			FontMarker: boxFormat.FontCommentMarker,
			Line:       *boxFormat.Line,
			Fill:       *boxFormat.Fill,
		}

	}
}

func (doc *BoxesDocument) newCommentContainer(
	text, label string,
	format *string,
	x, y int,
	moveMarker bool,
	dimensionsCalc types.TextDimensionCalculator,
	customMarker bool,
	connectionIndex *int) CommentContainer {
	if moveMarker {
		diff := doc.GlobalPadding + 2
		x, y = x+diff, y+diff
	}

	f := doc.GetCommentFormat(format)
	tw, th := dimensionsCalc.DimensionsWithMaxWidth(text, &f.FontText, doc.Width)
	mw, mh := dimensionsCalc.DimensionsWithMaxWidth(label, &f.FontMarker, doc.Width)

	if mw > doc.CommentMarkerRadius {
		doc.CommentMarkerRadius = mw
	}
	if mh > doc.CommentMarkerRadius {
		doc.CommentMarkerRadius = mh
	}

	return CommentContainer{
		Text:             text,
		Label:            label,
		Format:           f,
		MarkerX:          x,
		MarkerY:          y,
		TextWidth:        tw,
		TextHeight:       th,
		MarkerTextWidth:  mw,
		MarkerTextHeight: mh,
		CustomMarker:     customMarker,
		ConnectionIndex:  connectionIndex,
	}
}

func (doc *BoxesDocument) newLabel(label *string) (string, bool) {
	if label != nil {
		return *label, true
	} else {
		doc.CommentCurrent += 1
		return fmt.Sprintf("%d", doc.CommentCurrent), false
	}
}

func (doc *BoxesDocument) collectCommentsFromLayout(l *LayoutElement, dimensionsCalc types.TextDimensionCalculator) {
	currentX := l.X + l.Width - 4
	for i := range l.Comments {
		comment := l.Comments[i]
		label, customMarker := doc.newLabel(comment.Label)
		c := doc.newCommentContainer(comment.Text, label, comment.Format, currentX, l.Y+4, false, dimensionsCalc, customMarker, nil)
		c.Step = comment.Step
		doc.Comments = append(doc.Comments, c)
		currentX -= (2 * c.MarkerTextWidth) + doc.GlobalPadding + 2
		if i > 7 {
			break
		}
	}
	if l.HiddenComments {
		c := doc.createHiddenComment(l.X, l.Y, dimensionsCalc)
		doc.Comments = append(doc.Comments, c)
	}
	doc.collectCommentsFromLayoutCont(l.Horizontal, dimensionsCalc)
	doc.collectCommentsFromLayoutCont(l.Vertical, dimensionsCalc)
}

func (doc *BoxesDocument) createHiddenComment(x, y int, dimensionsCalc types.TextDimensionCalculator) CommentContainer {
	var text string
	if !doc.HasHiddenComments {
		text = "Has comments in hidden childs"
		doc.HasHiddenComments = true
	}
	return doc.newCommentContainer(text, "...", nil, x, y, false, dimensionsCalc, true, nil)
}

func (doc *BoxesDocument) collectCommentsFromLayoutCont(cont *LayoutElemContainer, dimensionsCalc types.TextDimensionCalculator) {
	if cont != nil {
		for i := range cont.Elems {
			doc.collectCommentsFromLayout(&cont.Elems[i], dimensionsCalc)
		}
	}
}
