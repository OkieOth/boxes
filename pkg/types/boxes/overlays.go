package boxes

import "github.com/okieoth/boxes/pkg/types"

func (doc *BoxesDocument) IncludeOverlays(c types.TextDimensionCalculator) error {
	if len(doc.Overlays) > 0 {
		minX, minY, maxX, maxY := 0, 0, doc.Width, doc.Height
		for overlayIndex := range doc.Overlays {
			overlay := doc.Overlays[overlayIndex]
			for layoutKey, _ := range overlay.Layouts {
				overlayEntry := overlay.Layouts[layoutKey]
				doc.initOverlayForLayout(&doc.Boxes, &overlay, layoutKey, &overlayEntry)
				radiusAsInt := int(overlayEntry.Radius)
				minX = getMin(overlayEntry.X-radiusAsInt, minX)
				minY = getMin(overlayEntry.Y-radiusAsInt, minY)
				maxX = getMax(overlayEntry.X+radiusAsInt, maxX)
				maxY = getMax(overlayEntry.Y+radiusAsInt, maxY)
				overlay.Layouts[layoutKey] = overlayEntry
			}
		}
		// maybe not needed for the general nice appearance
		// adjust in case the document size
		// if minX < 0 {
		// 	// extend the document to the right and move all elements to the right
		// 	diff := absInt(minX) + doc.GlobalPadding
		// 	doc.Width += diff
		// 	doc.MoveAllElementsToRight(diff)
		// }
		// if maxX < doc.Width {
		// 	// extend only the document to the right
		// 	diff := absInt(minX) + doc.GlobalPadding
		// 	doc.Boxes.Width += diff
		// 	doc.Width += diff
		// }
		// if minY < 0 {
		// 	diff := absInt(minY) + doc.GlobalPadding
		// 	doc.Height += diff
		// 	doc.MoveAllElementsDown(diff)
		// }
		// if maxY < 0 {
		// 	diff := absInt(maxY) + doc.GlobalPadding
		// 	doc.Height += diff
		// 	doc.Boxes.Height += diff
		// }
	}
	return nil
}

func (doc *BoxesDocument) initOverlayForLayout(l *LayoutElement, overlay *DocOverlay, entryKey string, overlayEntry *OverlayEntry) bool {
	if l.Id != "" && (l.Id == entryKey || l.Caption == entryKey) {
		// found the related layout element
		doc.initOverlayXY(l, overlay, overlayEntry)
		doc.initOverlayRadius(overlay, overlayEntry)
		doc.initOverlayFormat(overlay, overlayEntry)
		return true
	}
	if doc.initOverlaysFromLayoutCont(l.Horizontal, overlay, entryKey, overlayEntry) {
		return true
	}
	return doc.initOverlaysFromLayoutCont(l.Vertical, overlay, entryKey, overlayEntry)
}

func (doc *BoxesDocument) initOverlayXY(l *LayoutElement, overlay *DocOverlay, overlayEntry *OverlayEntry) {
	halfWidth := (l.Width / 2)
	halfHeight := (l.Height / 2)
	overlayEntry.X = l.X + halfWidth
	overlayEntry.Y = l.Y + halfHeight
	if overlay.CenterXOffset != 0 {
		diff := float64(halfWidth) * overlay.CenterXOffset
		overlayEntry.X += int(diff)
	}
	if overlay.CenterYOffset != 0 {
		diff := float64(halfHeight) * overlay.CenterYOffset
		overlayEntry.Y += int(diff)
	}
}

func (doc *BoxesDocument) initOverlayRadius(overlay *DocOverlay, overlayEntry *OverlayEntry) {
	if overlay.RefValue == 0.0 {
		overlayEntry.Radius = overlayEntry.Value
	} else {
		overlayEntry.Radius = (overlayEntry.Value / overlay.RefValue) * overlay.RadiusDefs.Max
	}
	if overlay.RadiusDefs.Min > 0 && overlayEntry.Radius < overlay.RadiusDefs.Min {
		overlayEntry.Radius = overlay.RadiusDefs.Min
	}
	if overlay.RadiusDefs.Max > 0 && overlayEntry.Radius > overlay.RadiusDefs.Max {
		overlayEntry.Radius = overlay.RadiusDefs.Max
	}
}

func (doc *BoxesDocument) getDefaultOverlayFormat() BoxFormat {
	return BoxFormat{
		Line: types.InitLineDef2("blue", 0.6),
		Fill: types.InitFillDef2("blue", 0.4),
	}
}

