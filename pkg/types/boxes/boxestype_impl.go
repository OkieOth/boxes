package boxes

import (
	"fmt"
	"slices"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (d *BoxesDocument) DrawBoxes(drawingImpl types.Drawing) error {
	return d.Boxes.Draw(drawingImpl)
}

func (doc *BoxesDocument) drawStartPositionsImpl(drawingImpl *types.Drawing, elem *LayoutElement, f *types.LineDef) {
	if elem.TopXToStart != nil && elem.BottomXToStart != nil && elem.LeftYToStart != nil && elem.RightYToStart != nil {
		(*drawingImpl).DrawLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y-2, *f)
		(*drawingImpl).DrawLine(*elem.BottomXToStart, elem.Y+elem.Height, *elem.BottomXToStart, elem.Y+elem.Height+2, *f)
		(*drawingImpl).DrawLine(elem.X, *elem.LeftYToStart, elem.X-2, *elem.LeftYToStart, *f)
		(*drawingImpl).DrawLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+2, *elem.RightYToStart, *f)
	}
	if elem.Vertical != nil {
		for i := range elem.Vertical.Elems {
			doc.drawStartPositionsImpl(drawingImpl, &elem.Vertical.Elems[i], f)
		}
	}
	if elem.Horizontal != nil {
		for i := range elem.Horizontal.Elems {
			doc.drawStartPositionsImpl(drawingImpl, &elem.Horizontal.Elems[i], f)
		}
	}
}

func (d *BoxesDocument) ShouldHandle(elem *LayoutElement) bool {
	if elem == &d.Boxes {
		return false
	}
	if elem.Caption == "" && elem.Text1 == "" && elem.Text2 == "" && elem.Id == "" {
		return false
	}
	return true
}

func (d *BoxesDocument) DrawStartPositions(drawingImpl types.Drawing) {
	w := 2.0
	b := "blue"
	f := types.LineDef{
		Width: &w,
		Color: &b,
	}
	d.drawStartPositionsImpl(&drawingImpl, &d.Boxes, &f)
}

func (d *BoxesDocument) DrawConnectionNodes(drawingImpl types.Drawing) {
	w1 := 4.0
	b1 := "grey"
	w2 := 4.0
	b2 := "purple"
	w3 := 2.0
	b3 := "pink"
	w4 := 0.5
	f1 := types.LineDef{
		Width: &w1,
		Color: &b1,
	}
	f2 := types.LineDef{
		Width: &w2,
		Color: &b2,
	}
	f3 := types.LineDef{
		Width: &w3,
		Color: &b3,
	}
	f4 := types.LineDef{
		Width: &w4,
		Color: &b1,
	}
	_ = f3
	for i := range d.ConnectionNodes {
		n := d.ConnectionNodes[i]
		if n.NodeId != nil && *n.NodeId != "" {
			drawingImpl.DrawLine(n.X, n.Y-2, n.X, n.Y+2, f2)
			drawingImpl.DrawLine(n.X-2, n.Y, n.X+2, n.Y, f2)
		} else {
			drawingImpl.DrawLine(n.X, n.Y-2, n.X, n.Y+2, f1)
			drawingImpl.DrawLine(n.X-2, n.Y, n.X+2, n.Y, f1)
		}
		for _, e := range n.Edges {
			nx, ny, ex, ey := n.X, n.Y, e.X, e.Y
			if n.X == e.X {
				// vertical line
				if n.Y > e.Y {
					// line up
					nx = nx - 2
				} else {
					// line down
					nx = nx + 2
				}
				ex = nx
			} else {
				// horizontal line
				if n.X > e.X {
					// line to left
					ny = ny - 2
				} else {
					// line to right
					ny = ny + 2
				}
				ey = ny
			}

			if e.DestNodeId != nil && e.Weight != 0 {
				drawingImpl.DrawLine(nx, ny, ex, ey, f3)
			} else {
				//fmt.Println("DEBUG: boxId=", *n.BoxId, ", id=", *n.NodeId, "edgeCount=", len(n.Edges), ", destNodeId=", e.DestNodeId)
				drawingImpl.DrawLine(nx, ny, ex, ey, f4)
			}
		}
	}
}

