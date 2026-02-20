package boxes

import (
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (l *LayoutElement) incrementX(xOffset int) {
	l.X += xOffset
	if l.XTextBox != nil {
		*l.XTextBox += xOffset
	}
	if l.Image != nil {
		l.Image.X += xOffset
	}
	if l.Vertical != nil {
		l.Vertical.X += xOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementX(xOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.X += xOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementX(xOffset)
		}
	}
}

func (l *LayoutElement) incrementY(yOffset int) {
	l.Y += yOffset
	if l.YTextBox != nil {
		*l.YTextBox += yOffset
	}
	if l.Image != nil {
		l.Image.Y += yOffset
	}
	if l.Vertical != nil {
		l.Vertical.Y += yOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementY(yOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.Y += yOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementY(yOffset)
		}
	}
}

func (l *LayoutElement) centerHorizontalElems() {
	if l.Horizontal != nil {
		cont := l.Horizontal
		contYOffset := (l.Height - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			offsetY := ((cont.Height - sub.Height) / 2)
			sub.incrementY(offsetY)
			sub.incrementX(contXOffset)
			//sub.centerHorizontalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) centerVerticalElems() {
	if l.Vertical != nil {
		cont := l.Vertical
		contYOffset := (l.Y + l.Height - cont.Y - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			sub.incrementY(contYOffset)
			offsetX := contXOffset + ((cont.Width - sub.Width) / 2)
			sub.incrementX(offsetX)
			//sub.centerVerticalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) AreOnTheSameVerticalLevel(otherElem *LayoutElement) bool {
	return ((l.CenterY > otherElem.Y) && (l.CenterY < otherElem.Y+otherElem.Height)) ||
		((otherElem.CenterY > l.Y) && (otherElem.CenterY < l.Y+l.Height))
}

type ConnDirection int

const (
	ConnDirectionLeft ConnDirection = iota
	ConnDirectionUp
	ConnDirectionRight
	ConnDirectionDown
)

func (l *LayoutElement) IsInYRange(y1, y2 int) bool {
	if y1 < y2 {
		return l.CenterY > y1 && l.CenterY < y2
	}
	return l.CenterY < y1 && l.CenterY > y2
}

func (l *LayoutElement) IsInXRange(x1, x2 int) bool {
	if x1 < x2 {
		return l.CenterX > x1 && l.CenterX < x2
	}
	return l.CenterX < x1 && l.CenterX > x2
}

func between(value, min, max int) bool {
	return value > min && value < max
}

func (l *LayoutElement) initVertical(c types.TextDimensionCalculator, yInnerOffset int) {
	if l.Vertical != nil && len(l.Vertical.Elems) > 0 {
		curX := l.X
		l.Vertical.X = curX
		curY := l.Y + yInnerOffset
		l.Vertical.Y = curY
		var w int
		lv := len(l.Vertical.Elems)
		margin := types.GlobalMinBoxMargin
		if l.Format != nil {
			margin = l.Format.MinBoxMargin
		}
		for i := range lv {
			sub := &l.Vertical.Elems[i]
			marginToUse := margin
			if sub.Caption == "" && sub.Text1 == "" && sub.Text2 == "" && sub.Image != nil {
				// in case sub contains only a picture, then the image margin overrides the format margin
				marginToUse = 0
			} else if sub.Format != nil && sub.Format.MinBoxMargin != 0 {
				marginToUse = sub.Format.MinBoxMargin
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c)
			if marginToUse > 0 {
				curY += (sub.Height + marginToUse)
			} else {
				curY += (sub.Height + types.GlobalPadding)
			}
			if sub.Width > w {
				w = sub.Width
			}
			if sub.Width > l.Vertical.Width {
				l.Vertical.Width = sub.Width
			}
		}
		for i := range l.Vertical.Elems {
			l.Vertical.Elems[i].Width = w
		}

		l.Vertical.Height = l.Vertical.Elems[lv-1].Y + l.Vertical.Elems[lv-1].Height - l.Vertical.Y
		l.Height += l.Vertical.Height
		padding := types.GlobalPadding
		if l.Format != nil {
			padding = l.Format.Padding
		}
		l.adjustDimensionsBasedOnNested(w, padding)
		l.Height += types.GlobalPadding
	}
}

func (l *LayoutElement) adjustDimensionsBasedOnNested(width, padding int) {
	// TODO: remove later if it proves as working
	// if l.Caption != "" || l.Text1 != "" || l.Text2 != "" {
	// 	l.Height += padding
	// }
	if width > l.Width {
		if l.Caption != "" || l.Text1 != "" || l.Text2 != "" || l.Image != nil {
			l.Width = width + (2 * padding)
		} else {
			l.Width = width
		}
	}
}

func (l *LayoutElement) initHorizontal(c types.TextDimensionCalculator, yInnerOffset int) {
	if l.Horizontal != nil && len(l.Horizontal.Elems) > 0 {
		curX := l.X
		l.Horizontal.X = curX
		curY := l.Y + yInnerOffset
		l.Horizontal.Y = curY
		var h, w int
		var hasChilds bool
		margin := types.GlobalMinBoxMargin
		if l.Format != nil {
			margin = l.Format.MinBoxMargin
		}
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			if (sub.Horizontal != nil && len(sub.Horizontal.Elems) > 0) || (sub.Vertical != nil && len(sub.Vertical.Elems) > 0) {
				hasChilds = true
			}
			marginToUse := margin
			if sub.Caption == "" && sub.Text1 == "" && sub.Text2 == "" && sub.Image != nil {
				// in case sub contains only a picture, then the image margin overrides the format margin
				marginToUse = 0
			} else if sub.Format != nil && sub.Format.MinBoxMargin != 0 {
				marginToUse = sub.Format.MinBoxMargin
			}
			if w > 0 {
				w += marginToUse
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c)
			curX += (sub.Width + marginToUse)
			w += sub.Width
			if sub.Height > h {
				h = sub.Height
			}
			if sub.Height > l.Horizontal.Height {
				l.Horizontal.Height = sub.Height
			}
		}
		if !hasChilds {
			for i := 0; i < len(l.Horizontal.Elems); i++ {
				sub := &l.Horizontal.Elems[i]
				sub.Height = h
			}
		}

		l.Horizontal.Height = h
		l.Horizontal.Width = w
		l.Height += l.Horizontal.Height

		l.adjustDimensionsBasedOnNested(w, margin)
		l.Height += types.GlobalPadding
	}
}

type FilterFunc func(l *LayoutElement, currentDepth int) bool

var dummyFilter = func(l *LayoutElement, currentDepth int) bool {
	return true
}

func getMax(v1, v2 int) int {
	if v2 > v1 {
		return v2
	}
	return v1
}

func getMin(v1, v2 int) int {
	if v2 < v1 {
		return v2
	}
	return v1
}

func (l *LayoutElement) InitDimensions(c types.TextDimensionCalculator) {
	var cW, cH, t1W, t1H, t2W, t2H, textWidth, textHeight int
	//var yCaptionOffset, yText1Offset, yText2Offset, yInnerOffset int
	var yInnerOffset int
	padding := types.GlobalPadding
	if l.Format != nil && l.Format.Padding > 0 {
		padding = l.Format.Padding
	}
	yTextBox := l.Y + padding
	if l.Format != nil && l.Format.Padding > 0 {
		padding = l.Format.Padding
	}
	if l.Image != nil {
		w := (l.Image.Width + (2 * l.Image.MarginLeftRight))
		h := l.Image.Height + (2 * l.Image.MarginTopBottom)
		l.Image.Y = l.Y + l.Image.MarginTopBottom
		l.Height += h
		if l.Width < w {
			l.Width = w
		}
		yInnerOffset += h
		yTextBox = l.Y + h
	}
	p := types.GlobalPadding
	if l.Format != nil && l.Format.Padding > 0 {
		p = l.Format.Padding
	}
	if l.Caption != "" || l.Text1 != "" || l.Text2 != "" {
		if l.Caption != "" {
			cW, cH = c.Dimensions(l.Caption, &l.Format.FontCaption)
			if !l.Format.VerticalTxt {
				l.Height += cH
				l.Height += l.Format.FontCaption.SpaceBottom
				l.Height += p
				if l.Text1 == "" && l.Text2 == "" && l.Vertical == nil && l.Horizontal == nil {
					p := l.Format.Padding
					if l.Format.FontCaption.SpaceBottom > 0 {
						p = l.Format.FontCaption.SpaceBottom
					}
					l.Height += p
				}
			} else {
				// vertical text
				cW, cH = cH, cW
				l.Width += cW + p + l.Format.FontCaption.SpaceBottom
				// if l.Text1 == "" && l.Text2 == "" {
				// 	l.Width += l.Format.Padding
				// }
			}
			textWidth = cW
			textHeight = cH
		}
		if l.Text1 != "" {
			p := l.Format.Padding
			if l.Format.FontText1.SpaceTop > 0 {
				p = l.Format.FontText1.SpaceTop
			}
			t1W, t1H = c.Dimensions(l.Text1, &l.Format.FontText1)
			if !l.Format.VerticalTxt {
				l.Height += t1H
				l.Height += l.Format.FontText1.SpaceBottom
				l.Height += p
				textHeight += p
				textWidth = getMax(textWidth, t1W)
				textHeight += t1H
				if l.Text2 == "" && l.Vertical == nil && l.Horizontal == nil {
					p := l.Format.Padding
					if l.Format.FontText1.SpaceBottom > 0 {
						p = l.Format.FontText1.SpaceBottom
					}
					l.Height += p
				}
			} else {
				t1W, t1H = t1H, t1W
				l.Width += t1W + p + l.Format.FontText1.SpaceBottom
				if l.Text2 == "" {
					l.Width += l.Format.Padding
				}
				textWidth += t1W
				textHeight = getMax(textHeight, t1H)
			}

		}
		if l.Text2 != "" {
			p := l.Format.Padding
			if l.Format.FontText2.SpaceTop > 0 {
				p = l.Format.FontText2.SpaceTop
			}
			t2W, t2H = c.Dimensions(l.Text2, &l.Format.FontText2)
			if !l.Format.VerticalTxt {
				l.Height += t2H
				l.Height += p
				textWidth = getMax(textWidth, t2W)
				textHeight += t2H
				if l.Vertical == nil && l.Horizontal == nil {
					p := l.Format.Padding
					if l.Format.FontText2.SpaceBottom > 0 {
						p = l.Format.FontText2.SpaceBottom
					}
					l.Height += p
				}
			} else {
				t2W, t2H = t2H, t2W
				l.Width += t2W + p + l.Format.FontText2.SpaceBottom
				l.Width += l.Format.Padding
				textWidth += t2W
				textHeight = getMax(textHeight, t2H)
			}
		}
		if !l.Format.VerticalTxt {
			// normal horizontal text
			h := l.Height
			if l.Format.FixedHeight != nil {
				h = *l.Format.FixedHeight
			}
			l.Height = l.adjustToRaster(h)
			yInnerOffset = l.Height + types.RasterSize
			var w int
			if l.Format.FixedWidth != nil {
				w = *l.Format.FixedWidth
			} else {
				w = max(cW, max(t1W, t2W))
			}
			l.Width = l.adjustToRaster(w + (2 * l.Format.Padding))
		} else {
			// vertical text
			w := l.Width
			if l.Format.FixedWidth != nil {
				w = *l.Format.FixedWidth
			}
			l.Width = l.adjustToRaster(w)
			yInnerOffset = l.Width
			var h int
			if l.Format.FixedHeight != nil {
				h = *l.Format.FixedHeight
			} else {
				h = max(cH, max(t1H, t2H)) + (2 * l.Format.Padding)
			}
			l.Height = l.adjustToRaster(h)
		}
		l.WidthTextBox = &textWidth
		l.HeightTextBox = &textHeight
	} else if l.Format != nil {
		if l.Format.FixedHeight != nil {
			l.Height = l.adjustToRaster(*l.Format.FixedHeight)
		}
		if l.Format.FixedWidth != nil {
			l.Width = l.adjustToRaster(*l.Format.FixedWidth)
		}
	}
	if l.Vertical != nil {
		l.initVertical(c, yInnerOffset)
	}
	if l.Horizontal != nil {
		l.initHorizontal(c, yInnerOffset)
	}
	l.adjustToParentWidth(l.Vertical)   // doesn't go recursive through the children
	l.adjustToParentWidth(l.Horizontal) // doesn't go recursive through the children

	xTextBox := l.X + (l.Width-textWidth)/2
	l.XTextBox = &xTextBox
	l.YTextBox = &yTextBox
	if l.Image != nil {
		l.Image.X = l.X + ((l.Width - l.Image.Width - l.Image.MarginLeftRight - l.Image.MarginLeftRight) / 2)
	}
}

func (l *LayoutElement) adjustToParentWidth(cont *LayoutElemContainer) {
	if cont != nil {
		for i := range len(cont.Elems) {
			e := &cont.Elems[i]
			if e.Format != nil && e.Format.WidthOfParent != nil && *e.Format.WidthOfParent {
				e.Width = cont.Width
				e.X = cont.X
				if e.XTextBox != nil {
					x := e.X + ((e.Width - *e.WidthTextBox) / 2)
					e.XTextBox = &x
				}
			}
		}
	}
}

func (l *LayoutElement) adjustToRaster(value int) int {
	if value > 0 {
		rasterRest := value % (types.RasterSize * 2)
		return value + ((types.RasterSize * 2) - rasterRest)
	}
	return value
}

func (l *LayoutElement) Center() {
	l.CenterX = l.X + (l.Width / 2)
	if l.WidthTextBox != nil {
		*l.XTextBox = l.X + ((l.Width - *l.WidthTextBox) / 2)
	}
	l.CenterY = l.Y + (l.Height / 2)
	l.centerHorizontalElems()
	l.centerVerticalElems()
}