func (doc *BoxesDocument) initOverlayFormat(overlay *DocOverlay, overlayEntry *OverlayEntry) {
	var format BoxFormat
	if overlay.Formats.Default != "" {
		f, ok := doc.Formats[overlay.Formats.Default]
		if !ok {
			f = doc.getDefaultOverlayFormat()
		}
		format = f
	} else {
		format = doc.getDefaultOverlayFormat()
	}

	if len(overlay.Formats.Gradations) > 0 {
		var lastFormat, format2Use string
		for i := range overlay.Formats.Gradations {
			g := overlay.Formats.Gradations[i]
			if g.Limit > overlayEntry.Value {
				format2Use = lastFormat
				break
			}
			lastFormat = g.Format
		}
		if format2Use != "" {
			if f, ok := doc.Formats[format2Use]; ok {
				format = f
			}
		}
	}
	overlayEntry.Format = format
}

func (doc *BoxesDocument) initOverlaysFromLayoutCont(cont *LayoutElemContainer, overlay *DocOverlay, entryKey string, overlayEntry *OverlayEntry) bool {
	if cont != nil {
		for i := range cont.Elems {
			if doc.initOverlayForLayout(&cont.Elems[i], overlay, entryKey, overlayEntry) {
				return true
			}
		}
	}
	return false
}

func (doc *BoxesDocument) MoveAllElementsToRight(diff int) {
	doc.moveLayoutToRight(&doc.Boxes, diff)
	for i := range doc.HorizontalLines {
		l := doc.HorizontalLines[i]
		l.StartX += diff
		l.EndX += diff
	}
	for i := range doc.VerticalLines {
		l := doc.VerticalLines[i]
		l.StartX += diff
		l.EndX += diff
	}
	for i := range doc.Comments {
		c := doc.Comments[i]
		c.MarkerX += diff
	}
	for i := range doc.Overlays {
		for j := range doc.Overlays[i].Layouts {
			oe := doc.Overlays[i].Layouts[j]
			oe.X += diff
		}
	}
}

func addIntIfNotNil(v *int, i int) {
	if v != nil {
		*v += i
	}
}

func (doc *BoxesDocument) moveLayoutToRight(l *LayoutElement, diff int) {
	addIntIfNotNil(l.BottomXToStart, diff)
	addIntIfNotNil(l.TopXToStart, diff)
	addIntIfNotNil(l.XTextBox, diff)
	l.CenterX += diff
	l.X += diff
	doc.moveLayoutContToRight(l.Horizontal, diff)
	doc.moveLayoutContToRight(l.Vertical, diff)
}

func (doc *BoxesDocument) moveLayoutContToRight(cont *LayoutElemContainer, diff int) {
	if cont != nil {
		for i := range cont.Elems {
			doc.moveLayoutToRight(&cont.Elems[i], diff)
		}
	}
}

func (doc *BoxesDocument) MoveAllElementsDown(diff int) {
	doc.moveLayoutDown(&doc.Boxes, diff)
	for i := range doc.HorizontalLines {
		l := doc.HorizontalLines[i]
		l.StartY += diff
		l.EndY += diff
	}
	for i := range doc.VerticalLines {
		l := doc.VerticalLines[i]
		l.StartY += diff
		l.EndY += diff
	}
	for i := range doc.Comments {
		c := doc.Comments[i]
		c.MarkerY += diff
	}
	for i := range doc.Overlays {
		for j := range doc.Overlays[i].Layouts {
			oe := doc.Overlays[i].Layouts[j]
			oe.Y += diff
		}
	}
}

func (doc *BoxesDocument) moveLayoutDown(l *LayoutElement, diff int) {
	addIntIfNotNil(l.LeftYToStart, diff)
	addIntIfNotNil(l.RightYToStart, diff)
	addIntIfNotNil(l.YTextBox, diff)
	l.CenterY += diff
	l.Y += diff
	doc.moveLayoutContDown(l.Horizontal, diff)
	doc.moveLayoutContDown(l.Vertical, diff)
}

func (doc *BoxesDocument) moveLayoutContDown(cont *LayoutElemContainer, diff int) {
	if cont != nil {
		for i := range cont.Elems {
			doc.moveLayoutDown(&cont.Elems[i], diff)
		}
	}
}