func (d *BoxesDocument) DrawRoads(drawingImpl types.Drawing) {
	w := 1.0
	b := "lightgray"
	f := types.LineDef{
		Width: &w,
		Color: &b,
	}
	for i := range d.VerticalRoads {
		r := d.VerticalRoads[i]
		drawingImpl.DrawLine(r.StartX, r.StartY, r.EndX, r.EndY, f)
	}
	for i := range d.HorizontalRoads {
		r := d.HorizontalRoads[i]
		drawingImpl.DrawLine(r.StartX, r.StartY, r.EndX, r.EndY, f)
	}
}

func (d *BoxesDocument) adjustLineToWidth(startX, startY, endX, endY, offset int, firstLine, lastLine bool) (int, int, int, int) {
	if startX == endX {
		// vertical line
		if startY < endY {
			// top/down
			if !firstLine {
				startY += -offset
			}
			if !lastLine {
				endY += offset
			}
		} else {
			//botton/up
			if !firstLine {
				endY += -offset
			}
			if !lastLine {
				startY += offset
			}
		}
	} else {
		// horizontal line
		if startX < endX {
			// left/right
			if !firstLine {
				startX += -offset
			}
			if !lastLine {
				endX += offset
			}
		} else {
			// right/left
			if !firstLine {
				startX += offset
			}
			if !lastLine {
				endX += -offset
			}
		}
	}
	return startX, startY, endX, endY
}

func (d *BoxesDocument) DrawConnections(drawingImpl types.Drawing) error {
	// DEBUG - start
	// b := "blue"
	// w := 2.0
	// format := types.LineDef{
	// 	Width: &w,
	// 	Color: &b,
	// }

	// for _, elem := range d.Connections {
	// 	// iterate over the connections of the document
	// 	lineFormat := format
	// 	if elem.Format != nil {
	// 		lineFormat = *elem.Format
	// 	}
	// 	offset := int(*lineFormat.Width / 2)
	// 	lastPos := len(elem.Parts) - 1
	// 	for i, l := range elem.Parts {
	// 		// drawing the connection lines
	// 		x1, y1, x2, y2 := d.adjustLineToWidth(l.StartX, l.StartY, l.EndX, l.EndY, offset, i == 0, i == lastPos)
	// 		// DEBUG - Start
	// 		if x1 == x2 {
	// 			// vertical line
	// 			y1 += 3
	// 			y2 -= 3
	// 		} else {
	// 			// horizontal line
	// 			x1 += 3
	// 			x2 -= 3
	// 		}
	// 		// DEBUG - End
	// 		drawingImpl.DrawLine(x1, y1, x2, y2, lineFormat)
	// 	}
	// }
	// DEBUG - end

	// this helps in the UI to highlight related connections
	d.AdjustConnectionIndexForCommentsWithSameLabel()
	d.DrawMovedConnectionLines(drawingImpl)
	return nil
}

func (d *BoxesDocument) AdjustConnectionIndexForCommentsWithSameLabel() {
	labelFindings := make(map[string][]int)
	for i := range d.Connections {
		c := d.Connections[i]
		if c.Comment != nil && c.Comment.Label != nil {
			l := *c.Comment.Label
			var newFindings []int
			if alreadyFound, ok := labelFindings[l]; ok {
				newFindings = append(alreadyFound, c.ConnectionIndex)
			} else {
				newFindings = []int{c.ConnectionIndex}
			}
			labelFindings[l] = newFindings
		}
	}
	for _, v := range labelFindings {
		if len(v) > 1 {
			newConnectionIndex, indexesToReplace := v[0], v[1:]
			d.replaceConnectionIds(newConnectionIndex, indexesToReplace)
		}
	}
}

func (d *BoxesDocument) replaceConnectionIdsNow(newConnectionIndex int, indexesToReplace []int, conLines *[]ConnectionLine) {
	for i := range *conLines {
		l := &(*conLines)[i]
		if slices.Contains(indexesToReplace, l.ConnectionIndex) {
			l.ConnectionIndex = newConnectionIndex
		}
	}
}

func (d *BoxesDocument) replaceConnectionIdsInComments(newConnectionIndex int, indexesToReplace []int) {
	for i := range d.Comments {
		c := &d.Comments[i]
		if c.ConnectionIndex != nil {
			if slices.Contains(indexesToReplace, *c.ConnectionIndex) {
				*c.ConnectionIndex = newConnectionIndex
			}

		}
	}
}

func (d *BoxesDocument) replaceConnectionIds(newConnectionIndex int, indexesToReplace []int) {
	d.replaceConnectionIdsNow(newConnectionIndex, indexesToReplace, &d.HorizontalLines)
	d.replaceConnectionIdsNow(newConnectionIndex, indexesToReplace, &d.VerticalLines)
	d.replaceConnectionIdsInComments(newConnectionIndex, indexesToReplace)
}

func (d *BoxesDocument) DrawMovedConnectionLines(drawingImpl types.Drawing) {
	b := "black"
	w := 1.0
	s := types.LineDefStyleEnum_dashed
	format := types.LineDef{
		Width: &w,
		Color: &b,
		Style: &s,
	}
	for i := range d.HorizontalLines {
		l := d.HorizontalLines[i]
		lineFormat := format
		if l.Format != nil {
			lineFormat = *l.Format
		}
		offset := int(*lineFormat.Width / 2)
		x1, y1, x2, y2 := d.adjustLineToWidth(l.StartX, l.StartY, l.EndX, l.EndY, offset, l.IsStart, l.IsEnd)
		classNames := fmt.Sprintf("connection conLine_%d", l.ConnectionIndex)
		if l.InverseDirection {
			x1, x2 = x2, x1
		}
		drawingImpl.DrawLineWithClass(x1, y1, x2, y2, lineFormat, classNames)
	}
	for i := range d.VerticalLines {
		l := d.VerticalLines[i]
		lineFormat := format
		if l.Format != nil {
			lineFormat = *l.Format
		}
		classNames := fmt.Sprintf("connection conLine_%d", l.ConnectionIndex)
		x1, y1, x2, y2 := l.StartX, l.StartY, l.EndX, l.EndY
		if l.InverseDirection {
			y1, y2 = y2, y1
		}
		drawingImpl.DrawLineWithClass(x1, y1, x2, y2, lineFormat, classNames)
	}
}

func (doc *BoxesDocument) AdjustDocHeight(le *LayoutElement, currentMax int) int {
	if le != &doc.Boxes {
		if le.Y+le.Height > currentMax {
			currentMax = le.Y + le.Height
		}
	} else {
		currentMax = doc.Boxes.Height
	}
	if le.Vertical != nil {
		for _, elem := range le.Vertical.Elems {
			currentMax = doc.AdjustDocHeight(&elem, currentMax)
		}
	}
	if le.Horizontal != nil {
		for _, elem := range le.Horizontal.Elems {
			currentMax = doc.AdjustDocHeight(&elem, currentMax)
		}
	}
	return currentMax
}

func (doc *BoxesDocument) GetTitleFormat() types.FontDef {
	format2Use := doc.Formats["default"].FontCaption // risky doesn't check if default exists
	if doc.TitleFormat != nil {
		if titleFormat, ok := doc.Formats[*doc.TitleFormat]; ok {
			format2Use = titleFormat.FontCaption
		}
	}
	format2Use.Anchor = types.FontDefAnchorEnum_left
	return format2Use
}

func (doc *BoxesDocument) DrawTitle(drawing types.Drawing, c types.TextDimensionCalculator) error {
	if doc.Title != "" {
		format2Use := doc.GetTitleFormat()
		//_, h := c.Dimensions(doc.Title, &format2Use)
		x := doc.Boxes.X
		y := doc.Boxes.Y + doc.Boxes.Height
		drawing.DrawText(doc.Title, x, y, 0, &format2Use)
	}
	return nil
}
func (doc *BoxesDocument) DrawLegend(drawing types.Drawing, c types.TextDimensionCalculator) error {
	if doc.Legend != nil {
		currentY := doc.Boxes.Y + doc.Boxes.Height
		if len(doc.Legend.Entries) > 0 {
			format2Use := doc.Formats["default"].FontCaption // risky doesn't check if default exists
			if doc.Legend.Format != nil {
				if legendFormat, ok := doc.Formats[*doc.Legend.Format]; ok {
					format2Use = legendFormat.FontCaption
				}
			}
			format2Use.Anchor = types.FontDefAnchorEnum_left
			currentX := doc.GlobalPadding
			lineW := 0.5
			lineC := "#9a9a9a"
			lineFormat := types.LineDef{
				Width: &lineW,
				Color: &lineC,
			}
			currentY += (2 * types.GlobalPadding)
			drawing.DrawLine(currentX, currentY-10, currentX+doc.Boxes.Width, currentY-10, lineFormat)
			for i := range doc.Legend.Entries {
				e := doc.Legend.Entries[i]
				if e.Format == "" {
					// no format found - the entry isn't included in the legend
					continue
				}
				legendFormat, ok := doc.Formats[e.Format]
				if !ok {
					// no suitable format found - the entry isn't included in the legend
					continue
				}
				text := e.Text + ": "
				w, _ := c.Dimensions(text, &format2Use)
				if (currentX + w) >= (doc.Boxes.X + doc.Boxes.Width) {
					currentX = doc.GlobalPadding
					currentY += types.GlobalPadding
				}
				drawing.DrawText(text, currentX, currentY, 0, &format2Use)
				currentX += w
				if legendFormat.Fill != nil {
					// draw example box
					currentX += 5
					c := "#606060"
					line := types.LineDef{
						Color: &c,
					}
					drawing.DrawSolidRect(currentX, currentY, 6, 10, legendFormat.Fill, &line)
					currentX += 8
				} else if legendFormat.Line != nil {
					// draw line example
					currentX += 5
					drawing.DrawLine(currentX, currentY+2, currentX+6, currentY+2, *legendFormat.Line)
					drawing.DrawLine(currentX+6, currentY+2, currentX+6, currentY+6, *legendFormat.Line)
					drawing.DrawLine(currentX+6, currentY+6, currentX+12, currentY+6, *legendFormat.Line)
					currentX += 14
				} else {
					// dummy output
					currentX += 5
					drawing.DrawText("?", currentX, currentY, 0, &format2Use) // Dummy
					currentX += 6
				}
				drawing.DrawText(",", currentX, currentY, 0, &format2Use)
				currentX += doc.GlobalPadding
			}
			doc.LegendEndY = currentY + doc.GlobalPadding
		}
	}
	return nil
}

func (doc *BoxesDocument) AdjustDocHeightTLegend(c types.TextDimensionCalculator) error {
	if doc.Legend != nil {
		if len(doc.Legend.Entries) > 0 {
			format2Use := doc.Formats["default"].FontCaption // risky doesn't check if default exists
			if doc.Legend.Format != nil {
				if legendFormat, ok := doc.Formats[*doc.Legend.Format]; ok {
					format2Use = legendFormat.FontCaption
				}
			}
			currentX := doc.GlobalPadding
			currentY := doc.Boxes.Y + doc.Boxes.Height + (2 * types.GlobalPadding)
			for i := range doc.Legend.Entries {
				e := doc.Legend.Entries[i]
				text := e.Text + ": "
				w, _ := c.Dimensions(text, &format2Use)
				if (currentX + w) >= (doc.Boxes.X + doc.Boxes.Width) {
					currentX = doc.GlobalPadding
					currentY += types.GlobalPadding
				}
				currentX += w
				currentX += 6
				currentX += doc.GlobalPadding
			}
			if currentY > doc.Height {
				doc.Height = currentY + doc.GlobalPadding
			}
		}
	}
	return nil
}

func (b *LayoutElement) Draw(drawing types.Drawing) error {
	if b.Format != nil {
		f := types.RectWithTextFormat{
			FontCaption:  b.Format.FontCaption,
			FontText1:    b.Format.FontText1,
			FontText2:    b.Format.FontText2,
			Padding:      b.Format.Padding,
			Border:       b.Format.Line,
			Fill:         b.Format.Fill,
			VerticalTxt:  b.Format.VerticalTxt,
			CornerRadius: b.Format.CornerRadius,
		}
		textYOffset := 0
		if b.Image != nil {
			textYOffset = (b.Image.Y - b.Y) + b.Image.Height + b.Image.MarginTopBottom
		}
		if err := drawing.DrawRectWithText(b.Id, b.Caption, b.Text1, b.Text2, b.X, b.Y, b.Width, b.Height, textYOffset, f); err != nil {
			return fmt.Errorf("Error drawing element %s: %w", b.Id, err)
		}
	}
	if b.Image != nil {
		if b.DataLink != nil && *b.DataLink != "" {
			if err := drawing.DrawPngWithAdditionalLink(b.Image.X, b.Image.Y, b.Image.ImgId, *b.DataLink); err != nil {
				return fmt.Errorf("Error drawing image %s: %w", b.Id, err)
			}

		} else {
			if err := drawing.DrawPng(b.Image.X, b.Image.Y, b.Image.ImgId); err != nil {
				return fmt.Errorf("Error drawing image %s: %w", b.Id, err)
			}
		}
	}
	if b.Vertical != nil {
		for i := range b.Vertical.Elems {
			elem := b.Vertical.Elems[i]
			if err := elem.Draw(drawing); err != nil {
				return err
			}
		}
	}
	if b.Horizontal != nil {
		// adjust here horizontal images, that are partial moved based on overlapping connections ...
		// ... this is only a cosmetical correction
		b.adjustHorizontalImagePositionsInCase()
		for i := range b.Horizontal.Elems {
			elem := b.Horizontal.Elems[i]
			if err := elem.Draw(drawing); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *LayoutElement) hasOnlyHorizontalImages() (bool, int) {
	w := 0
	for i := range b.Horizontal.Elems {
		elem := b.Horizontal.Elems[i]
		if elem.Vertical != nil || elem.Horizontal != nil {
			return false, 0
		}
		if elem.Image == nil || elem.Caption != "" || elem.Text1 != "" || elem.Text2 != "" {
			return false, 0
		}
		w += elem.Image.Width
	}
	return true, w
}

func (b *LayoutElement) adjustHorizontalImagePositionsInCase() {
	if ok, width := b.hasOnlyHorizontalImages(); ok {
		m := types.GlobalMinBoxMargin
		if b.Format.MinBoxMargin > 0 {
			m = b.Format.MinBoxMargin
		}
		width += (len(b.Horizontal.Elems) - 1) * m
		xStart := b.X + ((b.Width - width) / 2)
		for i := range b.Horizontal.Elems {
			elem := b.Horizontal.Elems[i]
			if elem.Image != nil {
				elem.Image.X = xStart
				xStart += elem.Image.Width + m
			}
		}
	}
}

func (b *LayoutElement) DrawTextBoxes(drawing types.Drawing) error {
	if b.Format != nil {
		c := "black"
		w := 2.0
		o := 0.4
		f := types.LineDef{
			Color:   &c,
			Width:   &w,
			Opacity: &o,
		}
		if b.XTextBox != nil {
			if err := drawing.DrawSolidRect(*b.XTextBox, *b.YTextBox, *b.WidthTextBox, *b.HeightTextBox, nil, &f); err != nil {
				return fmt.Errorf("Error drawing text boxes element %s: %w", b.Id, err)
			}
		}
	}
	if b.Vertical != nil {
		for _, elem := range b.Vertical.Elems {
			if err := elem.DrawTextBoxes(drawing); err != nil {
				return err
			}
		}
	}
	if b.Horizontal != nil {
		for _, elem := range b.Horizontal.Elems {
			if err := elem.DrawTextBoxes(drawing); err != nil {
				return err
			}
		}
	}
	return nil
}

func (doc *BoxesDocument) DrawComments(drawing types.Drawing, c types.TextDimensionCalculator) error {
	if len(doc.Comments) > 0 {
		doc.drawCommentMarkers(drawing)
		currentY, _ := doc.drawCommentTextsCustomLabels(drawing, c)
		_ = currentY
		doc.drawCommentTextsStdLabels(currentY, drawing, c)
	}
	return nil
}

func (doc *BoxesDocument) DrawOverlays(drawing types.Drawing, c types.TextDimensionCalculator) error {
	if len(doc.Overlays) > 0 {
		for i := range doc.Overlays {
			o := doc.Overlays[i]
			for j := range o.Layouts {
				oe := doc.Overlays[i].Layouts[j]
				if oe.Format.Fill == nil || oe.Format.Line == nil {
					continue
				}
				drawing.DrawCircleWithBorder(oe.X, oe.Y, int(oe.Radius), oe.Format.Fill, oe.Format.Line)
				if o.PrintValue {
					label := fmt.Sprintf("%d", int(oe.Value))
					size := int(oe.Radius / 5)
					if size < 14 {
						size = 14
					}
					labelFont := types.InitFontDef(nil, "sans", size, true, false, 0)
					labelFont.Color = *oe.Format.Fill.Color
					labelFont.Anchor = types.FontDefAnchorEnum_middle
					_, h := c.Dimensions(label, &labelFont)
					drawing.DrawText(label, oe.X, oe.Y-(h/2), 0, &labelFont)
				}
			}
		}
	}
	return nil
}

func (doc *BoxesDocument) drawCommentMarkers(drawing types.Drawing) error {
	for i := range doc.Comments {
		c := doc.Comments[i]
		if c.ConnectionIndex != nil {
			className := fmt.Sprintf("connection conLine_%d", *c.ConnectionIndex)
			drawing.DrawCircleWithBorderTextAndClass(c.Label, c.MarkerX, c.MarkerY, doc.CommentMarkerRadius, &c.Format.Fill, &c.Format.Line, &c.Format.FontMarker, className)
		} else {
			drawing.DrawCircleWithBorderAndText(c.Label, c.MarkerX, c.MarkerY, doc.CommentMarkerRadius, &c.Format.Fill, &c.Format.Line, &c.Format.FontMarker)
		}
	}
	return nil
}

func (doc *BoxesDocument) drawCommentTextsCustomLabels(drawing types.Drawing, c types.TextDimensionCalculator) (int, error) {
	markerX := doc.GlobalPadding + doc.CommentMarkerRadius
	textX := markerX + (2 * doc.CommentMarkerRadius)
	currentY := doc.LegendEndY
	if currentY == 0 {
		currentY = doc.Boxes.Y + doc.Boxes.Height
	}
	currentY += doc.CommentMarkerRadius
	neededMarkerSpace := doc.CommentMarkerRadius + doc.GlobalPadding + 3
	customLabeledComments := make([]*CommentContainer, 0)
	for i := range doc.Comments {
		c := doc.Comments[i]
		if c.CustomMarker {
			customLabeledComments = append(customLabeledComments, &c)
		}
	}
	slices.SortFunc(customLabeledComments, func(c1, c2 *CommentContainer) int {
		if c1.Label == c2.Label {
			return 0
		} else if c1.Label < c2.Label {
			return -1
		} else {
			return 1
		}
	})

	var lastLabel, lastText string
	for i := range customLabeledComments {
		c := customLabeledComments[i]
		if c.Text == "" {
			continue
		}
		if c.Text == lastText && c.Label == lastLabel {
			continue
		}
		lastLabel, lastText = c.Label, c.Text
		className := "comment"
		if c.ConnectionIndex != nil {
			className = fmt.Sprintf("%s connection conLine_%d", className, *c.ConnectionIndex)
		}
		drawing.DrawCircleWithBorderTextAndClass(c.Label, markerX, currentY, doc.CommentMarkerRadius, &c.Format.Fill, &c.Format.Line, &c.Format.FontMarker, className)
		c.Format.FontText.Anchor = types.FontDefAnchorEnum_left
		c.Format.FontText.MaxLenBeforeBreak = doc.Boxes.Width
		drawing.DrawText(c.Text, textX, currentY-(2*doc.CommentMarkerRadius), 0, &c.Format.FontText)
		currentY += getMax(c.TextHeight, neededMarkerSpace)
	}
	return currentY, nil
}

func (doc *BoxesDocument) drawCommentTextsStdLabels(currentY int, drawing types.Drawing, c types.TextDimensionCalculator) (int, error) {
	markerX := doc.GlobalPadding + doc.CommentMarkerRadius
	textX := markerX + (2 * doc.CommentMarkerRadius)
	currentY += doc.CommentMarkerRadius
	neededMarkerSpace := doc.CommentMarkerRadius + doc.GlobalPadding + 3
	stdLabeledComments := make([]*CommentContainer, 0)
	for i := range doc.Comments {
		c := doc.Comments[i]
		if !c.CustomMarker {
			stdLabeledComments = append(stdLabeledComments, &c)
		}
	}
	slices.SortFunc(stdLabeledComments, func(c1, c2 *CommentContainer) int {
		if c1.Label == c2.Label {
			return 0
		} else if c1.Label < c2.Label {
			return -1
		} else {
			return 1
		}
	})

	for i := range stdLabeledComments {
		c := stdLabeledComments[i]
		if c.Text == "" {
			continue
		}
		c.Format.FontText.Anchor = types.FontDefAnchorEnum_left
		c.Format.FontText.MaxLenBeforeBreak = doc.Boxes.Width
		className := "comment"
		if c.ConnectionIndex != nil {
			className = fmt.Sprintf("%s connection conLine_%d", className, *c.ConnectionIndex)
		}
		drawing.DrawCircleWithBorderTextAndClass(c.Label, markerX, currentY, doc.CommentMarkerRadius, &c.Format.Fill, &c.Format.Line, &c.Format.FontMarker, className)
		drawing.DrawText(c.Text, textX, currentY-(2*doc.CommentMarkerRadius), 0, &c.Format.FontText)
		currentY += getMax(c.TextHeight, neededMarkerSpace)
	}
	return currentY, nil
}
